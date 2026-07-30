package handlers

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/rs/zerolog/log"
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

	server.registerHandlers(enrollmentHandler, seatHandler)

	return server, nil
}

func (ms *MessagingServer) registerHandlers(
	enrollmentHandler EnrollmentMessageHandler,
	seatHandler SeatMessageHandler,
) {
	ms.registerConsumer(fun.TopicEnrollCmd, enrollmentHandler.HandleEnrollCmd)
	ms.registerConsumer(fun.TopicSeatReservedEvt, seatHandler.HandleSeatReservedEvt)
	ms.registerConsumer(fun.TopicSeatWaitlistedEvt, seatHandler.HandleSeatWaitlistedEvt)
	ms.registerConsumer(fun.TopicEnrollmentCancelledEvt, enrollmentHandler.HandleEnrollmentCancelledEvt)

	ms.registerRetryConsumer(fun.TopicSeatAllocationFailedEvt, seatHandler.HandleSeatAllocationFailedEvt)
	ms.registerRetryConsumer(fun.TopicAllocateSeatCmd, seatHandler.HandleAllocateSeatCmd)
	ms.registerRetryConsumer(util.DeadLetterTopic(fun.TopicAllocateSeatCmd), seatHandler.HandleDeadLetteredAllocateSeatCmd)
	ms.registerRetryConsumer(fun.TopicEnrollmentConfirmedEvt, enrollmentHandler.HandleEnrollmentConfirmedEvt)
}

// registerConsumer wires a consumer with metadata normalization and recovery.
func (ms *MessagingServer) registerConsumer(topic string, handler message.NoPublishHandlerFunc) {
	h := ms.router.AddConsumerHandler(topic, topic, ms.subscriber, handler)
	h.AddMiddleware(util.SagaMetadataMiddleware())
	h.AddMiddleware(middleware.Recoverer)
}

// registerRetryConsumer wires retrying consumers and their dead-letter destinations.
func (ms *MessagingServer) registerRetryConsumer(topic string, handler message.NoPublishHandlerFunc) {
	h := ms.router.AddConsumerHandler(topic, topic, ms.subscriber, handler)
	poisonMiddleware, err := middleware.PoisonQueue(ms.publisher, util.DeadLetterTopic(topic))
	if err != nil {
		log.Fatal().Err(err).Str("topic", topic).Msg("failed to create dead-letter middleware")
	}
	h.AddMiddleware(poisonMiddleware)
	h.AddMiddleware(util.SagaMetadataMiddleware())
	h.AddMiddleware(util.DefaultRetryConfig().Middleware)
	h.AddMiddleware(middleware.Recoverer)
}

// Router exposes the configured Watermill router for lifecycle control.
func (ms *MessagingServer) Router() *message.Router { return ms.router }
