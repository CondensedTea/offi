package service

import (
	"context"
	"offi/internal/gen/api"

	ht "github.com/ogen-go/ogen/http"
)

func (s *Service) GetTeam(_ context.Context, _ api.GetTeamParams) (api.GetTeamRes, error) {
	// recruitments, err := s.db.GetRecruitments(ctx, db.Team, p.ID)
	// if err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return &api.GetTeamNotFound{}, nil
	// 	}
	//
	// 	return nil, fmt.Errorf("getting recruitments from db: %w", err)
	// }
	//
	// res := make([]*api.RecruitmentInfo, len(recruitments))
	// for i, r := range recruitments {
	// 	res[i] = &api.RecruitmentInfo{
	// 		Skill:    r.SkillLevel,
	// 		URL:      fmt.Sprintf("https://etf2l.org/recruitment/%d/", r.RecruitmentID),
	// 		GameMode: r.TeamType,
	// 	}
	// }

	return nil, ht.ErrNotImplemented
}
