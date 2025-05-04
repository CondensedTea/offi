package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"offi/internal/cache"
	"offi/internal/db"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"
	"time"
)

var ErrTooManyPlayers = errors.New("too many players, 18 or less allowed")

const maxPlayers = 18

func (s *Service) GetLogsForMatch(ctx context.Context, params gen.GetLogsForMatchParams) (r gen.GetLogsForMatchRes, _ error) {
	logs, err := s.getLogsForMatch(ctx, params.MatchID)
	if err != nil {
		if errors.Is(err, cache.ErrCached) {
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
	exists, err := s.db.MatchExists(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match %d from cache: %w", matchID, err)
	}

	if !exists {
		if err = s.cache.CheckLogError(ctx, matchID); err != nil {
			return nil, err
		}
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

func (s *Service) saveNewMatch(ctx context.Context, matchID int) ([]db.Log, error) {
	match, err := s.etf2l.GetMatch(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get etf2l match: %w", err)
	}

	if len(match.PlayerSteamIDs) > maxPlayers {
		// logs.tf does not allow searching for more than 18 players, so we save only save the match
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

	matchLogs, secondaryLogs, err := s.logs.SearchLogs(ctx, match.PlayerSteamIDs, match.Maps, match.SubmittedAt)
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
