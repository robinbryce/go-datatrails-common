package tracing

import (
	"context"

	"github.com/datatrails/go-datatrails-common/logger"
	opentracing "github.com/opentracing/opentracing-go"
)

// LogFromContext takes the trace ID from the current span and adds it to the logger:
//
// returns:
//   - the new logger with a context metadata value for traceID
//
// This will be called on entry to a method or a function that has a context.Context.
func LogFromContext(ctx context.Context, log logger.Logger) logger.Logger {
	traceID := TraceIDFromContext(ctx, log)
	if traceID != "" {
		return log.WithIndex(TraceID, traceID)
	}
	return log
}

func TraceIDFromContext(ctx context.Context, log logger.Logger) string {

	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		log.Infof("FromContext: span is nil - this should not happen - the context where this happened is missing tracing content - probably a middleware problem")
		return ""
	}
	carrier := opentracing.TextMapCarrier{}
	err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, carrier)
	if err != nil {
		log.Debugf("FromContext: can't inject span: %v", err)
		return ""
	}

	traceID, found := carrier[TraceID]
	if !found || traceID == "" {
		log.Debugf("%s not found", TraceID)
		return ""
	}

	return traceID
}
