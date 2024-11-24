// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// GetLogsForMatch implements GetLogsForMatch operation.
//
// Get logs associated with ETF2L match.
//
// GET /match/{match_id}
func (UnimplementedHandler) GetLogsForMatch(ctx context.Context, params GetLogsForMatchParams) (r GetLogsForMatchRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetMatchForLog implements GetMatchForLog operation.
//
// Get logs associated with given ETF2L match ID.
//
// GET /log/{log_id}
func (UnimplementedHandler) GetMatchForLog(ctx context.Context, params GetMatchForLogParams) (r GetMatchForLogRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetPlayers implements GetPlayers operation.
//
// Get players by Steam IDs.
//
// GET /players
func (UnimplementedHandler) GetPlayers(ctx context.Context, params GetPlayersParams) (r *GetPlayersOK, _ error) {
	return r, ht.ErrNotImplemented
}

// GetTeam implements GetTeam operation.
//
// Get team details.
//
// GET /team/{id}
func (UnimplementedHandler) GetTeam(ctx context.Context, params GetTeamParams) (r *GetTeamOK, _ error) {
	return r, ht.ErrNotImplemented
}

// NewError creates *ErrorStatusCode from error returned by handler.
//
// Used for common default response.
func (UnimplementedHandler) NewError(ctx context.Context, err error) (r *ErrorStatusCode) {
	r = new(ErrorStatusCode)
	return r
}
