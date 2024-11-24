package service

import (
	"context"
	"net/http"
	"offi/internal/cache"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"
)

type Cache interface {
	Ping() error

	GetLogs(matchId int) (cache.LogSet, error)
	SetLogs(matchId int, match *cache.LogSet) error

	SetLogError(matchId int, err error) error
	CheckLogError(matchId int) error

	DeleteLogs(matchId int) (*cache.LogSet, error)

	GetPlayer(playerID int) (cache.Player, error)
	SetPlayer(playerID int, player cache.Player) error

	GetTeam(teamID int) (cache.Team, error)
	SetTeam(teamID int, team cache.Team) error

	GetMatch(logId int) (cache.MatchPage, error)
	SetMatch(logIds []int, matchPage *cache.MatchPage) error

	GetAllKeys(hashKey string) ([]string, error)

	IncrementViews(object string, id int) (int64, error)
	GetViews(object string, id int) (int64, error)
}

type Service struct {
	gen.UnimplementedHandler

	cache              Cache
	etf2l              *etf2l.Client
	enableErrorCaching bool
}

func NewService(cache Cache, etf2lClient *etf2l.Client, cacheErrors bool) *Service {
	return &Service{
		cache:              cache,
		etf2l:              etf2lClient,
		enableErrorCaching: cacheErrors,
	}
}

func (s *Service) NewError(_ context.Context, err error) (r *gen.ErrorStatusCode) {
	return &gen.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response:   gen.Error{Error: err.Error()},
	}
}