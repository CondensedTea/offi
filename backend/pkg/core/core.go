package core

import (
	"fmt"
	"offi/pkg/cache"
	"offi/pkg/etf2l"
	"os"
	"time"

	"github.com/go-co-op/gocron"
)

type Core struct {
	cache     cache.Cache
	etf2l     *etf2l.Client
	scheduler *gocron.Scheduler

	enableErrorCaching bool
}

func New(cache cache.Cache, etf2l *etf2l.Client) *Core {
	_, ok := os.LookupEnv("DISABLE_ERROR_CACHE")

	return &Core{
		cache:     cache,
		etf2l:     etf2l,
		scheduler: gocron.NewScheduler(time.UTC),

		enableErrorCaching: !ok,
	}
}

func (c Core) StartScheduler() error {
	_, err := c.scheduler.Every(4).Hours().Do(c.loadPlayersRecruitmentPosts)
	if err != nil {
		return fmt.Errorf("failed to schedule load players rec posts scrapper: %v", err)
	}

	_, err = c.scheduler.Every(4).Hours().Do(c.loadTeamsRecruitmentPosts)
	if err != nil {
		return fmt.Errorf("failed to schedule load team rec posts scrapper job: %v", err)
	}

	c.scheduler.StartAsync()
	return nil
}
