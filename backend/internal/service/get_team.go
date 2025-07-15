package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"offi/internal/db"
	"offi/internal/gen/api"
	"unsafe"

	"github.com/jackc/pgx/v5"
)

func (s *Service) GetTeam(ctx context.Context, p api.GetTeamParams) (api.GetTeamRes, error) {
	recruitment, err := s.db.GetLastRecruitmentForAuthor(ctx, db.Team, int64(p.ID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &api.ErrorStatusCode{StatusCode: http.StatusNotFound, Response: api.Error{Error: "team not found"}}, nil
		}

		return nil, fmt.Errorf("getting recruitments from db: %w", err)
	}

	return &api.GetTeamOK{
		Team: api.Team{
			Recruitment: api.RecruitmentInfo{
				Skill:    recruitment.SkillLevel,
				URL:      fmt.Sprintf("https://etf2l.org/recruitment/%d/", recruitment.RecruitmentID),
				Classes:  *(*[]api.GameClass)(unsafe.Pointer(&recruitment.Classes)), //nolint:gosec
				GameMode: recruitment.TeamType,
			},
		},
	}, nil
}
