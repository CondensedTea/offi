package http

import (
	"net/http"
	info "offi/internal/build_info"
	"offi/internal/tracing"

	"github.com/go-chi/transport"
)

func Transport(withRetries bool) http.RoundTripper {
	uaTransport := transport.SetHeader("User-Agent", "offi-backend/"+info.Version)

	base := []func(http.RoundTripper) http.RoundTripper{
		tracing.OTelHTTPTransport,
		uaTransport,
	}

	if withRetries {
		base = append(base, transport.Retry(transport.Chain(http.DefaultTransport, uaTransport), 3))
	}

	return transport.Chain(http.DefaultTransport, base...)
}
