package azbus

import (
	"context"

	"github.com/datatrails/go-datatrails-common/spanner"
	"github.com/datatrails/go-datatrails-common/tracing"
)

func (r *Receiver) handleReceivedMessageWithTracingContext(ctx context.Context, log Logger, message *ReceivedMessage, handler Handler) (Disposition, context.Context, error) {
	var span spanner.Spanner
	span, ctx = tracing.NewSpanWithAttributes(ctx, "Receiver", log, message.ApplicationProperties)
	defer span.Close()
	return handler.Handle(ctx, message)
}

func (s *Sender) updateSendingMesssageForSpan(ctx context.Context, message *OutMessage, span spanner.Spanner) {
	log := tracing.LogFromContext(ctx, s.log)
	defer log.Close()

	for k, v := range span.Attributes(log) {
		OutMessageSetProperty(message, k, v)
	}
	log.Debugf("updateSendingMesssageForSpan(): ApplicationProperties %v", OutMessageProperties(message))
}
