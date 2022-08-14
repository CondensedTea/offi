package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	clientVersionCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "offi_client_version_requests_counter",
			Help: "Counter for client version in request",
		}, []string{"version"})

	clientBrowserCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "offi_client_browser_requests_counter",
			Help: "Counter for client browser types in request",
		}, []string{"browser"})
)

func init() {
	prometheus.MustRegister(clientVersionCounter)
	prometheus.MustRegister(clientBrowserCounter)
}

func clientVersionMiddleware(ctx *fiber.Ctx) error {
	clientVersionCounter.WithLabelValues(ctx.Query("version")).Inc()
	clientBrowserCounter.WithLabelValues(ctx.Query("browser")).Inc()
	return ctx.Next()
}
