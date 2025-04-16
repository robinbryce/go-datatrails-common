package spanner

import (
	"net/http"

	"github.com/datatrails/go-datatrails-common/logger"
)

// this interface is in a separate package such that azblob and tracing packages can share the
// same interface definition that is returned fro the StartSpanFromContext etc.. methods.
type Spanner interface {
	Close()
	SetTag(string, any)
	SetSpanHTTPHeader(*http.Request)
	Attributes(logger.Logger) map[string]any
	LogField(string, any)
	TraceID() string
}
