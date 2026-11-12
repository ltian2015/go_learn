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
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc" // Re-enable this import

	// "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp" // Remove this import
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc" // Re-enable this import
	// "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp" // Remove this import
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
	fmt.Println("Initializing Tracer Provider...")
	traceExporter, err := otlptracegrpc.New(ctx, // ✅ 改回 gRPC 导出器
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("host.docker.internal:4317"), // ✅ 尝试 host.docker.internal
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
	fmt.Println("Initializing Logger Provider...")
	logExporter, err := otlploggrpc.New(ctx, // ✅ 改回 gRPC 导出器
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithEndpoint("host.docker.internal:4317"), // ✅ 尝试 host.docker.internal
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create log exporter: %w", err)
	}

	lp := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(res),
	)
	// otel.SetLoggerProvider(lp) // This was commented out already
	return lp, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received a request to handler.")
	ctx := r.Context()
	tr := otel.Tracer(serviceName)
	ctx, span := tr.Start(ctx, "handler-request")
	fmt.Printf("Span ID: %s, Trace ID: %s\n", span.SpanContext().SpanID().String(), span.SpanContext().TraceID().String())
	defer span.End()

	logger := slog.Default().With("trace_id", span.SpanContext().TraceID().String(), "span_id", span.SpanContext().SpanID().String())

	logger.InfoContext(ctx, "Received a request",
		attribute.String("method", r.Method),
		attribute.String("path", r.URL.Path),
	)
	fmt.Println("Logged info context.")

	time.Sleep(50 * time.Millisecond) // Simulate some work

	slog.WarnContext(ctx, "Processing request with a warning",
		attribute.Int("processing_time_ms", 50),
	)
	fmt.Println("Logged warn context.")

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
		fmt.Println("Shutting down Tracer Provider...") // Added for debugging
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
		fmt.Println("Shutting down Logger Provider...") // Added for debugging
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
	<-ctx.Done() // 阻塞，直到接收到中断信号
	slog.Default().Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Default().Error("Server shutdown failed", "error", err)
	}
	fmt.Println("Server shutdown complete.") // Added for debugging
	slog.Default().Info("Server gracefully stopped.")

	// Add a small delay to ensure all buffered data is flushed before main exits
	// This is a last resort for debugging, not for production
	time.Sleep(2 * time.Second) // 增加一个短延迟，确保数据有时间发出
}
