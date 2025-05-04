package service

import (
	"context"
	"errors"
	"fmt"
	"offi/internal/cache"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func (s *Service) GetETF2LPlayers(ctx context.Context, p gen.GetETF2LPlayersParams) (r *gen.GetETF2LPlayersOK, _ error) {
	players, err := s.getETF2LPlayers(ctx, p.ID, p.WithRecruitmentStatus.Or(false))
	if err != nil {
		return nil, err
	}

	return &gen.GetETF2LPlayersOK{
		Players: players,
	}, nil
}

func (s *Service) getETF2LPlayers(ctx context.Context, playerIDs []int64, withRecruitmentStatus bool) ([]gen.ETF2LPlayer, error) {
	var players []gen.ETF2LPlayer

	for _, playerID := range playerIDs {
		player, err := s.cache.GetPlayer(ctx, cache.LeagueETF2L, playerID)
		switch {
		case errors.Is(err, redis.Nil):
			etf2lPlayer, err := s.etf2l.GetPlayer(ctx, int(playerID))
			switch {
			case errors.Is(err, etf2l.ErrPlayerNotFound):
				if cacheErr := s.cache.SetPlayer(ctx, cache.LeagueETF2L, playerID, cache.Player{DoesntExists: true}); cacheErr != nil {
					return nil, fmt.Errorf("failed to save unknown player to cache: %w", cacheErr)
				}
				continue
			case err != nil:
				return nil, fmt.Errorf("failed to get player %d from etf2l: %w", playerID, err)
			}

			player = etf2lPlayer.ToCache()
			if cacheErr := s.cache.SetPlayer(ctx, cache.LeagueETF2L, playerID, player); cacheErr != nil {
				return nil, fmt.Errorf("failed to save player to cache: %w", cacheErr)
			}
		case err != nil:
			return nil, fmt.Errorf("failed to get player from cache: %w", err)
		}

		if player.DoesntExists {
			continue
		}

		bans := make([]gen.PlayerBan, len(player.Bans))
		for i, ban := range player.Bans {
			bans[i] = gen.PlayerBan{
				Start:  ban.Start,
				End:    ban.End,
				Reason: ban.Reason,
			}
		}

		steamID64, err := strconv.ParseInt(player.SteamID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse steam ID %s: %w", player.SteamID, err)
		}

		apiPlayer := gen.ETF2LPlayer{
			ID:      int64(player.ID),
			SteamID: steamID64,
			Name:    player.Name,
			Bans:    bans,
		}

		if withRecruitmentStatus {
			apiPlayer.Recruitment, err = s.getRecruitmentStatusForPlayer(ctx, playerID)
			if err != nil {
				return nil, err
			}
		}

		players = append(players, apiPlayer)
	}

	return players, nil
}
