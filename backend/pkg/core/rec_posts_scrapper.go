package core

import (
	"offi/pkg/cache"
	"offi/pkg/etf2l"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func (c Core) loadTeamsRecruitmentPosts() {
	entries, err := c.etf2l.LoadRecruitmentPosts(etf2l.TeamPost)
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
	entries, err := c.etf2l.LoadRecruitmentPosts(etf2l.PlayerPost)
	if err != nil {
		logrus.Errorf("failed to load recruitment posts from etf2l: %v", err)
		return
	}

	var player etf2l.Player
	lo.ForEach[etf2l.Recruitment](entries, func(entry etf2l.Recruitment, i int) {
		player, err = c.etf2l.GetPlayer(entry.Id)
		if err != nil {
			logrus.Errorf("failed to get player from cache: %v", err)
			return
		}

		cachePlayer := player.ToCache()
		cachePlayer.Recruitment = entry.ToCache()
		if err = c.cache.SetPlayer(entry.Id, cachePlayer); err != nil {
			logrus.Errorf("failed to save recruitment team posts: %v", err)
		}
	})
	logrus.Info("loaded recruitment posts for players")
}
