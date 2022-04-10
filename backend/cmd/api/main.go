package main

import (
	"log"
	"offi/pkg/cache"
	"offi/pkg/core"
	"offi/pkg/etf2l"
	"offi/pkg/handler"
	"offi/pkg/logstf"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	if _, ok := os.LookupEnv("DEBUG"); ok {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logsTf := logstf.New()
	etf2lClient := etf2l.New()

	redisUrl := os.Getenv("REDIS_URL")
	cacheClient, err := cache.New(redisUrl)
	if err != nil {
		log.Fatal(err)
	}
	c := core.New(cacheClient, etf2lClient, logsTf)

	app := handler.CreateApp(c)

	if err = app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
