package etf2l

import (
	"errors"
	"fmt"
	"offi/internal/redis"
)

type Steam struct {
	ID64 string `json:"id64"`
}

type URLs struct {
	Player      string `json:"player,omitempty"`
	Team        string `json:"team,omitempty"`
	Recruitment string `json:"recruitment"`
}

type Comments struct {
	Count int `json:"count"`
	Last  int `json:"last"`
}

type Recruitment struct {
	ID       int      `json:"id"`
	Classes  []string `json:"classes"`
	Comments Comments `json:"comments"`
	Name     string   `json:"name"`
	Skill    string   `json:"skill"`
	Steam    Steam    `json:"steam"`
	Type     string   `json:"type"`
	URLs     URLs     `json:"urls"`
}

func (r Recruitment) RecruitmentID() (int, error) {
	if r.URLs.Recruitment == "" {
		return 0, errors.New("recruitment does not have URL")
	}

	var id int
	_, err := fmt.Sscanf(r.URLs.Recruitment, "https://etf2l.org/recruitment/%d", &id)
	if err != nil {
		return 0, fmt.Errorf("recruitment %q: parsing player ID: %w", r.URLs.Recruitment, err)
	}

	return id, nil
}

func (r Recruitment) AuthorID(t postType) (int, error) {
	var (
		rawURL     string
		urlPattern string
	)
	switch t {
	case PlayerPost:
		rawURL = r.URLs.Player
		urlPattern = "http://api-v2.etf2l.org/player/%d"
	case TeamPost:
		rawURL = r.URLs.Team
		urlPattern = "http://api-v2.etf2l.org/team/%d"
	default:
		return 0, fmt.Errorf("unknown post type %q", t)
	}

	var id int
	_, err := fmt.Sscanf(rawURL, urlPattern, &id)
	if err != nil {
		return 0, fmt.Errorf("recruitment %q: %w", r.URLs.Recruitment, err)
	}

	return id, nil
}

type RecruitmentResponse struct {
	Recruitments struct {
		NextPageURL string        `json:"next_page_url"`
		Data        []Recruitment `json:"data"`
	} `json:"recruitment"`
}

type Ban struct {
	Start  int    `json:"start"`
	End    int    `json:"end"`
	Reason string `json:"reason"`
}

type Player struct {
	ID    int `json:"id"`
	Steam struct {
		ID64 int64 `json:"id64,string"`
	} `json:"steam"`
	Name string `json:"name"`
	Bans []Ban  `json:"bans"`
}

type PlayerResponse struct {
	Player Player `json:"player"`
}

// ToCache converts a Player to a cache.Player
func (p Player) ToCache() redis.Player {
	var bans = make([]redis.PlayerBan, len(p.Bans))
	for i, v := range p.Bans {
		bans[i] = redis.PlayerBan{
			Start:  v.Start,
			End:    v.End,
			Reason: v.Reason,
		}
	}

	return redis.Player{
		ID:      p.ID,
		Bans:    bans,
		SteamID: p.Steam.ID64,
		Name:    p.Name,
	}
}
