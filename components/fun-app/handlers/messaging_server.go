package handlers

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/fun"
)

// MessagingServer builds and owns the Watermill router and saga handlers wiring.
type MessagingServer struct {
	router     *message.Router
	publisher  message.Publisher
	subscriber message.Subscriber
}

// NewMessagingServer constructs router, attaches middlewares, and registers topic consumers.
func NewMessagingServer(
	logger watermill.LoggerAdapter,
	publisher message.Publisher,
	subscriber message.Subscriber,
	enrollmentHandler EnrollmentMessageHandler,
	seatHandler SeatMessageHandler,
) (*MessagingServer, error) {
	router, err := util.NewRouter(logger)
	if err != nil {
		return nil, fmt.Errorf("new watermill router: %w", err)
	}

	server := &MessagingServer{
		router:     router,
		publisher:  publisher,
		subscriber: subscriber,
	}

	if err := server.registerHandlers(enrollmentHandler, seatHandler); err != nil {
		return nil, err
	}

	return server, nil
}

func (ms *MessagingServer) registerHandlers(
	enrollmentHandler EnrollmentMessageHandler,
	seatHandler SeatMessageHandler,
) error {
	ordinaryConsumers := []struct {
		topic   string
		handler message.NoPublishHandlerFunc
	}{
		{fun.TopicEnrollCmd, enrollmentHandler.HandleEnrollCmd},
		{fun.TopicSeatReservedEvt, seatHandler.HandleSeatReservedEvt},
		{fun.TopicSeatWaitlistedEvt, seatHandler.HandleSeatWaitlistedEvt},
		{fun.TopicEnrollmentCancelledEvt, enrollmentHandler.HandleEnrollmentCancelledEvt},
	}

	for _, consumer := range ordinaryConsumers {
		ms.registerConsumer(consumer.topic, consumer.handler)
	}

	retryConsumers := []struct {
		topic   string
		handler message.NoPublishHandlerFunc
	}{
		{fun.TopicSeatAllocationFailedEvt, seatHandler.HandleSeatAllocationFailedEvt},
		{fun.TopicAllocateSeatCmd, seatHandler.HandleAllocateSeatCmd},
		{util.DeadLetterTopic(fun.TopicAllocateSeatCmd), seatHandler.HandleDeadLetteredAllocateSeatCmd},
		{fun.TopicEnrollmentConfirmedEvt, enrollmentHandler.HandleEnrollmentConfirmedEvt},
	}

	for _, consumer := range retryConsumers {
		err := ms.registerRetryConsumer(consumer.topic, consumer.handler)
		if err != nil {
			return err
		}
	}
	return nil
}

// registerConsumer wires a consumer with only the server's standard recovery middleware.
func (ms *MessagingServer) registerConsumer(topic string, handler message.NoPublishHandlerFunc) {
	h := ms.router.AddConsumerHandler(topic, topic, ms.subscriber, handler)
	h.AddMiddleware(middleware.Recoverer)
}

// registerRetryConsumer wires retrying consumers and their dead-letter destinations.
func (ms *MessagingServer) registerRetryConsumer(topic string, handler message.NoPublishHandlerFunc) error {
	h := ms.router.AddConsumerHandler(topic, topic, ms.subscriber, handler)
	poisonMiddleware, err := middleware.PoisonQueue(ms.publisher, util.DeadLetterTopic(topic))
	if err != nil {
		return fmt.Errorf("create dead-letter middleware for %s: %w", topic, err)
	}
	h.AddMiddleware(poisonMiddleware)
	h.AddMiddleware(util.DefaultRetryConfig().Middleware)
	h.AddMiddleware(middleware.Recoverer)
	return nil
}

// Router exposes the configured Watermill router for lifecycle control.
func (ms *MessagingServer) Router() *message.Router { return ms.router }
