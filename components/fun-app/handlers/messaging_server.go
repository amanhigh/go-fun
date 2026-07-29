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

	// Register topic consumers with handler-level middleware.
	addEnrollmentCmdHandlers(router, subscriber, enrollmentHandler)
	addSeatHandlers(router, subscriber, seatHandler, allocPoisonMw, allocRetry)
	addEnrollmentEvtHandlers(router, subscriber, enrollmentHandler)
	addPoisonHandlers(router, subscriber, seatHandler)

	// Router-level Recoverer: innermost router middleware, applied innermost relative
	// to the handler-level PoisonQueue → Retry chain, producing PoisonQueue → Retry → Recoverer → handler.
	router.AddMiddleware(middleware.Recoverer)

	return &MessagingServer{router: router}, nil
}

// Router exposes the configured Watermill router for lifecycle control.
func (ms *MessagingServer) Router() *message.Router { return ms.router }

// addEnrollmentCmdHandlers registers enrollment command handlers.
func addEnrollmentCmdHandlers(router *message.Router, subscriber message.Subscriber, enrollmentHandler EnrollmentMessageHandler) {
	router.AddConsumerHandler(
		fun.RouterHandlerIDRequestSeatAllocation,
		fun.TopicEnrollCmd,
		subscriber,
		enrollmentHandler.HandleEnrollCmd,
	)
}

// addSeatHandlers registers seat handlers. The allocation command handler receives
// scoped PoisonQueue → Retry handler-level middleware; other seat handlers have none.
func addSeatHandlers(
	router *message.Router,
	subscriber message.Subscriber,
	seatHandler SeatMessageHandler,
	allocPoisonMw message.HandlerMiddleware,
	allocRetry message.HandlerMiddleware,
) {
	// Allocation command handler: PoisonQueue (outermost) → Retry → Recoverer (innermost) → handler.
	h := router.AddConsumerHandler(
		fun.RouterHandlerIDAllocateSeat,
		fun.TopicAllocateSeatCmd,
		subscriber,
		seatHandler.HandleAllocateSeatCmd,
	)
	router.AddConsumerHandler(
		fun.RouterHandlerIDConfirmEnrollment,
		fun.TopicSeatReservedEvt,
		subscriber,
		seatHandler.HandleSeatReservedEvt,
	)

	router.AddConsumerHandler(
		fun.RouterHandlerIDWaitlistEnrollment,
		fun.TopicSeatWaitlistedEvt,
		subscriber,
		seatHandler.HandleSeatWaitlistedEvt,
	)
}

// addEnrollmentEvtHandlers registers enrollment event handlers.
func addEnrollmentEvtHandlers(router *message.Router, subscriber message.Subscriber, enrollmentHandler EnrollmentMessageHandler) {
	router.AddConsumerHandler(
		fun.RouterHandlerIDRecordEnrollmentConfirmation,
		fun.TopicEnrollmentConfirmedEvt,
		subscriber,
		enrollmentHandler.HandleEnrollmentConfirmedEvt,
	)

	router.AddConsumerHandler(
		fun.RouterHandlerIDRecordEnrollmentCancellation,
		fun.TopicEnrollmentCancelledEvt,
		subscriber,
		enrollmentHandler.HandleEnrollmentCancelledEvt,
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
