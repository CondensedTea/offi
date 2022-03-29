package etf2l

import (
	"net/http"
	"time"
)

type Getter interface {
	Get(url string) (*http.Response, error)
}

type ETF2L struct {
	httpClient Getter
}

func New() *ETF2L {
	return &ETF2L{
		httpClient: &http.Client{Timeout: 5 * time.Minute},
	}
}
