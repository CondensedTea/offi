package http

import (
	"log/slog"
	"math"
	"net/http"
	info "offi/internal/build_info"
	"offi/internal/tracing"
	"strconv"
	"time"

	"github.com/go-chi/transport"
)

func Transport(withRetries bool) http.RoundTripper {
	uaTransport := transport.SetHeader("User-Agent", "offi-backend/"+info.Version)

	base := []func(http.RoundTripper) http.RoundTripper{
		tracing.OTelHTTPTransport,
		uaTransport,
	}

	if withRetries {
		base = append(base, Retry(transport.Chain(http.DefaultTransport, uaTransport), 3))
	}

	return transport.Chain(http.DefaultTransport, base...)
}

func Retry(baseTransport http.RoundTripper, maxRetries int) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return transport.RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
			defer func() {
				if isRetryable(resp) {
					for i := 1; i <= maxRetries; i++ {
						wait := backOff(resp, i)

						timer := time.NewTimer(wait)

						slog.DebugContext(req.Context(), "waiting before retrying request", "wait_time", wait.String())

						select {
						case <-req.Context().Done():
							timer.Stop()
							err = req.Context().Err()
							return
						case <-timer.C:
						}

						resp, err = baseTransport.RoundTrip(req)
						if !isRetryable(resp) {
							break
						}

						slog.InfoContext(req.Context(), "retrying request",
							"target_host", req.URL.Host,
							"attempt", i,
						)
					}
				}
			}()

			return next.RoundTrip(req)
		})
	}
}

func backOff(resp *http.Response, attempt int) time.Duration {
	minDuration := 1 * time.Second
	maxDuration := 16 * time.Second

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
		if s, ok := resp.Header["Retry-After"]; ok {
			if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
				return time.Second * time.Duration(sleep)
			}
		}
	}

	// simple exp. backoff
	mult := math.Pow(2, float64(attempt)) * float64(minDuration)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > maxDuration {
		sleep = maxDuration
	}
	return sleep
}

func isRetryable(resp *http.Response) bool {
	if resp == nil {
		return false
	}

	// 429 Too Many Requests is recoverable. Sometimes the server puts
	// Retry-After response header to indicate when the server is will be available again
	if resp.StatusCode == http.StatusTooManyRequests {
		return true
	}

	// We retry on 500-range responses to allow the server time to recover, as 500's are typically not permanent
	// errors and may relate to outages on the server side.
	if resp.StatusCode == 0 || (resp.StatusCode > 500 && resp.StatusCode != http.StatusNotImplemented) {
		return true
	}

	return false
}
