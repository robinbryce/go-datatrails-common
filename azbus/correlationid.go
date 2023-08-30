package azbus

import (
	"context"

	"github.com/rkvst/go-rkvstcommon/correlationid"
)

func ContextFromReceivedMessage(ctx context.Context, message *ReceivedMessage) context.Context {
	if message.ApplicationProperties == nil {
		return ctx
	}
	cid, cidFound := message.ApplicationProperties[correlationid.CorrelationIDKey]
	if !cidFound {
		return ctx
	}
	return correlationid.ContextWithCorrelationID(ctx, cid.(string))
}

func AddCorrelationIDOption(ctx context.Context, opts ...OutMessageOption) []OutMessageOption {
	correlationID := correlationid.FromContext(ctx)
	if correlationID == "" {
		return opts
	}
	return append(opts, WithProperty(correlationid.CorrelationIDKey, correlationID))
}
