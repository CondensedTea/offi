package service

import (
	"context"
	"fmt"
	"offi/internal/cache"
	gen "offi/internal/gen/api"
	"strconv"
)

func (s *Service) GetRGLPlayers(ctx context.Context, p gen.GetRGLPlayersParams) (r *gen.GetRGLPlayersOK, _ error) {
	var players []gen.RGLPlayer

	cachePlayers, err := s.cache.GetPlayers(ctx, cache.LeagueRGL, p.ID)
	if err != nil {
		return nil, fmt.Errorf("getting players from cache: %w", err)
	}

	var playersToResolve []int64
	for playerID, player := range cachePlayers {
		if player == nil {
			playersToResolve = append(playersToResolve, playerID)
		} else {
			steamID, _ := strconv.ParseInt(player.SteamID, 10, 64)
			players = append(players, gen.RGLPlayer{
				SteamID: steamID,
				Name:    player.Name,
			})
		}
	}

	resolvedPlayers, err := s.rgl.GetPlayers(ctx, playersToResolve)
	if err != nil {
		return nil, fmt.Errorf("getting players from RGL: %w", err)
	}

	for _, player := range resolvedPlayers {
		steamID, _ := strconv.ParseInt(player.SteamID, 10, 64)

		if err = s.cache.SetPlayer(ctx, cache.LeagueRGL, steamID, cache.Player{
			SteamID: player.SteamID,
			Name:    player.Name,
		}); err != nil {
			return nil, fmt.Errorf("saving player %d to cache: %w", steamID, err)
		}

		players = append(players, gen.RGLPlayer{
			SteamID: steamID,
			Name:    player.Name,
		})

		// TODO: decide if should cache DoesntExists players as well
	}

	return &gen.GetRGLPlayersOK{Players: players}, nil
}
