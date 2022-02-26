package logstf

import (
	"net/http"
	"time"
)

type Logs struct {
	httpClient *http.Client
}

func New() (*Logs, error) {
	return &Logs{
		httpClient: &http.Client{Timeout: 5 * time.Minute},
	}, nil
}
