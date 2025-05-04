package service

import (
	"context"
	"errors"
	"fmt"
	"offi/internal/cache"
	"offi/internal/db"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"
	"unsafe"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func (s *Service) GetPlayers(ctx context.Context, p gen.GetPlayersParams) (r *gen.GetPlayersOK, _ error) {
	players, err := s.getPlayers(ctx, p.ID, p.WithRecruitmentStatus.Or(false))
	if err != nil {
		return nil, err
	}

	return &gen.GetPlayersOK{
		Players: players,
	}, nil
}

func (s *Service) getPlayers(ctx context.Context, playerIDs []int, withRecruitmentStatus bool) ([]gen.Player, error) {
	var players []gen.Player

	for _, playerID := range playerIDs {
		player, err := s.cache.GetPlayer(ctx, cache.LeagueETF2L, int64(playerID))
		switch {
		case errors.Is(err, redis.Nil):
			etf2lPlayer, err := s.etf2l.GetPlayer(ctx, playerID)
			switch {
			case errors.Is(err, etf2l.ErrPlayerNotFound):
				if cacheErr := s.cache.SetPlayer(ctx, cache.LeagueETF2L, int64(playerID), cache.Player{DoesntExists: true}); cacheErr != nil {
					return nil, fmt.Errorf("failed to save unknown player to cache: %w", cacheErr)
				}
				continue
			case err != nil:
				return nil, fmt.Errorf("failed to get player %d from etf2l: %w", playerID, err)
			}

			player = etf2lPlayer.ToCache()
			if cacheErr := s.cache.SetPlayer(ctx, cache.LeagueETF2L, int64(playerID), player); cacheErr != nil {
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

		apiPlayer := gen.Player{
			ID:      player.ID,
			SteamID: player.SteamID,
			Name:    player.Name,
			Bans:    bans,
		}

		if withRecruitmentStatus {
			apiPlayer.Recruitment, err = s.getRecruitmentStatusForPlayer(ctx, int64(playerID))
			if err != nil {
				return nil, err
			}
		}

		players = append(players, apiPlayer)
	}

	return players, nil
}

func (s *Service) getRecruitmentStatusForPlayer(ctx context.Context, playerID int64) (gen.OptRecruitmentInfo, error) {
	recruitment, err := s.db.GetLastRecruitmentForAuthor(ctx, db.Player, playerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return gen.OptRecruitmentInfo{Set: false}, nil
		}

		return gen.OptRecruitmentInfo{}, fmt.Errorf("failed to get recruitments for player %d: %w", playerID, err)
	}

	return gen.NewOptRecruitmentInfo(gen.RecruitmentInfo{
		Skill:    recruitment.SkillLevel,
		URL:      fmt.Sprintf("https://etf2l.org/recruitment/%d/", recruitment.RecruitmentID),
		Classes:  *(*[]gen.GameClass)(unsafe.Pointer(&recruitment.Classes)),
		GameMode: recruitment.TeamType,
	}), nil
}
