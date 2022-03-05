package db

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

type Competition struct {
	ID          int
	LastMatchID *int `db:"last_match_id"`
	isCompleted bool `db:"is_completed"`
}

func New(ctx context.Context, dsn string) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Database{pool: pool}, nil
}

func (db Database) GetCompetitions(ctx context.Context) ([]Competition, error) {
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	const query = `
		select 
			id, max(m.etf2l_match_id) as last_match_id
		from
			 competitions
		left outer join matches as m on competitions.id = m.competition_id
		where
			  is_completed = false
		group by
				 competitions.id
		`

	var c []Competition

	err = pgxscan.Select(ctx, conn, &c, query)
	if err == pgx.ErrNoRows {
		return c, nil
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}
