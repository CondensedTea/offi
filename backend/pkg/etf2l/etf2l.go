package etf2l

import (
	"net/http"
	"time"
)

const etf2lApiURL = "https://api.etf2l.org/"

type ETF2L struct {
	httpClient *http.Client
}

func New() (*ETF2L, error) {
	return &ETF2L{
		httpClient: &http.Client{Timeout: 5 * time.Minute},
	}, nil
}
