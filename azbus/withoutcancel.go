package azbus

import (
	"context"
	"time"
)

// The azservicebus obeys a context.WithDeadline if present. However we have
// learned that retries are better so all disposition functions and the sender
// suppress any deadlines in the context. This has the added benefit of not
// reusing a context that has already cancelled - which means that otherwise
// disposition code will exit immediately.
//
// This problem of context having 2 responsibilities (breaking the single
// responsibility principle) is known - see this proposal:
//
//	https://github.com/golang/go/issues/40221
//
// context.WithoutCancel() was committed 2023-03-29 https://go-review.googlesource.com/c/go/+/479918
// So we can assume it will be available in 1.21 in August 2023
//
// meanwhile this is the workaround until Go 1.21 when this file will be deleted
type withoutCancelCtx struct {
	context.Context
}

func contextWithoutCancel(ctx context.Context) context.Context {
	return withoutCancelCtx{ctx}
}
func (withoutCancelCtx) Deadline() (deadline time.Time, ok bool) { return }
func (withoutCancelCtx) Done() <-chan struct{}                   { return nil }
func (withoutCancelCtx) Err() error                              { return nil }

// func (c withoutCancelCtx) Value(key any) any { return value(c, key) }
func (c withoutCancelCtx) String() string { return "WithoutCancel" }
