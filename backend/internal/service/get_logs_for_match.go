package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"offi/internal/db"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"
	"offi/internal/logstf"
	"offi/internal/redis"
	"time"
)

var ErrTooManyPlayers = errors.New("too many players, 18 or less allowed")

const maxPlayers = 18

func (s *Service) GetLogsForMatch(ctx context.Context, params gen.GetLogsForMatchParams) (r gen.GetLogsForMatchRes, _ error) {
	logs, err := s.getLogsForMatch(ctx, params.MatchID)
	if err != nil {
		if errors.Is(err, redis.ErrCached) {
			return &gen.ErrorStatusCode{
				StatusCode: http.StatusTooEarly,
				Response:   gen.Error{Error: err.Error()},
			}, nil
		}

		if errors.Is(err, ErrTooManyPlayers) {
			return &gen.ErrorStatusCode{
				StatusCode: http.StatusBadRequest,
				Response:   gen.Error{Error: err.Error()},
			}, nil
		}

		if errors.Is(err, etf2l.ErrMatchNotFound) {
			return &gen.ErrorStatusCode{
				StatusCode: http.StatusNotFound,
				Response:   gen.Error{Error: "no logs found for the match"},
			}, nil
		}

		return nil, err
	}

	res := make([]gen.Log, len(logs))
	for i, log := range logs {
		res[i] = gen.Log{
			ID:          log.LogID,
			Title:       log.Title,
			Map:         log.Map,
			PlayedAt:    log.PlayedAt,
			IsSecondary: log.IsSecondary,
			DemoID: gen.OptInt{
				Set:   log.DemoID.Valid,
				Value: log.DemoID.V,
			},
		}
	}

	return &gen.GetLogsForMatchOK{Logs: res}, nil
}

func (s *Service) getLogsForMatch(ctx context.Context, matchID int) (logs []db.Log, err error) {
	exists, err := s.matchExists(ctx, matchID)
	if err != nil {
		return nil, err
	}

	if !exists {
		if err = s.cache.CheckLogError(ctx, matchID); err != nil {
			return nil, err
		} // TODO: store together with match exists?
		logs, err = s.saveNewMatch(ctx, matchID)
		if err != nil {
			if cacheErr := s.cache.SetLogError(ctx, matchID, err); cacheErr != nil {
				slog.ErrorContext(ctx, "failed to cache log error", "error", cacheErr)
			}
			return nil, fmt.Errorf("failed to save parsed match %d: %w", matchID, err)
		}
		return logs, nil
	}

	logs, err = s.db.GetLogsByMatchID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for match %d: %w", matchID, err)
	}

	return logs, nil
}

func (s *Service) matchExists(ctx context.Context, matchID int) (bool, error) {
	existsCached, err := s.cache.MatchExists(ctx, matchID)
	if err != nil {
		return false, fmt.Errorf("checking if match %d exists in cache: %w", matchID, err)
	}

	if existsCached {
		return true, nil
	}

	exists, err := s.db.MatchExists(ctx, matchID)
	if err != nil {
		return false, fmt.Errorf("check if match %d exists in db: %w", matchID, err)
	}

	if exists {
		if err = s.cache.SaveMatchExists(ctx, matchID); err != nil {
			return false, fmt.Errorf("saving match exists in cache: %w", err)
		}
	}

	return exists, nil
}

func (s *Service) saveNewMatch(ctx context.Context, matchID int) ([]db.Log, error) {
	match, err := s.etf2l.GetMatch(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get etf2l match: %w", err)
	}

	if len(match.PlayerSteamIDs) > maxPlayers || match.CompetitionType == "1v1" {
		// logs.tf does not allow searching for more than 18 players, so we don't search logs and only save the match
		// TODO: api-v2 returns `team_id: null` for mercs. It is can be used to search logs without 19th/20th player.

		// There is no logs for MGE matches so we skip them altogether
		if err = s.db.SaveMatch(ctx, db.Match{
			MatchID:     matchID,
			Competition: match.Competition,
			Stage:       match.Stage,
			Tier:        match.Tier,
			CompletedAt: match.SubmittedAt,
		}); err != nil {
			return nil, fmt.Errorf("failed to save match %d: %w", matchID, err)
		}

		return []db.Log{}, nil
	}

	var format logstf.Format
	switch match.CompetitionType {
	case "Highlander":
		format = logstf.Format9v9
	case "6v6":
		format = logstf.Format6v6
	default:
		format = logstf.FormatUnknown
	}

	matchLogs, secondaryLogs, err := s.logs.SearchLogs(ctx, logstf.SearchLogsRequest{
		PlayerIDs: match.PlayerSteamIDs,
		Maps:      match.Maps,
		Format:    format,
		PlayedAt:  match.SubmittedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search logs: %w", err)
	}

	logs := make([]db.Log, 0, len(matchLogs)+len(secondaryLogs))
	for _, log := range matchLogs {
		logs = append(logs, db.Log{
			MatchID:     matchID,
			LogID:       log.ID,
			Title:       log.Title,
			Map:         log.Map,
			PlayedAt:    time.Unix(log.Date, 0),
			IsSecondary: false,
		})
	}

	for _, log := range secondaryLogs {
		logs = append(logs, db.Log{
			MatchID:     matchID,
			LogID:       log.ID,
			Title:       log.Title,
			Map:         log.Map,
			PlayedAt:    time.Unix(log.Date, 0),
			IsSecondary: true,
		})
	}

	if len(logs) > 20 {
		return nil, fmt.Errorf("too many logs found for match %d", matchID)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	for _, log := range logs {
		if err = s.db.SaveLog(ctx, tx, log); err != nil {
			return nil, fmt.Errorf("failed to save log %d: %w", log.LogID, err)
		}
	}

	err = s.db.SaveMatchTx(ctx, tx, db.Match{
		MatchID:     matchID,
		Competition: match.Competition,
		Stage:       match.Stage,
		Tier:        match.Tier,
		CompletedAt: match.SubmittedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save match %d: %w", matchID, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	if err = s.cache.SaveMatchExists(ctx, matchID); err != nil {
		return nil, fmt.Errorf("failed to save match exists in cache: %w", err)
	}

	go func() {
		for _, log := range logs {
			if !log.IsSecondary {
				s.resolveDemoQueue <- resolveDemoRequest{
					logID:          log.LogID,
					playerSteamIDs: match.PlayerSteamIDs,
					playedAt:       log.PlayedAt,
					mapName:        log.Map,
				}
			}
		}
	}()

	return logs, nil
}
