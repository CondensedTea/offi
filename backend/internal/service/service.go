package service

import (
	"context"
	"log/slog"
	"net/http"
	"offi/internal/cache"
	"offi/internal/db"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"

	"errors"
)

type Cache interface {
	GetLogs(matchId int) (cache.LogSet, error)
	SetLogs(matchId int, match *cache.LogSet) error

	SetLogError(matchId int, err error) error
	CheckLogError(matchId int) error

	GetPlayer(playerID int) (cache.Player, error)
	SetPlayer(playerID int, player cache.Player) error

	GetMatch(logId int) (cache.MatchPage, error)
	SetMatch(logIds []int, matchPage *cache.MatchPage) error

	GetAllKeys(hashKey string) ([]string, error)
}

type database interface {
	GetLastRecruitmentForAuthor(ctx context.Context, postType db.Post, authorID int) (db.Recruitment, error)
}

type Service struct {
	gen.UnimplementedHandler

	cache              Cache
	db                 database
	etf2l              *etf2l.Client
	enableErrorCaching bool
}

func NewService(cache Cache, db database, etf2lClient *etf2l.Client, cacheErrors bool) *Service {
	return &Service{
		cache:              cache,
		db:                 db,
		etf2l:              etf2lClient,
		enableErrorCaching: cacheErrors,
	}
}

func (s *Service) NewError(_ context.Context, err error) (r *gen.ErrorStatusCode) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return &gen.ErrorStatusCode{
			StatusCode: http.StatusRequestTimeout,
			Response:   gen.Error{Error: err.Error()},
		}
	case errors.Is(err, etf2l.ErrIncompleteMatch):
		return &gen.ErrorStatusCode{
			StatusCode: http.StatusTooEarly,
			Response:   gen.Error{Error: err.Error()},
		}
	default:
		slog.Error("unexpected error", "error", err, "component", "api")

		return &gen.ErrorStatusCode{
			StatusCode: http.StatusInternalServerError,
			Response:   gen.Error{Error: err.Error()},
		}
	}
}
