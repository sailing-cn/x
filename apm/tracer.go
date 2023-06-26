package apm

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.9.0"
	"google.golang.org/grpc"
)

type Config struct {
	Service  string `json:"service" yaml:"service"`   //服务名称
	Version  string `json:"version" yaml:"version"`   //服务版本
	Exporter string `json:"exporter" yaml:"exporter"` //存储器
	Host     string `json:"host" yaml:"host"`         //存储地址
}

func NewConfig(paths ...string) *Config {
	viper.AddConfigPath("./conf.d")
	viper.SetConfigName("conf")
	err := viper.ReadInConfig()
	if err != nil {
		panic("读取配置文件出错:" + err.Error())
	}
	conf := &Config{
		Service:  viper.GetString("grpc.name"),
		Version:  viper.GetString("grpc.version"),
		Exporter: viper.GetString("apm.exporter"),
		Host:     viper.GetString("apm.host"),
	}
	return conf
}

func NewTracer(cfg *Config) func(ctx context.Context) error {
	provider := getTracerProvider(cfg)
	return provider.Shutdown
}

func NewGrpcTraceOptions() []grpc.ServerOption {
	var options = []grpc.ServerOption{
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	}
	return options
}

func getTracerProvider(cfg *Config) *sdk.TracerProvider {
	var exporter sdk.SpanExporter
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
	tp := sdk.NewTracerProvider(
		sdk.WithSampler(sdk.AlwaysSample()),
		sdk.WithBatcher(exporter),
		sdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Service),
			semconv.ServiceVersionKey.String(cfg.Version),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
