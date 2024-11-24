package service

import (
	"context"
	"errors"
	"net/http"
	"offi/internal/cache"
	gen "offi/internal/gen/api"
)

func (s *Service) GetMatchForLog(_ context.Context, params gen.GetMatchForLogParams) (r gen.GetMatchForLogRes, _ error) {
	matchPage, err := s.cache.GetMatch(params.LogID)
	if err != nil {
		if errors.Is(err, cache.ErrCached) {
			return &gen.ErrorStatusCode{
				StatusCode: http.StatusTooEarly,
				Response:   gen.Error{Error: err.Error()},
			}, nil
		}

		return nil, err
	}

	return &gen.GetMatchForLogOK{
		Match: gen.Match{
			MatchID:     matchPage.Id,
			Competition: matchPage.Competition,
			Stage:       matchPage.Stage,
		},
	}, nil
}