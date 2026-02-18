package common

// https://github.com/omniti-labs/jsend
// Envelope response types for all API responses.
// All API responses follow this standard format.
// EnvelopeStatus represents the status field in envelope responses.
type EnvelopeStatus string

const (
	// EnvelopeSuccess indicates the request was successful.
	EnvelopeSuccess EnvelopeStatus = "success"
	// EnvelopeFail indicates the request failed due to client error (4xx).
	EnvelopeFail EnvelopeStatus = "fail"
	// EnvelopeError indicates the request failed due to server error (5xx).
	EnvelopeError EnvelopeStatus = "error"
)

// Envelope is the generic response envelope for all responses.
// Use this for 2xx responses where data contains the actual payload.
// For 4xx responses, data should contain field-specific error messages.
// For 5xx responses, data should contain error details.
type Envelope[T any] struct {
	Status EnvelopeStatus `json:"status"`
	Data   T              `json:"data"`
}

// NewEnvelope creates a success response with the given data.
func NewEnvelope[T any](data T) Envelope[T] {
	return Envelope[T]{
		Status: EnvelopeSuccess,
		Data:   data,
	}
}

// NewFailEnvelope creates a fail response with field-specific errors.
func NewFailEnvelope(errors map[string]string) Envelope[map[string]string] {
	return Envelope[map[string]string]{
		Status: EnvelopeFail,
		Data:   errors,
	}
}

// NewErrorEnvelope creates an error response with a message and optional code.
func NewErrorEnvelope(message string, code int) Envelope[map[string]any] {
	return Envelope[map[string]any]{
		Status: EnvelopeError,
		Data: map[string]any{
			"message": message,
			"code":    code,
		},
	}
}
