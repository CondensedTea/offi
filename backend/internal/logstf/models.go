package logstf

type Response struct {
	Logs []Log `json:"logs"`
}

type Log struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Map   string `json:"map"`
	Date  int64  `json:"date"`
}
