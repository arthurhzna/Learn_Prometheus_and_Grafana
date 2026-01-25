package instrumentation

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func InitTraceProvider(ctx context.Context, serviceName string, otlpEndpoint string) (func(context.Context) error, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
			attribute.String("environment", "development"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(otlpEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(1)),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tracerProvider.Shutdown, nil
}

func GetTracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

type SpanState string

const (
	StateStarted    SpanState = "started"
	StateProcessing SpanState = "processing"
	StateSuccess    SpanState = "success"
	StateFailed     SpanState = "failed"
	StateRetrying   SpanState = "retrying"
	StateConnecting SpanState = "connecting"
	StateExecuting  SpanState = "executing"
)

func AddEvent(span trace.Span, eventName string, state SpanState, attrs ...attribute.KeyValue) {
	allAttrs := append([]attribute.KeyValue{
		attribute.String("state", string(state)),
	}, attrs...)
	span.AddEvent(eventName, trace.WithAttributes(allAttrs...))
}

func RecordError(span trace.Span, err error, errorType, component string, attrs ...attribute.KeyValue) {
	allAttrs := append([]attribute.KeyValue{
		attribute.String("error.type", errorType),
		attribute.String("error.message", err.Error()),
		attribute.String("error.component", component),
	}, attrs...)

	span.RecordError(err, trace.WithAttributes(allAttrs...))
	span.SetAttributes(attribute.Bool("error", true))
}

func DBAttributes(system, operation, table, query string) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("db.system", system),
		attribute.String("db.operation", operation),
		attribute.String("db.table", table),
		attribute.String("db.query", query),
	}
}

func HTTPAttributes(method, url string, statusCode int) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("http.method", method),
		attribute.String("http.url", url),
		attribute.Int("http.status_code", statusCode),
	}
}

func GRPCAttributes(service, method string, statusCode int) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("rpc.system", "grpc"),
		attribute.String("rpc.service", service),
		attribute.String("rpc.method", method),
		attribute.Int("rpc.grpc.status_code", statusCode),
	}
}

func StartSpanWithState(ctx context.Context, tracer trace.Tracer, spanName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, spanName, trace.WithAttributes(attrs...))
	AddEvent(span, spanName+".started", StateStarted)
	return ctx, span
}

func EndSpanWithSuccess(span trace.Span, spanName string, duration time.Duration, attrs ...attribute.KeyValue) {
	span.SetAttributes(append(attrs,
		attribute.Bool("success", true),
		attribute.Int64("duration_ms", duration.Milliseconds()),
	)...)
	AddEvent(span, spanName+".completed", StateSuccess,
		attribute.Int64("duration_ms", duration.Milliseconds()),
	)
	span.End()
}

func EndSpanWithError(span trace.Span, spanName string, err error, errorType string, duration time.Duration) {
	RecordError(span, err, errorType, spanName)
	span.SetAttributes(
		attribute.Bool("success", false),
		attribute.Int64("duration_ms", duration.Milliseconds()),
	)
	AddEvent(span, spanName+".failed", StateFailed,
		attribute.String("error", err.Error()),
	)
	span.End()
}
