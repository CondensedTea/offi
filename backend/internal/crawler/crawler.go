package crawler

import (
	"context"
	"fmt"
	"log/slog"
	"offi/internal/db"
	"offi/internal/etf2l"
	"time"
)

type database interface {
	GetLasRecruitmentID(ctx context.Context, postType db.Post) (int, error)
	SaveRecruitments(ctx context.Context, recs []db.Recruitment) error
	CleanupOldRecruitments(ctx context.Context, postType db.Post) (int64, error)
}

type Crawler struct {
	etf2l *etf2l.Client
	db    database
}

func NewCrawler(etf2l *etf2l.Client, db database) *Crawler {
	return &Crawler{
		etf2l: etf2l,
		db:    db,
	}
}

func (c *Crawler) CrawlTeamRecruitments() error {
	err := c.crawlTeamRecruitments()
	if err != nil {
		slog.Error("crawling team recruitments", "error", err)
	}

	return err
}

func (c *Crawler) crawlTeamRecruitments() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	lastID, err := c.db.GetLasRecruitmentID(ctx, db.Team)
	if err != nil {
		return fmt.Errorf("getting last team recruitment ID: %w", err)
	}

	posts, err := c.etf2l.LoadRecruitmentPosts(ctx, etf2l.TeamPost, lastID)
	if err != nil {
		return fmt.Errorf("loading team recruitments: %w", err)
	}

	var recs = make([]db.Recruitment, len(posts))
	for i, p := range posts {
		recruitmentID, err := p.RecruitmentID()
		if err != nil {
			return fmt.Errorf("parsing recrutiment ID: %w", err)
		}

		teamID, err := p.AuthorID(etf2l.TeamPost)
		if err != nil {
			return fmt.Errorf("parsing team ID: %w", err)
		}

		recs[i] = db.Recruitment{
			RecruitmentID: recruitmentID,
			AuthorID:      teamID,
			PostType:      db.Team,
			TeamType:      p.Type,
			Classes:       p.Classes,
			SkillLevel:    p.Skill,
		}
	}

	if err = c.db.SaveRecruitments(ctx, recs); err != nil {
		return fmt.Errorf("savin team recruitments: %w", err)
	}

	deleted, err := c.db.CleanupOldRecruitments(ctx, db.Team)
	if err != nil {
		return fmt.Errorf("cleaning up old team recruitments: %w", err)
	}

	slog.Info("loaded team recruitments", "new_count", len(recs), "deleted_count", deleted)

	return nil
}

func (c *Crawler) CrawlPlayerRecruitments() error {
	err := c.crawlPlayerRecruitments()
	if err != nil {
		slog.Error("crawling player recruitments", "error", err)
	}

	return err
}

func (c *Crawler) crawlPlayerRecruitments() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	lastID, err := c.db.GetLasRecruitmentID(ctx, db.Player)
	if err != nil {
		return fmt.Errorf("getting last player recruitment ID: %w", err)
	}

	posts, err := c.etf2l.LoadRecruitmentPosts(ctx, etf2l.PlayerPost, lastID)
	if err != nil {
		return fmt.Errorf("loading player recruitments: %w", err)
	}

	var recs = make([]db.Recruitment, len(posts))
	for i, p := range posts {
		recruitmentID, err := p.RecruitmentID()
		if err != nil {
			return fmt.Errorf("parsing recruitment ID: %w", err)
		}

		playerID, err := p.AuthorID(etf2l.PlayerPost)
		if err != nil {
			return fmt.Errorf("parsing author ID: %w", err)
		}

		recs[i] = db.Recruitment{
			RecruitmentID: recruitmentID,
			AuthorID:      playerID,
			PostType:      db.Player,
			TeamType:      p.Type,
			Classes:       p.Classes,
			SkillLevel:    p.Skill,
		}
	}

	if err = c.db.SaveRecruitments(ctx, recs); err != nil {
		return fmt.Errorf("savin player recruitments: %w", err)
	}

	deleted, err := c.db.CleanupOldRecruitments(ctx, db.Player)
	if err != nil {
		return fmt.Errorf("cleaning up old team recruitments: %w", err)
	}

	slog.Info("loaded player recruitments", "new_count", len(recs), "deleted_count", deleted)

	return nil
}
