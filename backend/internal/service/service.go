package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"offi/internal/cache"
	"offi/internal/db"
	"offi/internal/demostf"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"
	"offi/internal/logstf"
	"time"

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
	UpdateDemoIDForLog(ctx context.Context, logID int, demoID int) error
}

type demostfClient interface {
	FindDemo(ctx context.Context, req demostf.FindDemoRequest) (demostf.Demo, error)
}

type Service struct {
	gen.UnimplementedHandler

	cache Cache
	db    database
	etf2l *etf2l.Client
	logs  *logstf.Client
	demos demostfClient

	resolveDemoQueue chan resolveDemoRequest
}

func NewService(ctx context.Context, cache Cache, db database, etf2lClient *etf2l.Client, logs *logstf.Client, demo demostfClient) *Service {
	s := &Service{
		cache:            cache,
		db:               db,
		etf2l:            etf2lClient,
		logs:             logs,
		demos:            demo,
		resolveDemoQueue: make(chan resolveDemoRequest, 3),
	}

	s.startDemoResolver(ctx)

	return s
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

type resolveDemoRequest struct {
	logID          int
	playerSteamIDs []int
	playedAt       time.Time
	mapName        string
}

func (s *Service) startDemoResolver(ctx context.Context) {
	for range 2 {
		go func() {
			for {
				select {
				case req := <-s.resolveDemoQueue:
					cctx, cancel := context.WithTimeout(ctx, 5*time.Second)

					if err := s.resolveDemo(cctx, req); err != nil {
						slog.Error("failed to resolve demo", "error", err, "log_id", req.logID)
					}

					cancel()

				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func (s *Service) resolveDemo(ctx context.Context, req resolveDemoRequest) error {
	demo, err := s.demos.FindDemo(ctx, demostf.FindDemoRequest{
		PlayerIDs: req.playerSteamIDs,
		PlayedAt:  req.playedAt,
		Map:       req.mapName,
	})
	if err != nil {
		if errors.Is(err, demostf.ErrNotFound) {
			return nil
		}

		return fmt.Errorf("finding demo at demos.tf: %w", err)

	}

	if err = s.db.UpdateDemoIDForLog(ctx, req.logID, demo.ID); err != nil {
		return fmt.Errorf("updating demo id: %w", err)
	}

	return nil
}
