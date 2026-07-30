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
	router      *message.Router
	publisher   message.Publisher
	subscriber  message.Subscriber
	diagnostics util.MetadataDiagnostics
}

// NewMessagingServer constructs router, attaches middlewares, and registers topic consumers.
func NewMessagingServer(
	logger watermill.LoggerAdapter,
	publisher message.Publisher,
	subscriber message.Subscriber,
	enrollmentHandler EnrollmentMessageHandler,
	seatHandler SeatMessageHandler,
	diagnostics util.MetadataDiagnostics,
) (*MessagingServer, error) {
	router, err := util.NewRouter(logger)
	if err != nil {
		return nil, fmt.Errorf("new watermill router: %w", err)
	}

	server := &MessagingServer{
		router:      router,
		publisher:   publisher,
		subscriber:  subscriber,
		diagnostics: diagnostics,
	}

	server.registerHandlers(enrollmentHandler, seatHandler)

	return server, nil
}

func (ms *MessagingServer) registerHandlers(
	enrollmentHandler EnrollmentMessageHandler,
	seatHandler SeatMessageHandler,
) {
	ms.registerConsumer(fun.TopicEnrollCmd, enrollmentHandler.HandleEnrollCmd, util.RootMessageRole)
	ms.registerConsumer(fun.TopicSeatReservedEvt, seatHandler.HandleSeatReservedEvt, util.DownstreamMessageRole)
	ms.registerConsumer(fun.TopicSeatWaitlistedEvt, seatHandler.HandleSeatWaitlistedEvt, util.DownstreamMessageRole)
	ms.registerConsumer(fun.TopicEnrollmentCancelledEvt, enrollmentHandler.HandleEnrollmentCancelledEvt, util.DownstreamMessageRole)

	ms.registerRetryConsumer(fun.TopicSeatAllocationFailedEvt, seatHandler.HandleSeatAllocationFailedEvt, util.DownstreamMessageRole)
	ms.registerRetryConsumer(fun.TopicAllocateSeatCmd, seatHandler.HandleAllocateSeatCmd, util.DownstreamMessageRole)
	ms.registerRetryConsumer(util.DeadLetterTopic(fun.TopicAllocateSeatCmd), seatHandler.HandleDeadLetteredAllocateSeatCmd, util.DownstreamMessageRole)
	ms.registerRetryConsumer(fun.TopicEnrollmentConfirmedEvt, enrollmentHandler.HandleEnrollmentConfirmedEvt, util.DownstreamMessageRole)
}

// registerConsumer wires a consumer with metadata normalization and recovery.
func (ms *MessagingServer) registerConsumer(topic string, handler message.NoPublishHandlerFunc, role util.MessageRole) {
	h := ms.router.AddConsumerHandler(topic, topic, ms.subscriber, handler)
	h.AddMiddleware(util.SagaMetadataMiddleware(role, ms.diagnostics))
	h.AddMiddleware(middleware.Recoverer)
}

// registerRetryConsumer wires retrying consumers and their dead-letter destinations.
func (ms *MessagingServer) registerRetryConsumer(topic string, handler message.NoPublishHandlerFunc, role util.MessageRole) {
	h := ms.router.AddConsumerHandler(topic, topic, ms.subscriber, handler)
	poisonMiddleware, err := middleware.PoisonQueue(ms.publisher, util.DeadLetterTopic(topic))
	if err != nil {
		log.Fatal().Err(err).Str("topic", topic).Msg("failed to create dead-letter middleware")
	}
	h.AddMiddleware(poisonMiddleware)
	h.AddMiddleware(util.SagaMetadataMiddleware(role, ms.diagnostics))
	h.AddMiddleware(util.DefaultRetryConfig().Middleware)
	h.AddMiddleware(middleware.Recoverer)
}

// Router exposes the configured Watermill router for lifecycle control.
func (ms *MessagingServer) Router() *message.Router { return ms.router }
