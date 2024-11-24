package etf2l

import (
	"net/http"
	"time"
)

type Getter interface {
	Get(url string) (*http.Response, error)
}

type Client struct {
	httpClient Getter
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}
