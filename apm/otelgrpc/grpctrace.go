package otelgrpc

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

const (
	// instrumentationName is the name of this instrumentation package.
	instrumentationName = "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	// GRPCStatusCodeKey is convention for numeric status code of a gRPC request.
	GRPCStatusCodeKey = attribute.Key("rpc.grpc.status_code")
)

// Filter is a predicate used to determine whether a given request in
// interceptor info should be traced. A Filter must return true if
// the request should be traced.
type Filter func(*InterceptorInfo) bool

// config is a group of options for this instrumentation.
type config struct {
	Filter         Filter
	Propagators    propagation.TextMapPropagator
	TracerProvider trace.TracerProvider
}

// Option applies an option value for a config.
type Option interface {
	apply(*config)
}

// newConfig returns a config configured with all the passed Options.
func newConfig(opts []Option) *config {
	c := &config{
		Propagators:    otel.GetTextMapPropagator(),
		TracerProvider: otel.GetTracerProvider(),
	}
	for _, o := range opts {
		o.apply(c)
	}
	return c
}

type propagatorsOption struct{ p propagation.TextMapPropagator }

func (o propagatorsOption) apply(c *config) {
	if o.p != nil {
		c.Propagators = o.p
	}
}

// WithPropagators returns an Option to use the Propagators when extracting
// and injecting trace context from requests.
func WithPropagators(p propagation.TextMapPropagator) Option {
	return propagatorsOption{p: p}
}

type tracerProviderOption struct{ tp trace.TracerProvider }

func (o tracerProviderOption) apply(c *config) {
	if o.tp != nil {
		c.TracerProvider = o.tp
	}
}

// WithInterceptorFilter returns an Option to use the request filter.
func WithInterceptorFilter(f Filter) Option {
	return interceptorFilterOption{f: f}
}

type interceptorFilterOption struct {
	f Filter
}

func (o interceptorFilterOption) apply(c *config) {
	if o.f != nil {
		c.Filter = o.f
	}
}

// WithTracerProvider returns an Option to use the TracerProvider when
// creating a Tracer.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return tracerProviderOption{tp: tp}
}

type metadataSupplier struct {
	metadata *metadata.MD
}

// assert that metadataSupplier implements the TextMapCarrier interface.
var _ propagation.TextMapCarrier = &metadataSupplier{}

func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (s *metadataSupplier) Set(key string, value string) {
	s.metadata.Set(key, value)
}

func (s *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(*s.metadata))
	for key := range *s.metadata {
		out = append(out, key)
	}
	return out
}

// Inject injects correlation context and span context into the gRPC
// metadata object. This function is meant to be used on outgoing
// requests.
// Deprecated: Unnecessary public func.
func Inject(ctx context.Context, md *metadata.MD, opts ...Option) {
	c := newConfig(opts)
	c.Propagators.Inject(ctx, &metadataSupplier{
		metadata: md,
	})
}

func inject(ctx context.Context, propagators propagation.TextMapPropagator) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	propagators.Inject(ctx, &metadataSupplier{
		metadata: &md,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

// Extract returns the correlation context and span context that
// another service encoded in the gRPC metadata object with Inject.
// This function is meant to be used on incoming requests.
// Deprecated: Unnecessary public func.
func Extract(ctx context.Context, md *metadata.MD, opts ...Option) (baggage.Baggage, trace.SpanContext) {
	c := newConfig(opts)
	ctx = c.Propagators.Extract(ctx, &metadataSupplier{
		metadata: md,
	})

	return baggage.FromContext(ctx), trace.SpanContextFromContext(ctx)
}

func extract(ctx context.Context, propagators propagation.TextMapPropagator) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	return propagators.Extract(ctx, &metadataSupplier{
		metadata: &md,
	})
}
