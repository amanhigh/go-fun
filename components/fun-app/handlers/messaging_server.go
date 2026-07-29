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

	// Global middlewares: retry and poison queue
	router.AddMiddleware(
		middleware.Retry{MaxRetries: wmRetryMax, InitialInterval: wmRetryInterval}.Middleware,
	)
	poisonMw, perr := middleware.PoisonQueue(publisher, wmPoisonTopic)
	if perr != nil {
		return nil, fmt.Errorf("poison middleware: %w", perr)
	}
	router.AddMiddleware(poisonMw)

	// Register topic consumers
	addEnrollmentCommandHandlers(router, subscriber, enrollmentHandler)
	addSeatHandlers(router, subscriber, seatHandler)
	addEnrollmentEventHandlers(router, subscriber, enrollmentHandler)
	addPoisonHandlers(router, subscriber, seatHandler)

	return &MessagingServer{router: router}, nil
}

// Router exposes the configured Watermill router for lifecycle control.
func (ms *MessagingServer) Router() *message.Router { return ms.router }

func addEnrollmentCommandHandlers(router *message.Router, subscriber message.Subscriber, enrollmentHandler EnrollmentCommandHandler) {
	router.AddConsumerHandler(
		"enrollment_enroll_requested",
		fun.TopicEnrollCmd,
		subscriber,
		enrollmentHandler.EnrollCmd,
	)
}

func addSeatHandlers(router *message.Router, subscriber message.Subscriber, seatHandler SeatCommandHandler) {
	router.AddConsumerHandler(
		"seat_allocate",
		fun.TopicAllocateSeatCmd,
		subscriber,
		seatHandler.AllocateSeatCmd,
	)
	router.AddConsumerHandler(
		"seat_reserved_evt",
		fun.TopicSeatReservedEvt,
		subscriber,
		seatHandler.SeatReservedEvt,
	)
	router.AddConsumerHandler(
		"seat_waitlisted_evt",
		fun.TopicSeatWaitlistedEvt,
		subscriber,
		seatHandler.SeatWaitlistedEvt,
	)
}

func addEnrollmentEventHandlers(router *message.Router, subscriber message.Subscriber, enrollmentHandler EnrollmentCommandHandler) {
	router.AddConsumerHandler(
		"enrollment_confirmed_evt",
		fun.TopicEnrollmentConfirmedEvt,
		subscriber,
		enrollmentHandler.EnrollmentConfirmedEvt,
	)
	router.AddConsumerHandler(
		"enrollment_cancelled_evt",
		fun.TopicEnrollmentCancelledEvt,
		subscriber,
		enrollmentHandler.EnrollmentCancelledEvt,
	)
}

// addPoisonHandlers consumes messages from poison topic to perform final cancellation.
func addPoisonHandlers(router *message.Router, subscriber message.Subscriber, seatHandler SeatCommandHandler) {
	router.AddConsumerHandler(
		"poison_allocate",
		wmPoisonTopic,
		subscriber,
		func(msg *message.Message) error {
			// BUG: Is Poison Handler Worknig ?
			// Try decode as AllocateSeatCmdV1 and cancel via handler
			var cmd fun.AllocateSeatCmdV1
			if err := json.Unmarshal(msg.Payload, &cmd); err == nil {
				return seatHandler.PoisonAllocate(msg)
			}
			// Unknown poison payload: ack as no-op
			return nil
		},
	)
}
