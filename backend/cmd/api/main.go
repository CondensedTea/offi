package main

import (
	"log"
	"offi/pkg/cache"
	"offi/pkg/core"
	"offi/pkg/etf2l"
	"offi/pkg/handler"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	if _, ok := os.LookupEnv("DEBUG"); ok {
		logrus.SetLevel(logrus.DebugLevel)
	}

	etf2lClient := etf2l.New()

	redisUrl := os.Getenv("REDIS_URL")
	cacheClient, err := cache.New(redisUrl)
	if err != nil {
		log.Fatal(err)
	}
	c := core.New(cacheClient, etf2lClient)

	if err = c.StartScheduler(); err != nil {
		logrus.Fatalf("failed to start scheduler: %v", err)
	}

	if err = handler.New(c).Run(); err != nil {
		logrus.Fatalf("failed to run handler: %v", err)
	}
}
