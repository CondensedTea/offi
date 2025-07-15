package etf2l

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrMatchNotFound   = errors.New("match does not exist")
	ErrIncompleteMatch = errors.New("match does not contain all required data")
)

type Competition struct {
	Category string `json:"category"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

type Division struct {
	ID           int    `json:"id"`
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
	ID          int         `json:"id"`
	Maps        []string    `json:"maps"`
	Round       string      `json:"round"`
	ScheduledAt int         `json:"time"`
	SubmittedAt int         `json:"submitted"`
	Week        int         `json:"week"`
	Players     []Player    `json:"players"`
	ByeWeek     bool        `json:"bye_week"`
	MapResults  []MapResult `json:"map_results"`
}

type Match struct {
	PlayerSteamIDs []int64
	Maps           []string
	SubmittedAt    time.Time

	ID          int
	Competition string
	Tier        string
	Stage       string
}

type MatchResponse struct {
	Match match `json:"match"`
}

func (c Client) GetMatch(ctx context.Context, id int) (*Match, error) {
	ctx, span := c.tracer.Start(ctx, "etf2l.GetMatch",
		trace.WithAttributes(attribute.Int("match_id", id)),
	)
	defer span.End()

	if err := c.limiter.Wait(ctx); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/matches/%d", c.apiURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get match from etf2l api: %w", err)
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

	if matchResponse.Match.Defaultwin || len(matchResponse.Match.Players) == 0 || matchResponse.Match.SubmittedAt == 0 {
		// Default win, match without players or match will be played in the future
		return nil, ErrIncompleteMatch
	}

	// etf2l returns duplicate players in the match response
	var playerIDSet = make(map[int64]struct{})
	for _, player := range matchResponse.Match.Players {
		if player.Steam.ID64 == 0 {
			// Special case for matches restored after Great Data loss of 2020 (?)
			return nil, ErrIncompleteMatch
		}

		playerIDSet[player.Steam.ID64] = struct{}{}
	}

	var playerIDs = make([]int64, 0, len(playerIDSet))
	for k := range playerIDSet {
		playerIDs = append(playerIDs, k)
	}

	return &Match{
		ID:             matchResponse.Match.ID,
		PlayerSteamIDs: playerIDs,
		Maps:           matchResponse.Match.Maps,
		SubmittedAt:    time.Unix(int64(matchResponse.Match.SubmittedAt), 0),
		Competition:    matchResponse.Match.Competition.Name,
		Tier:           matchResponse.Match.Division.Name,
		Stage:          matchResponse.Match.Round,
	}, nil
}
