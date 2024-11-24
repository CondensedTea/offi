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
	"offi/internal/logstf"
	"strconv"

	"github.com/go-redis/redis"
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
				Response:   gen.Error{Error: err.Error()},
			}, nil
		}

		return nil, err
	}

	res := lo.Map(logs, func(log cache.Log, _ int) gen.Log {
		return gen.Log{
			ID:          log.ID,
			Title:       log.Title,
			Map:         log.Map,
			PlayedAt:    log.PlayedAt,
			IsSecondary: log.IsSecondary,
		}
	})

	return &gen.GetLogsForMatchOK{
		Logs: res,
	}, nil
}

func (s *Service) getLogsForMatch(ctx context.Context, matchID int) ([]cache.Log, error) {
	logSet, err := s.cache.GetLogs(matchID)
	switch {
	case err == redis.Nil:
		if s.enableErrorCaching {
			if storedErr := s.cache.CheckLogError(matchID); storedErr != nil {
				return nil, storedErr
			}
		}
		logs, saveErr := s.saveNewMatch(ctx, matchID)
		if saveErr != nil {
			if s.enableErrorCaching {
				if cacheErr := s.cache.SetLogError(matchID, saveErr); cacheErr != nil {
					slog.Error("failed to cache log error", "error", cacheErr)
				}
			}
			return nil, saveErr
		}
		return logs, nil
	case err != nil:
		return nil, fmt.Errorf("failed to get match from cache: %v", err)
	}

	return logSet.Logs, nil
}

func (s *Service) saveNewMatch(ctx context.Context, matchId int) ([]cache.Log, error) {
	logIDs := make([]int, 0)

	match, err := s.etf2l.GetMatch(matchId)
	if err != nil {
		return nil, fmt.Errorf("failed to get players for etf2l match: %w", err)
	}

	if len(match.Players) > maxPlayers {
		return nil, ErrTooManyPlayers
	}

	playerIDs := lo.Map(match.Players, func(playerID string, _ int) int {
		id, _ := strconv.Atoi(playerID)
		return id
	})

	players, err := s.getPlayers(ctx, playerIDs)
	if err != nil {
		return nil, err
	}

	steamIDs := lo.Map(players, func(player gen.Player, _ int) string {
		return player.SteamID
	})

	var cacheLogs []cache.Log

	matchLogs, secondaryLogs, err := logstf.SearchLogs(steamIDs, match.Maps, match.PlayedAt)
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
	if err = s.cache.SetLogs(matchId, &cache.LogSet{Logs: cacheLogs}); err != nil {
		return nil, fmt.Errorf("failed to set match in cache: %v", err)
	}
	if err = s.cache.SetMatch(logIDs, &cache.MatchPage{
		Id:          match.ID,
		Competition: match.Competition,
		Stage:       match.Stage,
	}); err != nil {
		return nil, fmt.Errorf("failed to set logs in cache: %v", err)
	}
	return cacheLogs, nil
}
