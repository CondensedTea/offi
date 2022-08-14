package core

import (
	"context"
	"offi/pkg/cache"
	"offi/pkg/etf2l"
	"strconv"
	"time"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

const loadJobTimeout = 5 * time.Minute

func (c Core) loadTeamsRecruitmentPosts() {
	ctx, cancel := context.WithTimeout(context.Background(), loadJobTimeout)
	defer cancel()

	entries, err := etf2l.LoadRecruitmentPosts(ctx, etf2l.TeamPost)
	if err != nil {
		logrus.Errorf("failed to load recruitment posts from etf2l: %v", err)
		return
	}
	lo.ForEach[etf2l.Recruitment](entries, func(entry etf2l.Recruitment, i int) {
		team := cache.Team{
			ID:          entry.Id,
			Recruitment: entry.ToCache(),
		}
		if err = c.cache.SetTeam(entry.Id, team); err != nil {
			logrus.Errorf("failed to save recruitment team posts: %v", err)
		}
	})
	logrus.Info("loaded recruitment posts for teams")
}

func (c Core) loadPlayersRecruitmentPosts() {
	ctx, cancel := context.WithTimeout(context.Background(), loadJobTimeout)
	defer cancel()

	entries, err := etf2l.LoadRecruitmentPosts(ctx, etf2l.PlayerPost)
	if err != nil {
		logrus.Errorf("failed to load recruitment posts from etf2l: %v", err)
		return
	}

	var player etf2l.Player
	lo.ForEach[etf2l.Recruitment](entries, func(entry etf2l.Recruitment, i int) {
		player, err = etf2l.GetPlayer(ctx, entry.Id)
		if err != nil {
			logrus.Errorf("failed to get player from cache: %v", err)
			return
		}

		cachePlayer := player.ToCache()
		cachePlayer.Recruitment = entry.ToCache()

		steamID64, _ := strconv.Atoi(cachePlayer.SteamID)
		if err = c.cache.SetPlayer(steamID64, cachePlayer); err != nil {
			logrus.Errorf("failed to save recruitment team posts: %v", err)
		}
	})
	logrus.Info("loaded recruitment posts for players")
}
