package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// InitTracer mengonfigurasi OpenTelemetry untuk mencetak jejak (trace) ke Terminal
func InitTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	// 1. Buat Exporter (Tujuan pembuangan data jejak, dalam hal ini ke terminal/stdout)
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	// 2. Buat Resource tanpa menyertakan Schema URL agar tidak terjadi bentrok versi
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"", // 👈 URL Schema dikosongkan
			attribute.String("service.name", serviceName), // 👈 Menggunakan attribute bawaan
		),
	)
	if err != nil {
		return nil, err
	}

	// 3. Gabungkan menjadi Tracer Provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// 4. Set Tracer Provider ini sebagai standar global di seluruh aplikasi Go
	otel.SetTracerProvider(tp)

	return tp, nil
}
