package logstf

type Response struct {
	Logs []Log `json:"logs"`
}

type Log struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Map     string `json:"map"`
	Date    int    `json:"date"`
	Views   int    `json:"views"`
	Players int    `json:"players"`
}
