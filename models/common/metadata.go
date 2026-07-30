package common

import "context"

// Metadata contains identifiers used to trace a message through a saga.
type Metadata struct {
	MessageID     string
	CorrelationID string
	CausationID   string
}

// NewRootMetadata creates metadata for the first message in a correlation.
func NewRootMetadata(messageID, correlationID string) Metadata {
	return Metadata{
		MessageID:     messageID,
		CorrelationID: correlationID,
	}
}

// NewChildMetadata creates metadata linked to its parent message.
func NewChildMetadata(messageID string, parent Metadata) Metadata {
	return Metadata{
		MessageID:     messageID,
		CorrelationID: parent.CorrelationID,
		CausationID:   parent.MessageID,
	}
}

type metadataContextKey struct{}

// WithMetadata stores typed metadata in the context. A nil context is treated
// as context.Background().
func WithMetadata(ctx context.Context, metadata Metadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, metadataContextKey{}, metadata)
}

// MetadataFromContext returns typed metadata from the context, or its zero
// value when the context is nil or does not contain metadata.
func MetadataFromContext(ctx context.Context) Metadata {
	if ctx == nil {
		return Metadata{}
	}
	metadata, ok := ctx.Value(metadataContextKey{}).(Metadata)
	if !ok {
		return Metadata{}
	}
	return metadata
}
