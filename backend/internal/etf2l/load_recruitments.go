package etf2l

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type postType string

const (
	PlayerPost = "player"
	TeamPost   = "team"
)

func (c Client) LoadRecruitmentPosts(postType postType) ([]Recruitment, error) {
	var (
		entries []Recruitment
		url     string
	)

	switch postType {
	case PlayerPost:
		url = "https://api-v2.etf2l.org/recruitment/players?per_page=100"
	case TeamPost:
		url = "https://api.etf2l.org/recruitment/teams?per_page=100"
	}

	pageIsNotLast := true

	for pageIsNotLast {
		resp, err := c.httpClient.Get(url)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("etf2l api returned bad status: %d", resp.StatusCode)
		}

		var response RecruitmentResponse
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}
		if err = resp.Body.Close(); err != nil {
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
