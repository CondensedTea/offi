package core

import (
	"fmt"
	"offi/pkg/cache"
	"offi/pkg/etf2l"
	"offi/pkg/logstf"
	"time"

	"github.com/go-co-op/gocron"
)

type Core struct {
	cache     cache.Cache
	etf2l     *etf2l.Client
	logsTf    *logstf.Client
	scheduler *gocron.Scheduler
}

func New(cache cache.Cache, etf2l *etf2l.Client, logsTf *logstf.Client) (*Core, error) {
	c := &Core{
		cache:     cache,
		etf2l:     etf2l,
		logsTf:    logsTf,
		scheduler: gocron.NewScheduler(time.UTC),
	}
	return c, nil
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
