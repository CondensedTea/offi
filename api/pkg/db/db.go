package db

import "github.com/jackc/pgx"

type Database struct {
	pool *pgx.ConnPool
}

func New(dsn string) (*Database, error) {
	c, err := pgx.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	cfg := pgx.ConnPoolConfig{
		ConnConfig: c,
	}

	pool, err := pgx.NewConnPool(cfg)
	if err != nil {
		return nil, err
	}
	return &Database{pool: pool}, nil
}

func (db Database) name() {

}
