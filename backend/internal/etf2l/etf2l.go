package etf2l

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/time/rate"
)

type Client struct {
	apiURL     string
	httpClient *http.Client
	limiter    *rate.Limiter
	tracer     trace.Tracer
}

func New(rt http.RoundTripper) *Client {
	return &Client{
		apiURL:     "https://api-v2.etf2l.org",
		httpClient: &http.Client{Transport: rt},
		limiter:    rate.NewLimiter(rate.Every(time.Second), 5),
		tracer:     otel.Tracer("etf2l"),
	}
}
