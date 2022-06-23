package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var clientVersionCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "offi_client_version_requests_counter",
		Help: "Counter for client version in request",
	},
	[]string{"version"},
)

func clientVersionMiddleware(ctx *fiber.Ctx) error {
	clientVersionCounter.With(prometheus.Labels{"version": ctx.Params("version")}).Inc()
	return ctx.Next()
}
