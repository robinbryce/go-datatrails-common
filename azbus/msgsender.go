package azbus

import (
	"context"
)

type MsgSender interface {
	Open() error
	Send(context.Context, []byte, ...OutMessageOption) error
	SendMsg(context.Context, OutMessage, ...OutMessageOption) error
	Close(context.Context)
	String() string
	GetAZClient() AZClient
}
