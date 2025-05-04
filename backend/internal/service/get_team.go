package service

import (
	"context"
	"errors"
	"fmt"
	"offi/internal/db"
	"offi/internal/gen/api"
	"unsafe"

	"github.com/jackc/pgx/v5"
)

func (s *Service) GetTeam(ctx context.Context, p api.GetTeamParams) (api.GetTeamRes, error) {
	recruitment, err := s.db.GetLastRecruitmentForAuthor(ctx, db.Team, int64(p.ID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &api.GetTeamNotFound{}, nil
		}

		return nil, fmt.Errorf("getting recruitments from db: %w", err)
	}

	return &api.GetTeamOK{
		Team: api.Team{
			Recruitment: api.RecruitmentInfo{
				Skill:    recruitment.SkillLevel,
				URL:      fmt.Sprintf("https://etf2l.org/recruitment/%d/", recruitment.RecruitmentID),
				Classes:  *(*[]api.GameClass)(unsafe.Pointer(&recruitment.Classes)),
				GameMode: recruitment.TeamType,
			},
		},
	}, nil
}
