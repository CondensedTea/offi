package main

import (
	"log"
	"offi/pkg/cache"
	"offi/pkg/core"
	"offi/pkg/etf2l"
	"offi/pkg/logstf"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"
)

var CLI struct {
	UnsetMatch struct {
		ID int `arg:"" name:"match_id"`
	} `cmd:"" help:"deletes single log set"`
	PlayerQuery struct {
		ID int `arg:"" name:"match_id"`
	} `cmd:"" help:"builds logs.tf API URL for given match"`
	LinkLog struct {
		Secondary bool `help:"mark log as secondary"`
		MatchID   int  `arg:"" name:"match_id"`
		LogID     int  `arg:"" name:"log_id"`
	} `cmd:"" help:"links log to match"`
}

func main() {
	logsTf := logstf.New()
	etf2lClient := etf2l.New()

	redisUrl := os.Getenv("REDIS_URL")
	cacheClient, err := cache.New(redisUrl)
	if err != nil {
		log.Fatal(err)
	}
	c := core.New(cacheClient, etf2lClient, logsTf)

	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "unset-match <match_id>":
		// logSet, err := cacheClient.DeleteLogs(CLI.UnsetMatch.ID)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// log.Printf("%v\n", logSet)
	case "player-query <match_id>":
		match, err := etf2lClient.ParseMatchPage(CLI.PlayerQuery.ID)
		if err != nil {
			log.Fatal(err)
		}
		steamIds, err := c.GetSteamIDs(match)
		if err != nil {
			log.Fatal(err)
		}
		query := "player=" + strings.Join(steamIds, ",")
		logrus.Infof("https://logs.tf/api/v1/log?%s", query)
	case "link-log <match_id> <log_id>":
		l, err := logsTf.GetLog(CLI.LinkLog.LogID)
		if err != nil {
			log.Fatal(err)
		}
		cacheLog := l.ToCache(CLI.LinkLog.Secondary)

		logSet, err := cacheClient.GetLogs(CLI.LinkLog.MatchID)
		if err != nil {
			log.Fatal(err)
		}
		(&logSet).Logs = append(logSet.Logs, cacheLog)
		if err = cacheClient.SetLogs(CLI.LinkLog.MatchID, &logSet); err != nil {
			log.Fatal(err)
		}
	default:
		panic(ctx.Command())
	}
}
