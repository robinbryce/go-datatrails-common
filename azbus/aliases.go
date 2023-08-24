package azbus

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"

	"github.com/rkvst/go-rkvstcommon/logger"
)

type Logger = logger.Logger

// so we dont have to import the azure repo everywhere
type OutMessage = azservicebus.Message

type ReceivedMessage = azservicebus.ReceivedMessage

func NewOutMessage(data []byte) OutMessage {
	return azservicebus.Message{
		Body: data,
	}
}

// not an alias but its convenient here
type Handler interface {
	Handle(context.Context, *ReceivedMessage) error
}
