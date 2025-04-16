package service

import (
	"context"
	"log/slog"
	"net/http"
	"offi/internal/cache"
	"offi/internal/db"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"
	"offi/internal/logstf"

	"errors"

	"github.com/jackc/pgx/v5"
)

type Cache interface {
	SetLogError(ctx context.Context, matchId int, err error) error
	CheckLogError(ctx context.Context, matchId int) error

	GetPlayer(ctx context.Context, playerID int) (cache.Player, error)
	SetPlayer(ctx context.Context, playerID int, player cache.Player) error
}

type database interface {
	Begin(ctx context.Context) (pgx.Tx, error)

	GetLastRecruitmentForAuthor(ctx context.Context, postType db.Post, authorID int) (db.Recruitment, error)

	SaveLog(ctx context.Context, tx pgx.Tx, log db.Log) error
	MatchExists(ctx context.Context, mathcID int) (bool, error)
	GetLogsByMatchID(ctx context.Context, matchID int) ([]db.Log, error)
	GetMatchByLogID(ctx context.Context, logID int) (db.Match, error)
	SaveMatchTx(ctx context.Context, tx pgx.Tx, match db.Match) error
	SaveMatch(ctx context.Context, match db.Match) error
}

type Service struct {
	gen.UnimplementedHandler

	cache Cache
	db    database
	etf2l *etf2l.Client
	logs  *logstf.Client
}

func NewService(cache Cache, db database, etf2lClient *etf2l.Client, logs *logstf.Client) *Service {
	return &Service{
		cache: cache,
		db:    db,
		etf2l: etf2lClient,
		logs:  logs,
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
