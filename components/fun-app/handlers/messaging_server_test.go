package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/amanhigh/go-fun/common/util"
	daomocks "github.com/amanhigh/go-fun/components/fun-app/dao/mocks"
	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	managermocks "github.com/amanhigh/go-fun/components/fun-app/manager/mocks"
	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	common "github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

type recordedPublish struct {
	topic string
	msg   *message.Message
}

type recordingPublisher struct {
	delegate message.Publisher
	mu       sync.Mutex
	records  []recordedPublish
}

func (p *recordingPublisher) Publish(topic string, msg ...*message.Message) error {
	p.mu.Lock()
	for _, item := range msg {
		p.records = append(p.records, recordedPublish{topic: topic, msg: item})
	}
	p.mu.Unlock()
	return p.delegate.Publish(topic, msg...)
}

func (p *recordingPublisher) Close() error { return p.delegate.Close() }

func (p *recordingPublisher) messages(topics ...string) []recordedPublish {
	p.mu.Lock()
	defer p.mu.Unlock()
	var result []recordedPublish
	for _, record := range p.records {
		if slices.Contains(topics, record.topic) {
			result = append(result, record)
		}
	}
	return result
}

var _ = Describe("MessagingServer learning scenarios", func() {
	var (
		channel        *gochannel.GoChannel
		ms             *handlers.MessagingServer
		enrollmentMock *managermocks.EnrollmentManagerInterface
		seatMock       *managermocks.SeatManagerInterface
		logger         watermill.LoggerAdapter
		routerCtx      context.Context
		routerCancel   context.CancelFunc
	)

	BeforeEach(func() {
		logger = watermill.NewStdLogger(false, false)
		channel = gochannel.NewGoChannel(gochannel.Config{}, logger)
		enrollmentMock = managermocks.NewEnrollmentManagerInterface(GinkgoT())
		seatMock = managermocks.NewSeatManagerInterface(GinkgoT())

		enrollmentHandler := handlers.NewEnrollmentMessageHandler(enrollmentMock)
		seatHandler := handlers.NewSeatMessageHandler(seatMock, enrollmentMock)

		var err error
		ms, err = handlers.NewMessagingServer(logger, channel, channel, enrollmentHandler, seatHandler)
		Expect(err).ToNot(HaveOccurred())

		routerCtx, routerCancel = context.WithCancel(context.Background())
		go func() { _ = ms.Router().Run(routerCtx) }()
		<-ms.Router().Running()
	})

	AfterEach(func() {
		routerCancel()
		_ = ms.Router().Close()
		_ = channel.Close()
	})

	Context("allocation compensation", func() {
		var (
			cmd       fun.AllocateSeatCmdV1
			cancelled chan struct{}
		)

		BeforeEach(func() {
			cmd = fun.AllocateSeatCmdV1{
				EnrollmentID: "enr-1",
				PersonID:     "person-1",
				Grade:        3,
				RequestedAt:  time.Now().UTC(),
			}
			cancelled = make(chan struct{})
			seatMock.EXPECT().AllocateSeat(mock.Anything, cmd).
				Return(common.NewHttpError("seat service unavailable", http.StatusInternalServerError)).Times(3)
			seatMock.EXPECT().PublishSeatAllocationFailed(
				mock.Anything,
				fun.Enrollment{ID: cmd.EnrollmentID, PersonID: cmd.PersonID},
				mock.MatchedBy(func(reason string) bool {
					return reason != ""
				}),
			).Run(func(_ context.Context, enrollment fun.Enrollment, reason string) {
				evtPayload, err := json.Marshal(fun.SeatAllocationFailedEvtV1{
					EnrollmentID: enrollment.ID,
					PersonID:     enrollment.PersonID,
					Reason:       reason,
					FailedAt:     time.Now().UTC(),
				})
				if err != nil {
					return
				}
				Expect(channel.Publish(fun.TopicSeatAllocationFailedEvt, message.NewMessage("allocation-failed-1", evtPayload))).To(Succeed())
			}).Return(nil).Once()
			enrollmentMock.EXPECT().CancelEnrollmentAndPublish(mock.Anything, mock.MatchedBy(func(evt fun.EnrollmentCancelledEvtV1) bool {
				return evt.EnrollmentID == cmd.EnrollmentID && evt.PersonID == cmd.PersonID
			})).Run(func(context.Context, fun.EnrollmentCancelledEvtV1) {
				close(cancelled)
			}).Return(nil).Once()

			payload, err := json.Marshal(cmd)
			Expect(err).ToNot(HaveOccurred())
			msg := message.NewMessage("allocate-1", payload)
			Expect(channel.Publish(fun.TopicAllocateSeatCmd, msg)).To(Succeed())
			Eventually(cancelled, 10*time.Second).Should(BeClosed())
		})

		It("retries allocation twice before publishing one failure", func() {
			seatMock.AssertNumberOfCalls(GinkgoT(), "AllocateSeat", 3)
			seatMock.AssertNumberOfCalls(GinkgoT(), "PublishSeatAllocationFailed", 1)
		})
	})

	Context("malformed allocation", func() {
		var (
			deadLetterMessages <-chan *message.Message
			deadLetterMessage  *message.Message
		)

		BeforeEach(func() {
			var err error
			deadLetterMessages, err = channel.Subscribe(routerCtx, util.DeadLetterTopic(fun.TopicAllocateSeatCmd))
			Expect(err).ToNot(HaveOccurred())
			Expect(channel.Publish(fun.TopicAllocateSeatCmd, message.NewMessage("allocate-malformed", []byte("not-json")))).To(Succeed())
			Eventually(deadLetterMessages).Should(Receive(&deadLetterMessage))
		})

		It("quarantines malformed allocation without retry or compensation", func() {
			Expect(string(deadLetterMessage.Payload)).To(Equal("not-json"))
			Expect(deadLetterMessage.Metadata.Get(middleware.PoisonedTopicKey)).To(Equal(fun.TopicAllocateSeatCmd))
			Consistently(deadLetterMessages, "200ms", "50ms").ShouldNot(Receive())
			seatMock.AssertNotCalled(GinkgoT(), "AllocateSeat", mock.Anything, mock.Anything)
			seatMock.AssertNotCalled(GinkgoT(), "PublishSeatAllocationFailed", mock.Anything, mock.Anything, mock.Anything)
			enrollmentMock.AssertNotCalled(GinkgoT(), "CancelEnrollmentAndPublish", mock.Anything, mock.Anything)
		})
	})

	Context("allocation failure compensation", func() {
		var (
			deadLetterMessages <-chan *message.Message
			deadLetterMessage  *message.Message
		)

		BeforeEach(func() {
			failed := fun.SeatAllocationFailedEvtV1{
				EnrollmentID: "enr-1",
				PersonID:     "person-1",
				Reason:       "capacity unavailable",
				FailedAt:     time.Now().UTC(),
			}
			payload, err := json.Marshal(failed)
			Expect(err).ToNot(HaveOccurred())
			enrollmentMock.EXPECT().CancelEnrollmentAndPublish(mock.Anything, mock.Anything).
				Return(common.NewHttpError("database unavailable", http.StatusInternalServerError)).Times(3)
			deadLetterMessages, err = channel.Subscribe(routerCtx, util.DeadLetterTopic(fun.TopicSeatAllocationFailedEvt))
			Expect(err).ToNot(HaveOccurred())
			Expect(channel.Publish(fun.TopicSeatAllocationFailedEvt, message.NewMessage("failed-1", payload))).To(Succeed())
			Eventually(deadLetterMessages, 10*time.Second).Should(Receive(&deadLetterMessage))
		})

		It("bounds cancellation retries and preserves a terminal dead-letter record", func() {
			Expect(deadLetterMessage.Metadata.Get(middleware.PoisonedTopicKey)).To(Equal(fun.TopicSeatAllocationFailedEvt))
			enrollmentMock.AssertNumberOfCalls(GinkgoT(), "CancelEnrollmentAndPublish", 3)
		})
	})

	Context("terminal enrollment confirmation DLQ", func() {
		var (
			deadLetterMessages <-chan *message.Message
			deadLetterMessage  *message.Message
			err                error
		)

		BeforeEach(func() {
			deadLetterMessages, err = channel.Subscribe(routerCtx, util.DeadLetterTopic(fun.TopicEnrollmentConfirmedEvt))
			Expect(err).ToNot(HaveOccurred())
			payload := []byte("not-json")
			Expect(channel.Publish(fun.TopicEnrollmentConfirmedEvt, message.NewMessage("confirmation-malformed", payload))).To(Succeed())
			Eventually(deadLetterMessages).Should(Receive(&deadLetterMessage))
		})

		It("publishes one terminal dead-letter record with source metadata", func() {
			Expect(string(deadLetterMessage.Payload)).To(Equal("not-json"))
			Expect(deadLetterMessage.Metadata.Get(middleware.PoisonedTopicKey)).To(Equal(fun.TopicEnrollmentConfirmedEvt))
			enrollmentMock.AssertNotCalled(GinkgoT(), "OnEnrollmentConfirmedEvt", mock.Anything, mock.Anything)
			seatMock.AssertNotCalled(GinkgoT(), "AllocateSeat", mock.Anything, mock.Anything)
		})
	})
})

