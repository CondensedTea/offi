package etf2l

import (
	"context"
	"net/http"

	"github.com/carlmjohnson/requests"
)

type postType string

const (
	PlayerPost = "player"
	TeamPost   = "team"
)

func LoadRecruitmentPosts(ctx context.Context, postType postType) ([]Recruitment, error) {
	var (
		entries []Recruitment
		url     string
	)

	switch postType {
	case PlayerPost:
		url = "https://api.etf2l.org/recruitment/players.json?per_page=100"
	case TeamPost:
		url = "https://api.etf2l.org/recruitment/teams.json?per_page=100"
	}

	pageIsNotLast := true

	for pageIsNotLast {
		var response RecruitmentResponse
		err := requests.
			URL(url).
			ToJSON(&response).
			CheckStatus(http.StatusOK).
			Fetch(ctx)
		if err != nil {
			return nil, err
		}

		if response.Page.NextPageUrl == "" {
			pageIsNotLast = false
		} else {
			url = response.Page.NextPageUrl
		}
		entries = append(entries, response.Recruitments...)
	}
	return entries, nil
}
