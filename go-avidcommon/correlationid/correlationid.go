package correlationid

// This code is used in the logger subsystem. Do **NOT** import
// any other avid packages here as an import cycle may occur.
// Be very careful if you are contemplating adding logger
// statements here.
import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

const (
	CorrelationIDKey = "rkvst-correlation-id"
	RequestIDKey     = "x-request-id"
	TraceIDKey       = "x-b3-traceid"
)

// Order is important here - first one found is used
var keys = []string{TraceIDKey, RequestIDKey}

// FromContext gets the correlationId from the context if
// it exists. An empty string is returned otherwise.
func FromContext(ctx context.Context) string {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	cid, correlationIDFound := md[CorrelationIDKey]
	if correlationIDFound && len(cid) > 0 && cid[0] != "" {
		return cid[0]
	}
	return ""
}

// ContextWithCorrelationID adds the correlation id to a new ctx if not present
// currently only used in azbus receiver.
func ContextWithCorrelationID(ctx context.Context, correlationID string) context.Context {
	ctx = contextSetID(ctx, CorrelationIDKey, correlationID)

	// its ok to overwrite these values as they should not exist
	ctx = contextSetID(ctx, TraceIDKey, correlationID)
	ctx = contextSetID(ctx, RequestIDKey, correlationID)
	return ctx
}

func contextSetID(ctx context.Context, key, value string) context.Context {
	var md metadata.MD
	var newMD metadata.MD

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		// No metadata so create with ID
		md = metadata.Pairs(key, value)
		return metadata.NewIncomingContext(ctx, md)
	}

	cid, idFound := md[key]
	if idFound && len(cid) > 0 {
		// if correlationId is found and is valid we dont bother changing anything
		if cid[0] == value {
			return ctx
		}
		// if we have a different value then overwrite
		// md Should not be modified, may cause races
		// https://github.com/grpc/grpc-go/blob/89faf1c3e8283dd3c863b877bcf1631d1fe6f50c/metadata/metadata.go#L166
		md = md.Copy()
		md.Set(key, value)
		return metadata.NewIncomingContext(ctx, md)
	}

	// correlation id is not present - add using id value
	newMD = metadata.Pairs(key, value)
	return metadata.NewIncomingContext(
		ctx,
		metadata.Join(md, newMD),
	)
}

// Context gets or creates the correlationId and create a
// context with the correlation id in the metadata. The correlation ID is
// set to the trace id if it exists. Otherwise it is created from new uuid.
// This function is idempotent.
func Context(ctx context.Context) context.Context {
	ctx, cid := contextWithID(ctx)
	ctx = contextSetID(ctx, TraceIDKey, cid)
	ctx = contextSetID(ctx, RequestIDKey, cid)
	return ctx
}

func contextWithID(ctx context.Context) (context.Context, string) {
	var md metadata.MD
	var newMD metadata.MD

	newValue := uuid.New().String()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		// No metadata so create with correlation ID
		md = metadata.Pairs(CorrelationIDKey, newValue)
		return metadata.NewIncomingContext(ctx, md), newValue
	}

	// if correlationId is found and is valid we dont bother with TraceID
	cid, correlationIDFound := md[CorrelationIDKey]
	if correlationIDFound && len(cid) > 0 && cid[0] != "" {
		return ctx, cid[0]
	}

	for _, key := range keys {
		id, idFound := md[key]
		if idFound && len(id) > 0 && id[0] != "" {
			// correlation id is present but is invalid - replace with id value
			if correlationIDFound {
				// md Should not be modified, may cause races
				// https://github.com/grpc/grpc-go/blob/89faf1c3e8283dd3c863b877bcf1631d1fe6f50c/metadata/metadata.go#L166
				md = md.Copy()

				md.Set(CorrelationIDKey, id[0])
				return metadata.NewIncomingContext(ctx, md), id[0]
			}
			// correlation id is not present - add using id value
			newMD = metadata.Pairs(CorrelationIDKey, id[0])
			return metadata.NewIncomingContext(
				ctx,
				metadata.Join(md, newMD),
			), id[0]
		}
	}
	// correlation id is present but is invalid - replace with uuid
	if correlationIDFound {
		// md Should not be modified, may cause races
		// https://github.com/grpc/grpc-go/blob/89faf1c3e8283dd3c863b877bcf1631d1fe6f50c/metadata/metadata.go#L166
		md = md.Copy()
		md.Set(CorrelationIDKey, newValue)
		return metadata.NewIncomingContext(ctx, md), newValue
	}

	// correlation id is not present - add with uuid
	newMD = metadata.Pairs(CorrelationIDKey, newValue)
	return metadata.NewIncomingContext(
		ctx,
		metadata.Join(md, newMD),
	), newValue
}
