package tracing

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"offi/internal/closer"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func Init(ctx context.Context) {
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithHost(),
	)
	if err != nil {
		log.Fatal(err)
	}

	otlpEndpoint, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if !ok {
		return
	}

	if _, ok = os.LookupEnv("OTEL_EXPORTER_DISABLE_METRICS"); !ok {
		if err = initMeterProvider(ctx, res, otlpEndpoint); err != nil {
			log.Fatal(err)
		}
	}

	if err = initTracerProvider(ctx, res, otlpEndpoint); err != nil {
		log.Fatal(err)
	}
}

// Initializes an OTLP exporter, and configures the corresponding trace provider.
func initTracerProvider(ctx context.Context, res *resource.Resource, endpoint string) error {
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	closer.AddContext(tracerProvider.Shutdown)

	return nil
}

// Initializes an OTLP exporter, and configures the corresponding meter provider.
func initMeterProvider(ctx context.Context, res *resource.Resource, endpoint string) error {
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(5*time.Second))),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	closer.AddContext(meterProvider.Shutdown)

	return nil
}

// Server is a generic ogen server type.
type Server[R Route] interface {
	FindPath(method string, u *url.URL) (r R, _ bool)
}

// Route is a generic ogen route type.
type Route interface {
	Name() string
	OperationID() string
}

func NewMiddleware[R Route, S Server[R]](finder S) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "",
			otelhttp.WithServerName("offi"),
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				route, ok := finder.FindPath(r.Method, r.URL)
				if !ok {
					return route.OperationID()
				}
				return fmt.Sprintf("%s %s", route.Name(), operation)
			}),
		)
	}
}

func InjectTracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(w.Header()))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OTelHTTPTransport(next http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(next)
}
