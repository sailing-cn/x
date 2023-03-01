package apm

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
	"sailing.cn/apm/otelgrpc"
)

// TraceConfig 跟踪配置
type TraceConfig struct {
	Service  string `json:"service" yaml:"service"`   //服务名称
	Version  string `json:"version" yaml:"version"`   //服务版本
	Exporter string `json:"exporter" yaml:"exporter"` //存储器
	Host     string `json:"host" yaml:"host"`         //存储地址
}

func InitTracer(cfg *TraceConfig) func(ctx context.Context) error {
	provider := getTracerProvider(cfg)
	return provider.Shutdown
}

func GetGrpcServerTraceOptions() []grpc.ServerOption {
	var options = []grpc.ServerOption{
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	}
	return options
}

func GetGrpcClientTraceOption() []grpc.DialOption {
	var options = []grpc.DialOption{
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	}
	return options
}

func getTracerProvider(cfg *TraceConfig) *sdktrace.TracerProvider {
	var exporter trace.SpanExporter
	var err error
	switch cfg.Exporter {
	case "zipkin":
		{
			exporter, err = zipkin.New(cfg.Host)
			break
		}
	case "jaeger":
		{
			exporter, err = jaeger.New(
				jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Host)),
			)
			break
		}
	}

	if err != nil {
		log.Errorf("系统异常%s", err.Error())
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Service),
			semconv.ServiceVersionKey.String(cfg.Version),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
