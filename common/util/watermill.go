package util

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
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
		return nil, err
	}

	router.AddMiddleware(
		middleware.Recoverer,
		middleware.CorrelationID,
	)

	return router, nil
}

// PublishJSONMessage marshals payload, attaches metadata, and publishes to the topic.
func PublishJSONMessage(_ context.Context, publisher message.Publisher, topic string, payload any, metadata map[string]string) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	msg := message.NewMessage(watermill.NewUUID(), data)
	for key, value := range metadata {
		if value == "" {
			continue
		}
		msg.Metadata.Set(key, value)
	}

	if err = publisher.Publish(topic, msg); err != nil {
		return fmt.Errorf("publish topic %s: %w", topic, err)
	}
	return nil
}
