package tracing

import (
	"context"
	"net/http"

	"github.com/datatrails/go-datatrails-common/logger"
	"github.com/datatrails/go-datatrails-common/spanner"
	opentracing "github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
)

type Spanner interface {
	Close()
	SetTag(string, any)
	SetSpanHTTPHeader(http.Request)
	SetSpanField(string, string)
	TraceID() string
}

// Injecting a function StartSpanFomContext that returns an interface does not work as
// the Go compiler treats interfaces defined in separate packages as different even
// though the signature is identical. So this struct hides the interface and returns a
// concrete type instead.
// Conveniently it also hides calls to the opentracing-go package thsu making it
// easier to move to opentelemetry later.
type Span struct {
	span opentracing.Span
	log  logger.Logger
}

func (s *Span) Close() {
	if s.span != nil {
		s.span.Finish()
		s.span = nil
	}
}

func (s *Span) SetTag(key string, value any) {
	if s.span != nil {
		s.span.SetTag(key, value)
	}
}

func (s *Span) SetSpanHTTPHeader(r *http.Request) {
	// Transmit the span's TraceContext as HTTP headers on our request
	if s.span != nil {
		err := opentracing.GlobalTracer().Inject(
			s.span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header),
		)
		if err != nil {
			s.log.Infof("Tracer.Inject %v", err)
		}
	}
}

func (s *Span) LogField(key string, value any) {
	if s.span != nil {
		switch v := value.(type) {
		case bool:
			s.span.LogFields(opentracinglog.Bool(key, v))
		case error:
			s.span.LogFields(opentracinglog.Error(v))
		case int:
			s.span.LogFields(opentracinglog.Int(key, v))
		case int32:
			s.span.LogFields(opentracinglog.Int32(key, v))
		case int64:
			s.span.LogFields(opentracinglog.Int64(key, v))
		case uint32:
			s.span.LogFields(opentracinglog.Uint32(key, v))
		case uint64:
			s.span.LogFields(opentracinglog.Uint64(key, v))
		case float32:
			s.span.LogFields(opentracinglog.Float32(key, v))
		case float64:
			s.span.LogFields(opentracinglog.Float64(key, v))
		case string:
			s.span.LogFields(opentracinglog.String(key, v))
		}
	}
}

func valueFromCarrier(carrier opentracing.TextMapCarrier, key string) string {
	value, found := carrier[key]
	if !found || value == "" {
		return ""
	}
	return value
}

func (s *Span) TraceID() string {
	if s.span == nil {
		return ""
	}
	carrier := opentracing.TextMapCarrier{}
	err := opentracing.GlobalTracer().Inject(s.span.Context(), opentracing.TextMap, carrier)
	if err != nil {
		return ""
	}

	return valueFromCarrier(carrier, TraceID)
}

func (s *Span) Attributes(log logger.Logger) map[string]any {
	var attributes = make(map[string]any)
	if s.span == nil {
		return attributes
	}

	carrier := opentracing.TextMapCarrier{}
	err := opentracing.GlobalTracer().Inject(s.span.Context(), opentracing.TextMap, carrier)
	if err != nil {
		log.Infof("attributes(): Unable to inject span context: %v", err)
		return attributes
	}
	for k, v := range carrier {
		attributes[k] = v
	}
	log.Debugf("Attributes(): %v", attributes)
	return attributes
}

// Constructors...
func NewSpanWithAttributes(ctx context.Context, name string, log logger.Logger, attributes map[string]any) (spanner.Spanner, context.Context) {
	log.Debugf("NewSpanWithAttributes %s", name)
	var opts = []opentracing.StartSpanOption{}
	carrier := opentracing.TextMapCarrier{}
	// This just gets all the attributes into a string map. That map is then passed into the
	// open tracing constructor which extracts any bits it is interested in to use to setup the spans etc.
	// It will ignore anything it doesn't care about. So the filtering of the map is done for us and
	// we don't need to pre-filter it.
	for k, v := range attributes {
		// Tracing properties will be strings
		value, ok := v.(string)
		if ok {
			carrier.Set(k, value)
		}
	}
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, carrier)
	if err != nil {
		log.Infof("NewSpanWithAttributes(): Unable to extract span context: %v", err)
	} else {
		opts = append(opts, opentracing.ChildOf(spanCtx))
	}
	span := opentracing.StartSpan(name, opts...)
	ctx = opentracing.ContextWithSpan(ctx, span)
	return &Span{span: span, log: log}, ctx
}

func NewSpanContext(ctx context.Context, log logger.Logger, name string) (spanner.Spanner, context.Context) {
	log.Debugf("NewSpanContext %s", name)
	span := opentracing.StartSpan(name)
	if span == nil {
		return nil, ctx
	}
	ctx = opentracing.ContextWithSpan(ctx, span)
	return &Span{span: span, log: log}, ctx
}

func StartSpanFromContext(ctx context.Context, log logger.Logger, name string) (spanner.Spanner, context.Context) {
	log.Debugf("StartSpanFromContext %s", name)
	span, ctx := opentracing.StartSpanFromContext(ctx, name)
	return &Span{span: span, log: log}, ctx
}
