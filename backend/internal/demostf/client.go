package demostf

import (
	"net/http"
	info "offi/internal/build_info"
	"offi/internal/tracing"

	"github.com/go-chi/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	client *http.Client
	tracer trace.Tracer
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{
			Transport: transport.Chain(
				http.DefaultTransport,
				tracing.OTelHTTPTransport,
				transport.SetHeader("User-Agent", "offi-backend/"+info.Version),
			),
		},
		tracer: otel.Tracer("demostf"),
	}
}
