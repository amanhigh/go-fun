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
//
// FIXME: CorrelationID only propagates for messages that pass through the router
// (i.e. returned by handlers). Direct publishing via BasePublisher / stampCtx
// bypasses the router and must carry its own correlation metadata.
func NewRouter(logger watermill.LoggerAdapter) (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, errors.Wrap(err, "new watermill router")
	}

	router.AddMiddleware(
		middleware.CorrelationID,
	)

	return router, nil
}

// PoisonTopic returns the derived poison-topic name for the given source topic.
func PoisonTopic(topic string) string {
	return topic + ".poison"
}

// AddConsumerHandler registers a no-publish handler on the router that subscribes
// to topic using the supplied subscriber. The handler is identified by topic (used
// as both the Watermill handler name and the subscription topic). Any supplied
// middlewares are added in order, then middleware.Recoverer is appended last so
// recovery is the innermost layer. Returns the Watermill *Handler for further
// configuration if needed.
func AddConsumerHandler(
	router *message.Router,
	topic string,
	subscriber message.Subscriber,
	handler message.NoPublishHandlerFunc,
	middlewares ...message.HandlerMiddleware,
) *message.Handler {
	h := router.AddConsumerHandler(topic, topic, subscriber, handler)
	h.AddMiddleware(middlewares...)
	h.AddMiddleware(middleware.Recoverer)
	return h
}

// RetryPoisonConsumerConfig bundles the parameters for AddRetryPoisonConsumerHandler.
type RetryPoisonConsumerConfig struct {
	Topic         string
	Retry         middleware.Retry
	Handler       message.NoPublishHandlerFunc
	PoisonHandler message.NoPublishHandlerFunc
}

// AddRetryPoisonConsumerHandler wires up a retry + poison-queue pipeline for the
// given source topic. It derives a poison topic via PoisonTopic, creates a
// middleware.PoisonQueue that publishes failed messages there, and registers two
// handlers:
//  1. Source handler — runs the user handler wrapped by PoisonQueue (outer) and
//     Retry.Middleware (inner), plus the automatic Recoverer.
//  2. Poison consumer — subscribes to the poison topic and runs poisonHandler.
//
// Returns any PoisonQueue construction error. The dedicated poison consumer
// receives the automatic Recoverer from AddConsumerHandler.
func AddRetryPoisonConsumerHandler(
	router *message.Router,
	publisher message.Publisher,
	subscriber message.Subscriber,
	config RetryPoisonConsumerConfig,
) error {
	poisonTopic := PoisonTopic(config.Topic)

	poisonMiddleware, err := middleware.PoisonQueue(publisher, poisonTopic)
	if err != nil {
		return errors.Wrap(err, "create poison queue middleware")
	}

	// Source handler: PoisonQueue (outer) → Retry → Recoverer (inner) → handler
	AddConsumerHandler(router, config.Topic, subscriber, config.Handler, poisonMiddleware, config.Retry.Middleware)

	// Dedicated poison consumer on the derived poison topic
	AddConsumerHandler(router, poisonTopic, subscriber, config.PoisonHandler)

	return nil
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
