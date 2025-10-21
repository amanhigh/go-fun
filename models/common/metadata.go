package common

import "context"

// Key names shared across saga publishers/consumers.
const (
	// MetadataCorrelationIDKey is a saga-wide correlation identifier.
	MetadataCorrelationIDKey = "correlation_id"
	// MetadataCausationIDKey holds the triggering message ID for this emission.
	MetadataCausationIDKey = "causation_id"
	// MetadataMessageIDKey mirrors the transport message ID for consumers.
	MetadataMessageIDKey = "message_id"
)

// Metadata represents a mutable key/value map for event publishing.
type Metadata map[string]string

// MustBaseMetadata constructs metadata with a mandatory correlation id.
// It panics if correlationID is empty and should be used only when the caller is sure the id exists.
func MustBaseMetadata(correlationID string) Metadata {
	if correlationID == "" {
		panic("correlation id is required")
	}
	return Metadata{MetadataCorrelationIDKey: correlationID}
}

// With returns a clone containing the additional key/value when value is non-empty.
func (m Metadata) With(key, value string) Metadata {
	if value == "" {
		return m.Clone()
	}

	cloned := m.Clone()
	cloned[key] = value
	return cloned
}

// WithCausation appends the causation id when non-empty.
func (m Metadata) WithCausation(id string) Metadata {
	return m.With(MetadataCausationIDKey, id)
}

// WithPair appends a single key/value when non-empty.
func (m Metadata) WithPair(key, value string) Metadata { return m.With(key, value) }

// WithPairs merges multiple keys, skipping empty values.
func (m Metadata) WithPairs(pairs map[string]string) Metadata {
	out := m.Clone()
	for k, v := range pairs {
		if v == "" {
			continue
		}
		out[k] = v
	}
	return out
}

// Clone produces a shallow copy of the metadata map.
func (m Metadata) Clone() Metadata {
	cloned := make(Metadata, len(m))
	for k, v := range m {
		cloned[k] = v
	}
	return cloned
}

type metadataCtxKey string

const (
	correlationCtxKey metadataCtxKey = "metadata_correlation_id"
	causationCtxKey   metadataCtxKey = "metadata_causation_id"
)

// WithCorrelation stores the correlation identifier in the context. Panics if id is empty.
func WithCorrelation(ctx context.Context, id string) context.Context {
	if id == "" {
		panic("correlation id cannot be empty")
	}
	return context.WithValue(ctx, correlationCtxKey, id)
}

// CorrelationFrom extracts the correlation identifier from the context, returning empty string when absent.
func CorrelationFrom(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if val, ok := ctx.Value(correlationCtxKey).(string); ok {
		return val
	}
	return ""
}

// WithCausation stores the causation identifier in the context. Panics if id is empty.
func WithCausation(ctx context.Context, id string) context.Context {
	if id == "" {
		panic("causation id cannot be empty")
	}
	return context.WithValue(ctx, causationCtxKey, id)
}

// CausationFrom extracts the causation identifier from the context, returning empty string when absent.
func CausationFrom(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if val, ok := ctx.Value(causationCtxKey).(string); ok {
		return val
	}
	return ""
}
