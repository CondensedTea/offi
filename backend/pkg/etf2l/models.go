package etf2l

import (
	"encoding/json"
	"offi/pkg/cache"
	"time"
)

type Page struct {
	EntriesPerPage  int    `json:"entries_per_page"`
	NextPageUrl     string `json:"next_page_url"`
	PreviousPageUrl string `json:"previous_page_url"`
	Page            int    `json:"page"`
	TotalPages      int    `json:"total_pages"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Steam struct {
	Avatar string `json:"avatar"`
	Id     string `json:"id,omitempty"`
	Id3    string `json:"id3,omitempty"`
	Id64   string `json:"id64,omitempty"`
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
	Classes  []string `json:"classes"`
	Comments Comments `json:"comments"`
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Skill    string   `json:"skill"`
	Steam    Steam    `json:"steam"`
	Type     string   `json:"type"`
	Urls     URLs     `json:"urls"`
}

func (r Recruitment) ToCache() cache.Entry {
	return cache.Entry{
		ID:       r.Id,
		Skill:    r.Skill,
		URL:      r.Urls.Recruitment,
		GameMode: r.Type,
		Classes:  r.Classes,
	}
}

type RecruitmentResponse struct {
	Page         Page          `json:"page"`
	Recruitments []Recruitment `json:"recruitment"`
	Status       Status        `json:"status"`
}

type Ban struct {
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Reason string    `json:"reason"`
}

func (b *Ban) UnmarshalJSON(data []byte) error {
	type rawBan struct {
		Start  int
		End    int
		Reason string
	}

	var banData rawBan
	if err := json.Unmarshal(data, &banData); err != nil {
		return err
	}

	b.Start = time.Unix(int64(banData.Start), 0)
	b.End = time.Unix(int64(banData.End), 0)
	b.Reason = banData.Reason
	return nil
}

type Player struct {
	ID   int   `json:"id"`
	Bans []Ban `json:"bans"`
}

type PlayerResponse struct {
	Player Player `json:"player"`
	Status Status `json:"status"`
}
