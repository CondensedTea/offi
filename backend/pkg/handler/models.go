package handler

import "offi/pkg/cache"

type GetPlayersResponse struct {
	Players []cache.Player `json:"players"`
}

type GetTeamResponse struct {
	Team cache.Team `json:"team"`
}
