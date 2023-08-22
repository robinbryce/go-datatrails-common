package logger

// This is the external interface to the logger pacakage.
import (
	"context"
)

type Logger interface {
	Debugf(string, ...any)
	DebugR(string, ...any)

	// Not used as we wamt deferes to work everywhere
	//Fatalf(string, ...any)

	Infof(string, ...any)
	InfoR(string, ...any)

	Panicf(string, ...any)

	FromContext(context.Context) *WrappedLogger
	WithIndex(string, string) *WrappedLogger
	WithServiceName(string) *WrappedLogger
	Close()

	WithOptions(...Option) *WrappedLogger
}
