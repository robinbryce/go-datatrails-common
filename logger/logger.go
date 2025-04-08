package logger

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var (
	Sugar      *WrappedLogger
	undoLogger func()
	Recorded   *observer.ObservedLogs
)

const (
	serviceNameKey = "servicename"
)

// so we dont have to import zap everywhere
type Option = zap.Option

type WrappedLogger struct {
	*zap.SugaredLogger
}

func (l *WrappedLogger) ErrorR(msg string, args ...any) {
	keyVals := []any{}

	for i, v := range args {
		keyVals = append(keyVals, fmt.Sprintf("arg%d", i))
		keyVals = append(keyVals, v)
	}

	l.WithOptions(zap.AddCallerSkip(1)).Errorw(msg, keyVals...)
}

func (l *WrappedLogger) InfoR(msg string, args ...any) {
	keyVals := []any{}

	for i, v := range args {
		keyVals = append(keyVals, fmt.Sprintf("arg%d", i))
		keyVals = append(keyVals, v)
	}

	l.WithOptions(zap.AddCallerSkip(1)).Infow(msg, keyVals...)
}

func (l *WrappedLogger) DebugR(msg string, args ...any) {
	keyVals := []any{}

	for i, v := range args {
		keyVals = append(keyVals, fmt.Sprintf("arg%d", i))
		keyVals = append(keyVals, v)
	}

	l.WithOptions(zap.AddCallerSkip(1)).Debugw(msg, keyVals...)
}

func (l *WrappedLogger) Check(lvl string) bool {
	switch lvl {
	case InfoLevel:
		ce := l.Desugar().Check(zap.InfoLevel, "")
		return ce != nil
	case DebugLevel:
		fallthrough
	default:
		ce := l.Desugar().Check(zap.DebugLevel, "")
		return ce != nil
	}
}

// OnExit should be deferred immediately after calling the
// New() method.
func OnExit() {
	_ = Sugar.Sync()
	undoLogger()
	Recorded = nil
}

// Resource - the counter is initialised with a zero value which indicates that
// the uber correction is made (default).
type Resource struct {
	console  bool
	filename string
}

type ResourceOption func(*Resource)

func WithFile(filename string) ResourceOption {
	return func(r *Resource) {
		r.filename = filename
	}
}

func WithConsole() ResourceOption {
	return func(r *Resource) {
		r.console = true
	}
}

type Syncer func() error

// New creates 2 loggers (plain and sugared) as global variables according
// to the desired loglevel ("DEBUG", "NOOP", "TEST", default is "INFO").
// Additionally log output from other loggers in 3rd-party packages
// is redirected to the INFO label of these loggers.
// Both ResourceOption and zap.Option types are supported option types. The
// zap.Options are passed on the to zap logger.
func New(level string, opts ...any) {
	r := &Resource{}

	for _, iopt := range opts {
		if opt, ok := iopt.(ResourceOption); ok {
			opt(r)
		}
	}

	var zopts []zap.Option
	for _, opt := range opts {
		if opt, ok := opt.(zap.Option); ok {
			zopts = append(zopts, opt)
		}
	}

	var err error
	var plain *zap.Logger
	// Use opinionated presets for now.
	switch level {
	case "DEBUG":
		cfg := zap.NewDevelopmentConfig()
		if r.filename != "" {
			cfg.OutputPaths = []string{r.filename}
		}
		if r.console {
			cfg.Encoding = "console"
			cfg.EncoderConfig = zapcore.EncoderConfig{
				MessageKey: "message",
			}
		}
		plain, err = cfg.Build(zopts...)
		if err != nil {
			log.Panicf("cannot initialise zap logger: %v", err)
		}

	case "NOOP":
		plain = zap.NewNop()

	case "TEST":
		core, recorded := observer.New(zapcore.DebugLevel)

		ram := zap.WrapCore(
			func(zapcore.Core) zapcore.Core {
				return core
			},
		)

		cfg := zap.NewDevelopmentConfig()
		if r.filename != "" {
			cfg.OutputPaths = []string{r.filename}
		}
		if r.console {
			cfg.Encoding = "console"
			cfg.EncoderConfig = zapcore.EncoderConfig{
				MessageKey: "message",
			}
		}
		plain, err = cfg.Build(zopts...)
		if err != nil {
			log.Panicf("cannot initialise zap logger: %v", err)
		}
		plain = plain.WithOptions(ram)
		Recorded = recorded

	default:
		cfg := zap.NewProductionConfig()
		if r.filename != "" {
			cfg.OutputPaths = []string{r.filename}
		}
		if r.console {
			cfg.Encoding = "console"
			cfg.EncoderConfig = zapcore.EncoderConfig{
				MessageKey: "message",
			}
		}
		plain, err = cfg.Build(zopts...)
		if err != nil {
			log.Panicf("cannot initialise zap logger: %v", err)
		}
	}
	undoLogger = zap.RedirectStdLog(plain)
	Sugar = &WrappedLogger{
		plain.Sugar(),
	}
}

func (wl *WrappedLogger) WithServiceName(servicename string) *WrappedLogger {
	return wl.WithIndex(serviceNameKey, servicename)
}

func (wl *WrappedLogger) WithIndex(key, value string) *WrappedLogger {
	return &WrappedLogger{
		wl.SugaredLogger.With(zap.String(key, strings.ToLower(value))),
	}
}

func (wl *WrappedLogger) WithOptions(opts ...Option) *WrappedLogger {
	s := &WrappedLogger{
		wl.SugaredLogger.WithOptions(opts...),
	}
	return s
}

// Close attempts to flush any buffered log entries.
func (wl *WrappedLogger) Close() {
	err := wl.Sync()

	// not alot we can do other than log that we couldn't flush the log
	// This is usually an error 'sync /dev/stderr invalid argument'
	// which is pointless
	if err != nil && !errors.Is(err, syscall.EINVAL) {
		wl.Debugf("Close: Failed to flush log: %v", err)
	}
}
