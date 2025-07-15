package service

import (
	"context"
	"fmt"
	gen "offi/internal/gen/api"
	"offi/internal/redis"
)

func (s *Service) GetRGLPlayers(ctx context.Context, p gen.GetRGLPlayersParams) (r *gen.GetRGLPlayersOK, _ error) {
	var players []gen.RGLPlayer

	cachePlayers, err := s.cache.GetPlayers(ctx, redis.LeagueRGL, p.ID)
	if err != nil {
		return nil, fmt.Errorf("getting players from cache: %w", err)
	}

	var playersToResolve []int64
	for playerID, player := range cachePlayers {
		if player == nil {
			playersToResolve = append(playersToResolve, playerID)
		} else {
			players = append(players, gen.RGLPlayer{
				SteamID: player.SteamID,
				Name:    player.Name,
			})
		}
	}

	resolvedPlayers, err := s.rgl.GetPlayers(ctx, playersToResolve)
	if err != nil {
		return nil, fmt.Errorf("getting players from RGL: %w", err)
	}

	for _, player := range resolvedPlayers {
		if err = s.cache.SetPlayer(ctx, redis.LeagueRGL, player.SteamID, redis.Player{
			SteamID: player.SteamID,
			Name:    player.Name,
		}); err != nil {
			return nil, fmt.Errorf("saving player %d to cache: %w", player.SteamID, err)
		}

		players = append(players, gen.RGLPlayer{
			SteamID: player.SteamID,
			Name:    player.Name,
		})

		// TODO: decide if should cache DoesntExists players as well
	}

	return &gen.GetRGLPlayersOK{Players: players}, nil
}
