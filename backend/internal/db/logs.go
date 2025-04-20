package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Log struct {
	LogID       int
	MatchID     int
	Title       string
	Map         string
	PlayedAt    time.Time
	IsSecondary bool
	DemoID      sql.Null[int]
}

func (c *Client) SaveLog(ctx context.Context, tx pgx.Tx, log Log) error {
	const query = `insert into logs(log_id, match_id, title, map, played_at, is_secondary) values($1, $2, $3, $4, $5, $6)`

	_, err := tx.Exec(
		ctx, query,
		log.LogID,
		log.MatchID,
		log.Title,
		log.Map,
		log.PlayedAt,
		log.IsSecondary,
	)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	return nil
}

func (c *Client) UpdateDemoIDForLog(ctx context.Context, logID int, demoID int) error {
	const query = `update logs set demo_id = $1 where log_id = $2`

	_, err := c.pool.Exec(ctx, query, demoID, logID)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	return nil
}

func (c *Client) GetLogsByMatchID(ctx context.Context, matchID int) ([]Log, error) {
	const query = `select log_id, match_id, title, map, played_at, is_secondary, demo_id from logs where match_id = $1`

	rows, err := c.pool.Query(ctx, query, matchID)
	if err != nil {
		return nil, fmt.Errorf("querying rows: %w", err)
	}

	res, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Log])
	if err != nil {
		return nil, fmt.Errorf("collecting rows: %w", err)
	}

	return res, nil
}
