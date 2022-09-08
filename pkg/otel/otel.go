package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	"go.opentelemetry.io/otel/trace"
)

var serviceName string

func Tracer() trace.Tracer {
	return otel.Tracer(serviceName)
}

func SetTracerProvider(name, environment string) (tp *tracesdk.TracerProvider, err error) {
	serviceName = name

	exp, err := otlptrace.New(context.Background(), otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint("otel-collector:4318"),
		otlptracehttp.WithInsecure(),
	))
	if err != nil {
		return nil, err
	}

	tp = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", environment),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, err
}
