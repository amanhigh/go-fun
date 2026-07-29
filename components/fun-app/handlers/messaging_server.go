package handlers

import (
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/fun"
)

// Scoped retry configuration for the allocation command handler.
const (
	allocateSeatRetryMax      = 2
	allocateSeatRetryInterval = 2 * time.Second
)

// MessagingServer builds and owns the Watermill router and saga handlers wiring.
type MessagingServer struct {
	router *message.Router
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

	// Register ordinary topic consumers (no handler-level middleware beyond automatic Recoverer).
	util.AddConsumerHandler(router, fun.TopicEnrollCmd, subscriber, enrollmentHandler.HandleEnrollCmd)
	util.AddConsumerHandler(router, fun.TopicSeatReservedEvt, subscriber, seatHandler.HandleSeatReservedEvt)
	util.AddConsumerHandler(router, fun.TopicSeatWaitlistedEvt, subscriber, seatHandler.HandleSeatWaitlistedEvt)
	util.AddConsumerHandler(router, fun.TopicEnrollmentConfirmedEvt, subscriber, enrollmentHandler.HandleEnrollmentConfirmedEvt)
	util.AddConsumerHandler(router, fun.TopicEnrollmentCancelledEvt, subscriber, enrollmentHandler.HandleEnrollmentCancelledEvt)

	// Allocation command handler: PoisonQueue (outermost) → Retry → Recoverer → handler.
	// The poison topic is derived from the source topic via PoisonTopic.
	if err := util.AddRetryPoisonConsumerHandler(router, publisher, subscriber, util.RetryPoisonConsumerConfig{
		Topic:         fun.TopicAllocateSeatCmd,
		Retry:         middleware.Retry{MaxRetries: allocateSeatRetryMax, InitialInterval: allocateSeatRetryInterval},
		Handler:       seatHandler.HandleAllocateSeatCmd,
		PoisonHandler: seatHandler.HandlePoisonedAllocateSeatCmd,
	}); err != nil {
		return nil, fmt.Errorf("allocation retry-poison handler: %w", err)
	}

	return &MessagingServer{router: router}, nil
}

// Router exposes the configured Watermill router for lifecycle control.
func (ms *MessagingServer) Router() *message.Router { return ms.router }
