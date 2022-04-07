package main

import (
	"fmt"
	"log"
	"offi/pkg/cache"
	"offi/pkg/core"
	"offi/pkg/etf2l"
	"offi/pkg/logstf"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

var CLI struct {
	UnsetMatch struct {
		ID int `arg:"" name:"match_id"`
	} `cmd:"" help:"deletes single log set"`
	PlayerQuery struct {
		ID int `arg:"" name:"match_id"`
	} `cmd:"" help:"builds logs.tf API URL for given match"`
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
		fmt.Printf("https://logs.tf/api/v1/log?%s\n", query)
	default:
		panic(ctx.Command())
	}
}
