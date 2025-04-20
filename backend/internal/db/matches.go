package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type conn interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type Match struct {
	MatchID     int
	Competition string
	Stage       string
	Tier        string
	CompletedAt time.Time
	DemoID      sql.Null[int]
}

func (c *Client) MatchExists(ctx context.Context, mathcID int) (bool, error) {
	const query = `select exists(select 1 from matches where match_id = $1)`

	var exists bool
	if err := c.pool.QueryRow(ctx, query, mathcID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (c *Client) SaveMatchTx(ctx context.Context, tx pgx.Tx, match Match) error {
	return c.saveMatch(ctx, tx, match)
}

func (c *Client) SaveMatch(ctx context.Context, match Match) error {
	return c.saveMatch(ctx, c.pool, match)
}

func (c *Client) saveMatch(ctx context.Context, conn conn, match Match) error {
	const query = `
		insert into matches(match_id, competition, stage, tier, completed_at)
		values ($1, $2, $3, $4, $5)`

	_, err := conn.Exec(ctx, query, match.MatchID, match.Competition, match.Stage, match.Tier, match.CompletedAt)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	return nil
}

func (c *Client) GetMatchByLogID(ctx context.Context, logID int) (Match, error) {
	const query = `
		with cte as (
			select match_id, demo_id from logs where log_id = $1
		)
		select m.match_id, m.competition, m.stage, m.tier, m.completed_at, cte.demo_id
		from matches m
		join cte on m.match_id = cte.match_id`

	rows, err := c.pool.Query(ctx, query, logID)
	if err != nil {
		return Match{}, fmt.Errorf("querying rows: %w", err)
	}

	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[Match])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Match{}, ErrNotFound
		}

		return Match{}, fmt.Errorf("collecting rows: %w", err)
	}

	return res, nil
}
