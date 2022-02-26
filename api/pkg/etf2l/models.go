package etf2l

type Meta struct {
	PreviousPageUrl string `json:"previous_page_url"`
	NextPageUrl     string `json:"next_page_url"`
	Page            int    `json:"page"`
	TotalPages      int    `json:"total_pages"`
	EntriesPerPage  int    `json:"entries_per_page"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SteamInfo struct {
	Avatar string  `json:"avatar"`
	Group  *string `json:"group"`
}

type Clan struct {
	Url     string    `json:"url"`
	Country string    `json:"country"`
	Drop    int       `json:"drop"`
	Id      int       `json:"id"`
	Steam   SteamInfo `json:"steam"`
	Name    string    `json:"name"`
}

type Division struct {
	Name         string `json:"name"`
	Tier         int    `json:"tier"`
	Id           int    `json:"id"`
	SkillContrib int    `json:"skill_contrib"`
}

type Result struct {
	Round      string   `json:"round"`
	R1         int      `json:"r1"`
	Maps       []string `json:"maps"`
	Clan2      Clan     `json:"clan2"`
	Clan1      Clan     `json:"clan1"`
	Week       int      `json:"week"`
	DefaultWin int      `json:"defaultwin"` // bool
	Time       int      `json:"time"`
	Division   Division `json:"division"`
	R2         int      `json:"r2"`
	Id         int      `json:"id"`
}

type CompetitionResults struct {
	Status  Status   `json:"status"`
	Page    Meta     `json:"page"`
	Results []Result `json:"results"`
}