var _ = Describe("MessagingServer causal-chain scenario", func() {
	var (
		m1, m2, m3, m4 *message.Message
		recorder       *recordingPublisher
		chainChannel   *gochannel.GoChannel
		server         *handlers.MessagingServer
		routerCtx      context.Context
		routerCancel   context.CancelFunc
	)

	BeforeEach(func() {
		logger := watermill.NewStdLogger(false, false)
		chainChannel = gochannel.NewGoChannel(gochannel.Config{}, logger)
		recorder = &recordingPublisher{delegate: chainChannel}
		enrollmentDAO := daomocks.NewEnrollmentDaoInterface(GinkgoT())
		enrollment := fun.Enrollment{ID: "enr-chain", PersonID: "person-chain", Grade: 3, Status: fun.EnrollmentStatusSeatAllocationInitiated}
		enrollmentDAO.EXPECT().UseOrCreateTx(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, run util.DbRun, _ ...bool) common.HttpError {
			return run(ctx)
		}).Twice()
		findCalls := 0
		enrollmentDAO.EXPECT().FindById(mock.Anything, enrollment.ID, mock.Anything).Run(func(_ context.Context, _ any, entity any) {
			persisted := enrollment
			if findCalls > 0 {
				persisted.Status = fun.EnrollmentStatusConfirmed
			}
			entityValue, ok := entity.(*fun.Enrollment)
			Expect(ok).To(BeTrue())
			*entityValue = persisted
			findCalls++
		}).Return(nil).Twice()
		enrollmentDAO.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()

		enrollmentPublisher := publisher.NewEnrollmentPublisher(publisher.NewBasePublisher(recorder))
		seatPublisher := publisher.NewSeatAllocationPublisher(publisher.NewBasePublisher(recorder))
		seatManager := manager.NewSeatManager(seatPublisher)
		enrollmentManager := manager.NewEnrollmentManager(nil, enrollmentDAO, enrollmentPublisher, seatManager)
		var err error
		server, err = handlers.NewMessagingServer(logger, recorder, chainChannel,
			handlers.NewEnrollmentMessageHandler(enrollmentManager),
			handlers.NewSeatMessageHandler(seatManager, enrollmentManager),
		)
		Expect(err).ToNot(HaveOccurred())
		routerCtx, routerCancel = context.WithCancel(context.Background())
		go func() { _ = server.Router().Run(routerCtx) }()
		<-server.Router().Running()

		Expect(enrollmentPublisher.Enroll(context.Background(), enrollment)).ToNot(HaveOccurred())
		chainTopics := []string{fun.TopicEnrollCmd, fun.TopicAllocateSeatCmd, fun.TopicSeatReservedEvt, fun.TopicEnrollmentConfirmedEvt}
		var chain []recordedPublish
		Eventually(func() []recordedPublish {
			chain = recorder.messages(chainTopics...)
			return chain
		}).Should(HaveLen(4))
		m1, m2, m3, m4 = chain[0].msg, chain[1].msg, chain[2].msg, chain[3].msg
	})

	AfterEach(func() {
		routerCancel()
		_ = server.Router().Close()
		_ = chainChannel.Close()
	})

	It("preserves correlation and immediate-message causation across the saga", func() {
		// M1: every message has a distinct transport identity.
		Expect(m1.UUID).ToNot(BeEmpty())
		Expect(m2.UUID).ToNot(BeEmpty())
		Expect(m3.UUID).ToNot(BeEmpty())
		Expect(m4.UUID).ToNot(BeEmpty())
		Expect(m1.UUID).ToNot(Equal(m2.UUID))
		Expect(m1.UUID).ToNot(Equal(m3.UUID))
		Expect(m1.UUID).ToNot(Equal(m4.UUID))
		Expect(m2.UUID).ToNot(Equal(m3.UUID))
		Expect(m2.UUID).ToNot(Equal(m4.UUID))
		Expect(m3.UUID).ToNot(Equal(m4.UUID))

		// M2: correlation is propagated across the complete chain.
		correlationID := middleware.MessageCorrelationID(m1)
		Expect(correlationID).NotTo(BeEmpty())
		Expect(middleware.MessageCorrelationID(m2)).To(Equal(correlationID))
		Expect(middleware.MessageCorrelationID(m3)).To(Equal(correlationID))
		Expect(middleware.MessageCorrelationID(m4)).To(Equal(correlationID))

		// M3: the first downstream message records its immediate predecessor.
		Expect(m2.Metadata.Get("causation_id")).To(Equal(m1.UUID))

		// M4: the terminal message records the preceding event.
		Expect(m3.Metadata.Get("causation_id")).To(Equal(m2.UUID))
		Expect(m4.Metadata.Get("causation_id")).To(Equal(m3.UUID))
	})
})
