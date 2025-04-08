package tracing

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
)

// LogFromContext takes the trace ID from the current span and adds it to the logger:
//
// returns:
//   - the new logger with a context metadata value for traceID
//
// This will be called on entry to a method or a function that has a context.Context.
func LogFromContext(ctx context.Context, log Logger) Logger {

	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		log.Infof("FromContext: span is nil - this should not happen - the context where this happened is missing tracing content - probably a middleware problem")
		return log
	}
	carrier := opentracing.TextMapCarrier{}
	err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, carrier)
	if err != nil {
		log.Debugf("FromContext: can't inject span: %v", err)
		return log
	}

	traceID, found := carrier[TraceID]
	if !found || traceID == "" {
		log.Debugf("%s not found", TraceID)
		return log
	}

	return log.WithIndex(TraceID, traceID)
}
