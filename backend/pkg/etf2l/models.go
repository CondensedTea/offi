package etf2l

import (
	"offi/pkg/cache"
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
	Start  int    `json:"start"`
	End    int    `json:"end"`
	Reason string `json:"reason"`
}

type Player struct {
	ID   int   `json:"id"`
	Bans []Ban `json:"bans"`
}

type PlayerResponse struct {
	Player Player `json:"player"`
	Status Status `json:"status"`
}
