package common

// Metadata represents a mutable key/value map for event publishing.
type Metadata map[string]string

// NewMetadata constructs a Metadata map populated with the supplied entries.
func NewMetadata(entries map[string]string) Metadata {
	meta := make(Metadata, len(entries))
	for k, v := range entries {
		if v != "" {
			meta[k] = v
		}
	}
	return meta
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

// Clone produces a shallow copy of the metadata map.
func (m Metadata) Clone() Metadata {
	cloned := make(Metadata, len(m))
	for k, v := range m {
		cloned[k] = v
	}
	return cloned
}
