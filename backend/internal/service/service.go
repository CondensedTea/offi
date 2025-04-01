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
	GetLogs(ctx context.Context, matchId int) (cache.LogSet, error)
	SetLogs(ctx context.Context, matchId int, match *cache.LogSet) error

	SetLogError(ctx context.Context, matchId int, err error) error
	CheckLogError(ctx context.Context, matchId int) error

	GetPlayer(ctx context.Context, playerID int) (cache.Player, error)
	SetPlayer(ctx context.Context, playerID int, player cache.Player) error

	GetMatch(ctx context.Context, logId int) (cache.MatchPage, error)
	SetMatch(ctx context.Context, logIds []int, matchPage *cache.MatchPage) error

	GetAllKeys(ctx context.Context, hashKey string) ([]string, error)
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

func (s *Service) NewError(ctx context.Context, err error) (r *gen.ErrorStatusCode) {
	switch {
	case errors.Is(err, context.Canceled):
		return &gen.ErrorStatusCode{
			StatusCode: http.StatusTeapot,
			Response:   gen.Error{Error: err.Error()},
		}
	case errors.Is(err, etf2l.ErrIncompleteMatch):
		return &gen.ErrorStatusCode{
			StatusCode: http.StatusTooEarly,
			Response:   gen.Error{Error: err.Error()},
		}
	default:
		slog.ErrorContext(ctx, "unexpected error", "error", err, "component", "api")

		return &gen.ErrorStatusCode{
			StatusCode: http.StatusInternalServerError,
			Response:   gen.Error{Error: err.Error()},
		}
	}
}
