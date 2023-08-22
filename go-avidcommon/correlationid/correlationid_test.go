package correlationid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

// Test_Context tests that a suitable correlationID is emitted
// when traceID, correlationID may or may not exist in all possible combinations
func Test_Context(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		expected string
		length   int
	}{
		{
			name: "no metadata",
			args: args{
				ctx: context.Background(),
			},
			length: 36,
		},
		{
			name: "no traceid or request id or correlation id",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{"somekey": {"somekey"}}),
			},
			length: 36,
		},
		{
			name: "trace id",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{TraceIDKey: {"traceid"}}),
			},
			expected: "traceid",
		},
		{
			name: "empty traceid",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{TraceIDKey: {""}}),
			},
			length: 36,
		},
		{
			name: "null traceid",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{TraceIDKey: {}}),
			},
			length: 36,
		},
		{
			name: "request id",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{RequestIDKey: {"requestid"}}),
			},
			expected: "requestid",
		},
		{
			name: "empty requestid",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{RequestIDKey: {""}}),
			},
			length: 36,
		},
		{
			name: "null requestid",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{RequestIDKey: {}}),
			},
			length: 36,
		},
		{
			name: "empty correlationid",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{CorrelationIDKey: {""}}),
			},
			length: 36,
		},
		{
			name: "null correlationid",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{CorrelationIDKey: {}}),
			},
			length: 36,
		},
		{
			name: "correlationid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(), metadata.MD{CorrelationIDKey: {"correlationid"}}),
			},
			expected: "correlationid",
		},
		{
			name: "traceid and requestid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.MD{
						RequestIDKey: {"requestid"},
						TraceIDKey:   {"traceid"},
					},
				),
			},
			expected: "traceid",
		},
		{
			name: "empty traceid and requestid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.MD{
						RequestIDKey: {"requestid"},
						TraceIDKey:   {""},
					},
				),
			},
			expected: "requestid",
		},
		{
			name: "null traceid and requestid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.MD{
						RequestIDKey: {"requestid"},
						TraceIDKey:   {},
					},
				),
			},
			expected: "requestid",
		},
		{
			name: "traceid and empty requestid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.MD{
						RequestIDKey: {""},
						TraceIDKey:   {"traceid"},
					},
				),
			},
			expected: "traceid",
		},
		{
			name: "traceid and null requestid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.MD{
						RequestIDKey: {},
						TraceIDKey:   {"traceid"},
					},
				),
			},
			expected: "traceid",
		},
		{
			name: "traceid and correlationid and requestid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.MD{
						CorrelationIDKey: {"correlationid"},
						RequestIDKey:     {"requestid"},
						TraceIDKey:       {"traceid"},
					},
				),
			},
			expected: "correlationid",
		},
		{
			name: "traceid and null correlationid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.MD{
						CorrelationIDKey: {},
						TraceIDKey:       {"traceid"},
					},
				),
			},
			expected: "traceid",
		},
		{
			name: "requestid and null correlationid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.MD{
						CorrelationIDKey: {},
						RequestIDKey:     {"requestid"},
					},
				),
			},
			expected: "requestid",
		},
		{
			name: "traceid and empty correlationid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.MD{
						CorrelationIDKey: {""},
						TraceIDKey:       {"traceid"},
					},
				),
			},
			expected: "traceid",
		},
		{
			name: "requestid and empty correlationid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.MD{
						CorrelationIDKey: {""},
						RequestIDKey:     {"requestid"},
					},
				),
			},
			expected: "requestid",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := Context(test.args.ctx)
			actual := FromContext(ctx)

			if test.expected != "" {
				assert.Equal(t, test.expected, actual)
			}
			if test.length > 0 {
				assert.Equal(t, test.length, len(actual))
			}
		})
	}
}
