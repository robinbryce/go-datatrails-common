package azblob

import (
	"context"

	"github.com/datatrails/go-datatrails-common/logger"
	"github.com/datatrails/go-datatrails-common/spanner"
)

type Spanner = spanner.Spanner

type startSpanFromContextFunc func(context.Context, logger.Logger, string) (spanner.Spanner, context.Context)
