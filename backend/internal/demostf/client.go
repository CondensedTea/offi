package demostf

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	client *http.Client
	tracer trace.Tracer
}

func NewClient(rt http.RoundTripper) *Client {
	return &Client{
		client: &http.Client{Transport: rt},
		tracer: otel.Tracer("demostf"),
	}
}
