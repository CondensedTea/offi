package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Match struct {
	MatchID     int
	Competition string
	Stage       string
	Tier        string
	CompletedAt time.Time
}

func (c *Client) MatchExists(ctx context.Context, mathcID int) (bool, error) {
	const query = `select exists(select 1 from matches where match_id = $1)`

	var exists bool
	if err := c.pool.QueryRow(ctx, query, mathcID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (c *Client) SaveMatch(ctx context.Context, tx pgx.Tx, match Match) error {
	const query = `
		insert into matches(match_id, competition, stage, tier, completed_at)
		values($1, $2, $3, $4, $5)
		on conflict do nothing`

	_, err := tx.Exec(ctx, query, match.MatchID, match.Competition, match.Stage, match.Tier, match.CompletedAt)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	return nil
}
