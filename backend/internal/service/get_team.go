package service

import (
	"context"
	"net/http"
	"offi/internal/gen/api"
)

func (s *Service) GetTeam(_ context.Context, _ api.GetTeamParams) (r *api.GetTeamOK, _ error) {
	return nil, &api.ErrorStatusCode{
		StatusCode: http.StatusNotImplemented,
		Response:   api.Error{Error: "not implemented"},
	}
}
