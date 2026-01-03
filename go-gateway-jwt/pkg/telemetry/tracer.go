package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// Hàm khởi tạo Tracer
func InitTracer(serviceName string, collectorURL string) func(context.Context) error {
	ctx := context.Background()

	// 1. Tạo Exporter gửi dữ liệu về Jaeger qua gRPC
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(collectorURL),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Không thể khởi tạo Otel Exporter: %v", err)
	}

	// 2. Định nghĩa Resource (Tên service, version...)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		log.Fatalf("Không thể tạo Resource: %v", err)
	}

	// 3. Tạo Tracer Provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), // Gom nhiều trace gửi 1 lần cho nhẹ mạng
		sdktrace.WithResource(res),
	)

	// 4. Set làm Provider toàn cục
	otel.SetTracerProvider(tp)

	// 5. Cấu hình Propagator (Để truyền TraceID từ Service A sang Service B)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown
}
