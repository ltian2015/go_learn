package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const (
	serviceName    = "simple-web-service"
	serviceVersion = "1.0.0"
)

func initResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
	)
}

func initTracerProvider(ctx context.Context, res *resource.Resource) (*trace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:4317"), // OTLP gRPC endpoint for traces
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func initLoggerProvider(ctx context.Context, res *resource.Resource) (*log.LoggerProvider, error) {
	logExporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithEndpoint("localhost:4317"), // OTLP gRPC endpoint for logs
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create log exporter: %w", err)
	}

	lp := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(res),
	)

	//otel.SetLoggerProvider(lp)
	return lp, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := otel.Tracer(serviceName)
	ctx, span := tr.Start(ctx, "handler-request")
	defer span.End()

	logger := slog.Default().With("trace_id", span.SpanContext().TraceID().String(), "span_id", span.SpanContext().SpanID().String())

	logger.InfoContext(ctx, "Received a request",
		attribute.String("method", r.Method),
		attribute.String("path", r.URL.Path),
	)

	time.Sleep(50 * time.Millisecond) // Simulate some work

	slog.WarnContext(ctx, "Processing request with a warning",
		attribute.Int("processing_time_ms", 50),
	)

	fmt.Fprintf(w, "Hello, OpenTelemetry with Slog!", time.Now())
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// 1. Initialize OpenTelemetry Resource
	res := initResource()

	// 2. Initialize Tracer Provider
	tp, err := initTracerProvider(ctx, res)
	if err != nil {
		slog.Default().Error("Failed to initialize trace provider", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			slog.Default().Error("Error shutting down tracer provider", "error", err)
		}
	}()

	// 3. Initialize Logger Provider
	lp, err := initLoggerProvider(ctx, res)
	if err != nil {
		slog.Default().Error("Failed to initialize logger provider", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := lp.Shutdown(ctx); err != nil {
			slog.Default().Error("Error shutting down logger provider", "error", err)
		}
	}()

	// 4. Set Slog to use OpenTelemetry Logger
	slog.SetDefault(slog.New(otelslog.NewHandler("otel-slog-handler", otelslog.WithLoggerProvider(lp))))
	// 5. Start HTTP Server
	http.HandleFunc("/", handler)
	port := ":8080"
	slog.Default().Info("Starting server", "port", port)
	server := &http.Server{Addr: port}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Default().Error("Server error", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	<-ctx.Done()
	slog.Default().Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Default().Error("Server shutdown failed", "error", err)
	}
	slog.Default().Info("Server gracefully stopped.")
}
