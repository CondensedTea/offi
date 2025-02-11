package db

import (
	"context"
	"fmt"
	"time"
)

type Post uint

const (
	Unknown Post = iota
	Team
	Player
)

type Recruitment struct {
	RecruitmentID int       `json:"recruitment_id"`
	AuthorID      int       `json:"author_id"`
	PostType      Post      `json:"post_type"`
	TeamType      string    `json:"team_type"`
	Classes       []string  `json:"classes"`
	SkillLevel    string    `json:"skill_level"`
	CreatedAt     time.Time `json:"created_at"`
}

func (c *Client) GetLasRecruitmentID(ctx context.Context, postType Post) (int, error) {
	const query = `select max(recruitment_id) from recruitments where post_type = $1`

	var id int
	if err := c.pool.QueryRow(ctx, query, postType.String()).Scan(&id); err != nil {
		return 0, fmt.Errorf("running query: %w", err)
	}

	return id, nil
}

func (c *Client) SaveRecruitments(ctx context.Context, recs []Recruitment) error {
	const query = `
		insert into recruitments(recruitment_id, author_id, post_type, team_type, classes, skill_level)
		select recruitment_id, author_id, post_type, team_type, classes, skill_level from json_to_recordset($1::json) as t(recruitment_id bigint, author_id bigint, post_type text, team_type text, classes text[], skill_level text)`

	if _, err := c.pool.Exec(ctx, query, recs); err != nil {
		return fmt.Errorf("running exec: %w", err)
	}

	return nil
}

func (c *Client) CleanupOldRecruitments(ctx context.Context, postType Post) (int64, error) {
	const query = `delete from recruitments where now() - created_at > interval '2 weeks' and post_type = $1`

	res, err := c.pool.Exec(ctx, query, postType.String())
	if err != nil {
		return 0, fmt.Errorf("running exec: %w", err)
	}

	return res.RowsAffected(), nil
}

func (c *Client) GetLastRecruitmentForAuthor(ctx context.Context, postType Post, authorID int) (Recruitment, error) {
	const query = `
		select
			recruitment_id,
			author_id,
			post_type,
			team_type,
			skill_level,
			created_at
		from recruitments 
		where post_type = $1 and
		      author_id = $2 and
		      created_at > now() - interval '4 week'
	  	order by recruitment_id desc
		limit 1`

	var r Recruitment
	if err := c.pool.QueryRow(ctx, query, postType.String(), authorID).Scan(&r); err != nil {
		return Recruitment{}, fmt.Errorf("running query: %w", err)
	}

	return r, nil
}
