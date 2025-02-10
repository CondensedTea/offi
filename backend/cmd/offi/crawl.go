package main

import (
	"context"
	"fmt"
	"offi/internal/crawler"
	"offi/internal/db"
	"offi/internal/etf2l"
	"os"

	"github.com/urfave/cli/v3"
)

var crawlCommand = &cli.Command{
	Name:      "crawl",
	UsageText: "crawl <team|player>",
	Usage:     "loads the latest recruitment posts from etf2l to database",
	Action:    crawlAction,
}

// crawlAction is the action for the crawl command.
func crawlAction(ctx context.Context, cmd *cli.Command) error {
	dbClient, err := db.NewClient(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("failed to init db client: %w", err)
	}

	// etf2l client
	etf2lClient := etf2l.New()

	c := crawler.NewCrawler(etf2lClient, dbClient)

	switch cmd.Args().First() {
	case "team":
		err = c.CrawlTeamRecruitments()
	case "player":
		err = c.CrawlPlayerRecruitments()
	default:
		return fmt.Errorf("unknown recruitment type: %s", cmd.Args().First())
	}

	if err != nil {
		return fmt.Errorf("failed to crawl %q recruitments: %w", cmd.Args().First(), err)
	}

	return nil
}
