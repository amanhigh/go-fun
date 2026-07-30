package handlers

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/fun"
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

	// Ordinary consumers — no retry, no DLQ. Messages are processed once;
	// the built-in Recoverer requeues on panic but nothing beyond that.
	util.AddConsumerHandler(router, fun.TopicEnrollCmd, subscriber, enrollmentHandler.HandleEnrollCmd)
	util.AddConsumerHandler(router, fun.TopicSeatReservedEvt, subscriber, seatHandler.HandleSeatReservedEvt)
	util.AddConsumerHandler(router, fun.TopicSeatWaitlistedEvt, subscriber, seatHandler.HandleSeatWaitlistedEvt)
	util.AddConsumerHandler(router, fun.TopicSeatAllocationFailedEvt, subscriber, seatHandler.HandleSeatAllocationFailedEvt, util.ConsumerConfig{
		Retry:               util.DefaultRetryConfig(),
		DeadLetterPublisher: publisher,
	})
	util.AddConsumerHandler(router, fun.TopicEnrollmentCancelledEvt, subscriber, enrollmentHandler.HandleEnrollmentCancelledEvt)

	// Compensation DLQ — AllocateSeatCmd.
	// Retries exhausted → source-specific dead-letter topic → publish allocation failure.
	util.AddConsumerHandler(router, fun.TopicAllocateSeatCmd, subscriber, seatHandler.HandleAllocateSeatCmd, util.ConsumerConfig{
		Retry:               util.DefaultRetryConfig(),
		DeadLetterPublisher: publisher,
	})
	util.AddConsumerHandler(router, util.DeadLetterTopic(fun.TopicAllocateSeatCmd), subscriber, seatHandler.HandleDeadLetteredAllocateSeatCmd, util.ConsumerConfig{
		Retry:               util.DefaultRetryConfig(),
		DeadLetterPublisher: publisher,
	})

	// Terminal DLQ — EnrollmentConfirmedEvt.
	// Retries exhausted → source-specific dead-letter topic with NO domain handler.
	// The message lands in the terminal DLQ topic for manual inspection;
	// there is no automatic compensation because confirmation is a terminal event.
	util.AddConsumerHandler(router, fun.TopicEnrollmentConfirmedEvt, subscriber, enrollmentHandler.HandleEnrollmentConfirmedEvt, util.ConsumerConfig{
		Retry:               util.DefaultRetryConfig(),
		DeadLetterPublisher: publisher,
	})

	return &MessagingServer{router: router}, nil
}

// Router exposes the configured Watermill router for lifecycle control.
func (ms *MessagingServer) Router() *message.Router { return ms.router }
