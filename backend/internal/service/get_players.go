package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"offi/internal/cache"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"

	"github.com/go-redis/redis"
	"github.com/samber/lo"
)

func (s *Service) GetPlayers(ctx context.Context, params gen.GetPlayersParams) (r *gen.GetPlayersOK, _ error) {
	players, err := s.getPlayers(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &gen.GetPlayersOK{
		Players: players,
	}, nil
}

func (s *Service) getPlayers(_ context.Context, playerIDs []int) ([]gen.Player, error) {
	var players []gen.Player

	for _, playerID := range playerIDs {
		player, err := s.cache.GetPlayer(playerID)
		switch {
		case err == redis.Nil:
			etf2lPlayer, etf2lErr := s.etf2l.GetPlayer(playerID)
			switch {
			case errors.Is(etf2lErr, etf2l.ErrPlayerNotFound):
				slog.Warn("failed to get player from etf2l", "player_id", playerID, "error", etf2lErr)
				continue
			case etf2lErr != nil:
				return nil, fmt.Errorf("failed to get player from etf2l: %v", etf2lErr)
			}

			player = etf2lPlayer.ToCache()
			if cacheErr := s.cache.SetPlayer(playerID, player); cacheErr != nil {
				return nil, fmt.Errorf("failed to save player to cache: %v", cacheErr)
			}
		case err != nil:
			return nil, fmt.Errorf("failed to get player from cache: %v", err)
		}

		bans := lo.Map(player.Bans, func(ban cache.PlayerBan, i int) gen.PlayerBan {
			return gen.PlayerBan{
				Start:  ban.Start,
				End:    ban.End,
				Reason: ban.Reason,
			}
		})

		apiPlayer := gen.Player{
			ID:      player.ID,
			SteamID: player.SteamID,
			Name:    player.Name,
			Bans:    bans,
		}

		// TODO: recruitment info

		players = append(players, apiPlayer)
	}

	return players, nil
}
