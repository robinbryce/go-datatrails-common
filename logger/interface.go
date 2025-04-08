package logger

// This is the external interface to the logger package.

const (
	DebugLevel = "DEBUG"
	InfoLevel  = "INFO"
)

type Logger interface {
	Debugf(string, ...any)
	DebugR(string, ...any)

	// Not used as we want defers to work everywhere
	//Fatalf(string, ...any)

	Infof(string, ...any)
	InfoR(string, ...any)

	Panicf(string, ...any)
	Check(string) bool

	WithIndex(string, string) *WrappedLogger
	WithServiceName(string) *WrappedLogger
	Close()

	WithOptions(...Option) *WrappedLogger
}
