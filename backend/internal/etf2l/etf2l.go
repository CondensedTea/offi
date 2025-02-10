package etf2l

import (
	"net/http"
	info "offi/internal/build_info"
	"time"

	"github.com/go-chi/transport"
	"golang.org/x/time/rate"
)

type Client struct {
	apiURL     string
	httpClient *http.Client
	limiter    *rate.Limiter
}

func New() *Client {
	return &Client{
		apiURL: "https://api-v2.etf2l.org",
		httpClient: &http.Client{
			Transport: transport.Chain(
				http.DefaultTransport,
				transport.SetHeader("User-Agent", "offi-backend/"+info.Version),
			),
			Timeout: 5 * time.Second,
		},
		limiter: rate.NewLimiter(rate.Every(time.Second), 5),
	}
}
