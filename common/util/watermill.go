package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	modelcommon "github.com/amanhigh/go-fun/models/common"
)

// WatermillLifecycle abstracts router lifecycle management.
// NewStdWatermillLogger returns Watermill's default stdout logger.
func NewStdWatermillLogger() watermill.LoggerAdapter {
	return watermill.NewStdLogger(false, false)
}

// NewGoChannel constructs an in-memory pub/sub channel with the given logger.
func NewGoChannel(logger watermill.LoggerAdapter) *gochannel.GoChannel {
	return gochannel.NewGoChannel(gochannel.Config{}, logger)
}

// NewRouter creates a Watermill router without correlation middleware.
//
// This application uses NoPublishHandlerFunc handlers and publishes directly
// through managers and publishers, so correlation metadata is set at publish time.
func NewRouter(logger watermill.LoggerAdapter) (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, fmt.Errorf("new watermill router: %w", err)
	}

	return router, nil
}

// DeadLetterTopic returns the derived dead-letter topic name for the given source topic.
func DeadLetterTopic(topic string) string {
	return topic + ".dead-letter"
}

// DefaultRetryConfig returns the default retry configuration: two retries with
// a two-second initial interval, including the standard ShouldRetry
// classification. Callers may mutate fields directly on the returned value.
func DefaultRetryConfig() middleware.Retry {
	return middleware.Retry{
		MaxRetries:      2,
		InitialInterval: 2 * time.Second,
		ShouldRetry:     shouldRetry,
	}
}

// shouldRetry classifies errors for retry decisions. Transient technical
// failures should be retried, while permanent input or client failures should
// be sent directly to their final handling path.
//
// Non-retryable: malformed JSON (json.SyntaxError), type errors
// (json.UnmarshalTypeError), and ordinary 4xx common.HttpError values.
//
// Retryable: 408 Request Timeout, 429 Too Many Requests, 5xx common.HttpError,
// and unknown errors.
func shouldRetry(params middleware.RetryParams) bool {
	err := params.Err
	// HTTP error classification.
	var httpErr modelcommon.HttpError
	if errors.As(err, &httpErr) {
		code := httpErr.Code()
		// 408 (Request Timeout) and 429 (Too Many Requests) are retryable.
		if code == http.StatusRequestTimeout || code == http.StatusTooManyRequests {
			return true
		}
		// Other 4xx errors are not retryable.
		if code >= 400 && code < 500 {
			return false
		}
		// 5xx errors are retryable.
		return true
	}

	// JSON syntax errors are not retryable.
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return false
	}

	// JSON type mismatch errors are not retryable.
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return false
	}

	// Unknown errors are retryable.
	return true
}

// PublishJSONMessage marshals payload, attaches metadata, and publishes to the topic.
func PublishJSONMessage(_ context.Context, publisher message.Publisher, topic string, payload any, metadata modelcommon.Metadata) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	if metadata.MessageID == "" {
		return fmt.Errorf("missing message id")
	}
	if metadata.CorrelationID == "" {
		return fmt.Errorf("missing correlation id")
	}

	msg := message.NewMessage(metadata.MessageID, data)
	middleware.SetCorrelationID(metadata.CorrelationID, msg)
	if metadata.CausationID != "" {
		msg.Metadata.Set("causation_id", metadata.CausationID)
	}

	if err = publisher.Publish(topic, msg); err != nil {
		return fmt.Errorf("publish topic %s: %w", topic, err)
	}
	return nil
}

// metadataFromMessage normalizes transport identity and returns typed saga
// metadata. Malformed transport metadata is recovered rather than treated as
// a business failure.
func metadataFromMessage(msg *message.Message) modelcommon.Metadata {
	if msg == nil {
		return modelcommon.Metadata{}
	}

	if msg.UUID == "" {
		msg.UUID = watermill.NewUUID()
	}

	correlationID := middleware.MessageCorrelationID(msg)
	if correlationID == "" {
		correlationID = watermill.NewUUID()
		middleware.SetCorrelationID(correlationID, msg)
	}

	return modelcommon.NewMetadata(msg.UUID, correlationID, msg.Metadata.Get("causation_id"))
}

// SagaMetadataMiddleware attaches normalized typed metadata before the
// wrapped handler. Recovery is fail-open so metadata cannot block delivery.
func SagaMetadataMiddleware() message.HandlerMiddleware {
	return func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			if msg != nil {
				metadata := metadataFromMessage(msg)
				msg.SetContext(modelcommon.WithMetadata(msg.Context(), metadata))
			}
			return next(msg)
		}
	}
}
