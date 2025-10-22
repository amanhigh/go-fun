package util

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	modelcommon "github.com/amanhigh/go-fun/models/common"
	"github.com/pkg/errors"
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

// NewRouter creates a Watermill router with default middleware.
func NewRouter(logger watermill.LoggerAdapter) (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, errors.Wrap(err, "new watermill router")
	}

	router.AddMiddleware(
		middleware.Recoverer,
		middleware.CorrelationID,
	)

	return router, nil
}

// PublishJSONMessage marshals payload, attaches metadata, and publishes to the topic.
func PublishJSONMessage(_ context.Context, publisher message.Publisher, topic string, payload any, metadata modelcommon.Metadata) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	// Correlation is mandatory for all saga messages.
	if c := metadata[modelcommon.MetadataCorrelationIDKey]; c == "" {
		return fmt.Errorf("missing %s", modelcommon.MetadataCorrelationIDKey)
	}

	id := watermill.NewUUID()
	msg := message.NewMessage(id, data)

	// Copy provided metadata.
	for key, value := range metadata {
		msg.Metadata.Set(key, value)
	}

	// Always mirror the message id for downstream consumers.
	msg.Metadata.Set(modelcommon.MetadataMessageIDKey, id)

	if err = publisher.Publish(topic, msg); err != nil {
		return fmt.Errorf("publish topic %s: %w", topic, err)
	}
	return nil
}
