package main

import (
	"context"
	"errors"
	"fmt"
	"offi/internal/db"
	"offi/internal/demostf"
	"offi/internal/etf2l"
	internalHTTP "offi/internal/http"
	"offi/internal/logstf"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v3"
)

var linkCommand = &cli.Command{
	Name:      "link",
	UsageText: "link <match_id> [log_id1, log_id2, ...]",
	Usage:     "Fetches logs and links them to given match",
	Action:    linkAction,
}

func linkAction(ctx context.Context, cmd *cli.Command) error {
	dbClient, err := db.NewClient(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("initing db client: %w", err)
	}

	transport := internalHTTP.Transport(false)

	logstfClient := logstf.NewClient(transport)
	demosTfClient := demostf.NewClient(transport)
	etf2lClient := etf2l.New(transport)

	matchID, err := strconv.Atoi(cmd.Args().First())
	if err != nil {
		return fmt.Errorf("parsing match id: %w", err)
	}

	tx, err := dbClient.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginig transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, logID := range cmd.Args().Slice()[1:] {
		id, err := strconv.Atoi(logID)
		if err != nil {
			return fmt.Errorf("parsing log id: %w", err)
		}

		log, err := logstfClient.GetLog(ctx, id)
		if err != nil {
			return fmt.Errorf("fetching log %q: %w", logID, err)
		}

		if err = dbClient.SaveLog(ctx, tx, db.Log{
			LogID:    id,
			MatchID:  matchID,
			Title:    log.Title,
			Map:      log.Map,
			PlayedAt: time.Unix(log.Date, 0),
		}); err != nil {
			return fmt.Errorf("saving log: %w", err)
		}

		match, err := etf2lClient.GetMatch(ctx, matchID)
		if err != nil {
			return fmt.Errorf("fetching match %d: %w", matchID, err)
		}

		demo, err := demosTfClient.FindDemo(ctx, demostf.FindDemoRequest{
			PlayerSteamIDs: match.PlayerSteamIDs,
			Map:            log.Map,
			PlayedAt:       time.Unix(log.Date, 0),
		})
		if err != nil {
			if errors.Is(err, demostf.ErrNotFound) {
				fmt.Printf("no demo found for log %d\n", id)
				continue
			}
			return fmt.Errorf("finding demo: %w", err)
		}

		if err = dbClient.UpdateDemoIDForLogTx(ctx, tx, id, demo.ID); err != nil {
			return fmt.Errorf("updating demo id for log %d: %w", id, err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	fmt.Print("ok")

	return nil
}
