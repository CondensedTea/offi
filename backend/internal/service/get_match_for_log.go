package service

import (
	"context"
	"errors"
	"net/http"
	"offi/internal/db"
	gen "offi/internal/gen/api"
)

func (s *Service) GetMatchForLog(ctx context.Context, params gen.GetMatchForLogParams) (r gen.GetMatchForLogRes, _ error) {
	match, err := s.db.GetMatchByLogID(ctx, params.LogID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &gen.ErrorStatusCode{
				StatusCode: http.StatusNotFound,
				Response:   gen.Error{Error: err.Error()},
			}
		}

		return nil, err
	}

	return &gen.GetMatchForLogOK{
		Match: gen.Match{
			MatchID:     match.MatchID,
			Competition: match.Competition,
			Stage:       match.Stage,
			Tier:        match.Tier,
		},
		Log: gen.GetMatchForLogOKLog{
			DemoID: gen.OptInt{
				Set:   match.DemoID.Valid,
				Value: match.DemoID.V,
			},
		},
	}, nil
}
