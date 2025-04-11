package etf2l

import (
	"net/http"
	info "offi/internal/build_info"
	"offi/internal/tracing"
	"time"

	"github.com/go-chi/transport"
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

func New() *Client {
	return &Client{
		apiURL: "https://api-v2.etf2l.org",
		httpClient: &http.Client{
			Transport: transport.Chain(
				http.DefaultTransport,
				tracing.OTelHTTPTransport,
				transport.SetHeader("User-Agent", "offi-backend/"+info.Version),
			),
			Timeout: 5 * time.Second,
		},
		limiter: rate.NewLimiter(rate.Every(time.Second), 6),
		tracer:  otel.Tracer("etf2l"),
	}
}
