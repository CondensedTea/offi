package core

import (
	"offi/pkg/cache"
	"offi/pkg/etf2l"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func (c Core) loadTeamsRecruitmentPosts() {
	logrus.Info("loading recruitment posts for teams")

	entries, err := c.etf2l.LoadRecruitmentPosts(etf2l.TeamPost)
	if err != nil {
		logrus.Errorf("failed to load recruitment posts from etf2l: %v", err)
		return
	}
	cacheEntries := lo.Map[etf2l.Recruitment, cache.Entry](entries, func(entry etf2l.Recruitment, i int) cache.Entry {
		return entry.ToCache()
	})

	if err = c.cache.SaveRecruitmentPosts("team", cacheEntries); err != nil {
		logrus.Errorf("failed to save recruitment team posts: %v", err)
		return
	}
}

func (c Core) loadPlayersRecruitmentPosts() {
	logrus.Info("loading recruitment posts for players")

	entries, err := c.etf2l.LoadRecruitmentPosts(etf2l.PlayerPost)
	if err != nil {
		logrus.Errorf("failed to load recruitment posts from etf2l: %v", err)
		return
	}
	cacheEntries := lo.Map[etf2l.Recruitment, cache.Entry](entries, func(entry etf2l.Recruitment, i int) cache.Entry {
		return entry.ToCache()
	})

	if err = c.cache.SaveRecruitmentPosts(etf2l.PlayerPost, cacheEntries); err != nil {
		logrus.Errorf("failed to save recruitment team posts: %v", err)
		return
	}
}
