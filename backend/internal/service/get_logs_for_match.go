package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"offi/internal/cache"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
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
			ID:          log.ID,
			Title:       log.Title,
			Map:         log.Map,
			PlayedAt:    log.PlayedAt,
			IsSecondary: log.IsSecondary,
		}
	}

	return &gen.GetLogsForMatchOK{
		Logs: res,
	}, nil
}

func (s *Service) getLogsForMatch(ctx context.Context, matchID int) ([]cache.Log, error) {
	logSet, err := s.cache.GetLogs(ctx, matchID)
	switch {
	case errors.Is(err, redis.Nil):
		if s.enableErrorCaching { // todo make error checking default
			if storedErr := s.cache.CheckLogError(ctx, matchID); storedErr != nil {
				return nil, storedErr
			}
		}
		logs, saveErr := s.saveNewMatch(ctx, matchID)
		if saveErr != nil {
			if s.enableErrorCaching {
				if cacheErr := s.cache.SetLogError(ctx, matchID, saveErr); cacheErr != nil {
					slog.ErrorContext(ctx, "failed to cache log error", "error", cacheErr)
				}
			}
			return nil, fmt.Errorf("failed to save parsed match %d: %w", matchID, saveErr)
		}
		return logs, nil
	case err != nil:
		return nil, fmt.Errorf("failed to get match %d from cache: %w", matchID, err)
	}

	return logSet.Logs, nil
}

func (s *Service) saveNewMatch(ctx context.Context, matchId int) ([]cache.Log, error) {
	logIDs := make([]int, 0)

	match, err := s.etf2l.GetMatch(ctx, matchId)
	if err != nil {
		return nil, fmt.Errorf("failed to get players for etf2l match: %w", err)
	}

	if len(match.PlayerSteamIDs) > maxPlayers {
		return nil, ErrTooManyPlayers
	}

	players, err := s.getPlayers(ctx, match.PlayerSteamIDs, false)
	if err != nil {
		return nil, err
	}

	steamIDs := lo.Map(players, func(player gen.Player, _ int) string {
		return player.SteamID
	})

	var cacheLogs []cache.Log

	matchLogs, secondaryLogs, err := s.logs.SearchLogs(ctx, steamIDs, match.Maps, match.SubmittedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to search logs: %v", err)
	}
	for _, log := range matchLogs {
		cacheLog := log.ToCache(false)
		logIDs = append(logIDs, log.Id)
		cacheLogs = append(cacheLogs, cacheLog)
	}
	for _, log := range secondaryLogs {
		cacheLog := log.ToCache(true)
		logIDs = append(logIDs, log.Id)
		cacheLogs = append(cacheLogs, cacheLog)
	}
	if err = s.cache.SetLogs(ctx, matchId, &cache.LogSet{Logs: cacheLogs}); err != nil {
		return nil, fmt.Errorf("failed to set match in cache: %v", err)
	}
	if err = s.cache.SetMatch(ctx, logIDs, &cache.MatchPage{
		Id:          match.ID,
		Competition: match.Competition,
		Stage:       match.Stage,
		Tier:        match.Tier,
	}); err != nil {
		return nil, fmt.Errorf("failed to set logs in cache: %v", err)
	}
	return cacheLogs, nil
}
