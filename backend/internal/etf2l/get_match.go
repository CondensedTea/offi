package etf2l

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrMatchNotFound = errors.New("match does not exist")
)

type Clan struct {
	Country string `json:"country"`
	Drop    bool   `json:"drop"`
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Steam   struct {
		Avatar string      `json:"avatar"`
		Group  interface{} `json:"group"`
	} `json:"steam"`
	Url string `json:"url"`
}

type Competition struct {
	Category string `json:"category"`
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Url      string `json:"url"`
}

type Division struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	SkillContrib int    `json:"skill_contrib"`
	Tier         int    `json:"tier"`
}

type MapResult struct {
	MatchOrder int    `json:"match_order"`
	Clan1      int    `json:"clan1"`
	Clan2      int    `json:"clan2"`
	Map        string `json:"map"`
}

type match struct {
	Competition Competition `json:"competition"`
	Defaultwin  bool        `json:"defaultwin"`
	Division    Division    `json:"division"`
	Id          int         `json:"id"`
	Maps        []string    `json:"maps"`
	Round       string      `json:"round"`
	Time        int         `json:"time"`
	Submitted   *int        `json:"submitted"`
	Week        int         `json:"week"`
	Players     []Player    `json:"players"`
	ByeWeek     bool        `json:"bye_week"`
	MapResults  []MapResult `json:"map_results"`
}

type Match struct {
	Players  []string
	Maps     []string
	PlayedAt time.Time

	ID          int
	Competition string
	Tier        string
	Stage       string
}

type MatchResponse struct {
	Match match `json:"match"`
}

func (c Client) GetMatch(ctx context.Context, id int) (*Match, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/matches/%d", c.apiURL, id)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get match from etf2l api: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrMatchNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("etf2l api returned non-200 status: %d", resp.StatusCode)
	}

	var matchResponse MatchResponse
	if err = json.NewDecoder(resp.Body).Decode(&matchResponse); err != nil {
		return nil, err
	}

	if matchResponse.Match.Defaultwin || len(matchResponse.Match.Players) == 0 || matchResponse.Match.Submitted == nil {
		// Default win, match without players or match will be played in the future
		return nil, ErrMatchNotFound
	}

	// etf2l returns duplicate players in the match response
	var playerIDSet = make(map[string]struct{})
	for _, player := range matchResponse.Match.Players {
		playerIDSet[player.Steam.ID64] = struct{}{}
	}

	var playerIDs = make([]string, 0, len(playerIDSet))
	for k := range playerIDSet {
		playerIDs = append(playerIDs, k)
	}

	playedAt := time.Unix(int64(*matchResponse.Match.Submitted), 0)

	return &Match{
		ID:          matchResponse.Match.Id,
		Players:     playerIDs,
		Maps:        matchResponse.Match.Maps,
		PlayedAt:    playedAt,
		Competition: matchResponse.Match.Competition.Name,
		Tier:        matchResponse.Match.Division.Name,
		Stage:       matchResponse.Match.Round,
	}, nil
}
