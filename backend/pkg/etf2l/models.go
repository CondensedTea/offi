package etf2l

import "offi/pkg/cache"

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

type Entry struct {
	Classes  []string `json:"classes"`
	Comments Comments `json:"comments"`
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Skill    string   `json:"skill"`
	Steam    Steam    `json:"steam"`
	Type     string   `json:"type"`
	Urls     URLs     `json:"urls"`
}

func (e Entry) ToCache() cache.Entry {
	return cache.Entry{
		ID:       e.Id,
		Skill:    e.Skill,
		URL:      e.Urls.Recruitment,
		GameMode: e.Type,
		Classes:  e.Classes,
	}
}

type RecruitmentResponse struct {
	Page        Page    `json:"page"`
	Recruitment []Entry `json:"recruitment"`
	Status      Status  `json:"status"`
}
