package etf2l

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type postType string

const (
	// PlayerPost represents a recruitment post for ETF2L players.
	PlayerPost = "player"
	// TeamPost represents a recruitment post for ETF2L teams.
	TeamPost = "team"
)

func (c Client) LoadRecruitmentPosts(ctx context.Context, postType postType, lastID int) ([]Recruitment, error) {
	var (
		entries []Recruitment
		url     string
	)

	switch postType {
	case PlayerPost:
		url = c.apiURL + "/recruitment/players?limit=100"
	case TeamPost:
		url = c.apiURL + "/recruitment/teams?limit=100"
	}

	for {
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, err
		}

		resp, err := c.httpClient.Get(url)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("etf2l api returned status %d: %s", resp.StatusCode, string(b))
		}

		var response RecruitmentResponse
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}
		if err = resp.Body.Close(); err != nil {
			return nil, err
		}

		for _, recruitment := range response.Recruitments.Data {
			id, err := recruitment.RecruitmentID()
			if err != nil {
				return nil, fmt.Errorf("parsing recruitment ID: %w", err)
			}
			if id > lastID {
				entries = append(entries, recruitment)
			}
		}

		id, err := response.Recruitments.Data[len(response.Recruitments.Data)-1].RecruitmentID()
		if err != nil {
			return nil, fmt.Errorf("parsing last recruitment ID: %w", err)
		}

		if response.Recruitments.NextPageURL == "" || id <= lastID {
			break
		}

		url = response.Recruitments.NextPageURL
	}
	return entries, nil
}
