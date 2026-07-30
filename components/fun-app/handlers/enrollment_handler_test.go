//nolint:dupl
package handlers_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	common "github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var _ = Describe("Enrollments", func() {
	var (
		ctx               context.Context
		db                *gorm.DB
		dbSQL             *sql.DB
		channel           *gochannel.GoChannel
		router            *gin.Engine
		person            fun.Person
		enrollmentDao     dao.EnrollmentDaoInterface
		enrollmentManager manager.EnrollmentManagerInterface
	)

	BeforeEach(func() {
		ctx = context.Background()
		channel = gochannel.NewGoChannel(gochannel.Config{}, watermill.NewStdLogger(false, false))

		var err error
		db, err = util.CreateTestDb(gormlogger.Warn)
		Expect(err).ToNot(HaveOccurred())
		Expect(db.AutoMigrate(&fun.Person{}, &fun.PersonAudit{}, &fun.Enrollment{})).To(Succeed())
		dbSQL, err = db.DB()
		Expect(err).ToNot(HaveOccurred())

		baseRepository := util.NewBaseDbRepository(db)
		tracer := otel.Tracer("fun-app-handler-test")
		personManager := manager.NewPersonManager(dao.NewPersonDao(baseRepository), tracer)
		enrollmentDao = dao.NewEnrollmentDao(baseRepository)
		enrollmentPublisher := publisher.NewEnrollmentPublisher(publisher.NewBasePublisher(channel))
		seatManager := manager.NewSeatManager(publisher.NewSeatAllocationPublisher(publisher.NewBasePublisher(channel)))
		enrollmentManager = manager.NewEnrollmentManager(
			personManager,
			enrollmentDao,
			enrollmentPublisher,
			seatManager,
		)

		person, err = personManager.CreatePerson(ctx, fun.PersonRequest{
			Name:   "REST Benchmark Person",
			Age:    10,
			Gender: "MALE",
		})
		Expect(err).ToNot(HaveOccurred())

		router = util.CreateTestGinRouter()
		router.POST("/v1/enrollments", handlers.NewEnrollmentHandler(enrollmentManager, tracer).CreateEnrollment)
	})

	AfterEach(func() {
		Expect(dbSQL.Close()).To(Succeed())
		Expect(channel.Close()).To(Succeed())
	})

	Context("POST /v1/enrollments", func() {
		var (
			request          fun.EnrollmentRequest
			response         fun.Enrollment
			persisted        fun.Enrollment
			persistedErr     common.HttpError
			command          fun.EnrollCmdV1
			commandMessage   *message.Message
			responseRecorder *httptest.ResponseRecorder
			commandMessages  <-chan *message.Message
		)

		Context("Happy Path", func() {
			Context("with existing person", func() {
				BeforeEach(func() {
					request = fun.EnrollmentRequest{PersonID: person.Id, Grade: 4}
					var err error
					commandMessages, err = channel.Subscribe(ctx, fun.TopicEnrollCmd)
					Expect(err).ToNot(HaveOccurred())

					req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", request)
					responseRecorder = recorder
					router.ServeHTTP(responseRecorder, req)
					Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
					persisted, err = enrollmentManager.GetEnrollment(ctx, person.Id)
					Expect(err).ToNot(HaveOccurred())
					Eventually(commandMessages, time.Second).Should(Receive(&commandMessage))
					Expect(json.Unmarshal(commandMessage.Payload, &command)).To(Succeed())
				})

				It("accepts, persists, and publishes the initial enrollment command", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusAccepted))
					Expect(response.ID).To(Equal(persisted.ID))
					Expect(response.PersonID).To(Equal(person.Id))
					Expect(response.Grade).To(Equal(request.Grade))
					Expect(response.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
					Expect(responseRecorder.Header().Get("Location")).To(Equal("/v1/enrollments/" + person.Id))
					Expect(persisted.PersonID).To(Equal(person.Id))
					Expect(persisted.Grade).To(Equal(request.Grade))
					Expect(persisted.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
					Expect(command.EnrollmentID).To(Equal(persisted.ID))
					Expect(command.PersonID).To(Equal(persisted.PersonID))
					Expect(command.Grade).To(Equal(persisted.Grade))
					Expect(command.Status).To(Equal(persisted.Status))
					Expect(command.RequestedAt).ToNot(Equal(time.Time{}))
				})
			})
		})

		Context("Field Validations", func() {
			Context("PersonID Field", func() {
				Context("Bad Values", func() {
					Context("with a missing personId", func() {
						BeforeEach(func() {
							var err error
							commandMessages, err = channel.Subscribe(ctx, fun.TopicEnrollCmd)
							Expect(err).ToNot(HaveOccurred())
							req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", fun.EnrollmentRequest{Grade: 4})
							responseRecorder = recorder
							router.ServeHTTP(responseRecorder, req)
							persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, person.Id)
						})

						It("returns a required PersonID validation error without persistence or publication", func() {
							util.AssertError(responseRecorder, "PersonID", "required")
							Expect(persisted).To(Equal(fun.Enrollment{}))
							Expect(persistedErr).To(Equal(common.ErrNotFound))
							Consistently(commandMessages, 100*time.Millisecond).ShouldNot(Receive())
						})
					})
				})
			})

			Context("Grade Field", func() {
				Context("Allowed Values", func() {
					Context("with grade 1", func() {
						BeforeEach(func() {
							request = fun.EnrollmentRequest{PersonID: person.Id, Grade: 1}
							var err error
							commandMessages, err = channel.Subscribe(ctx, fun.TopicEnrollCmd)
							Expect(err).ToNot(HaveOccurred())
							req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", request)
							responseRecorder = recorder
							router.ServeHTTP(responseRecorder, req)
							Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
							persisted, err = enrollmentManager.GetEnrollment(ctx, person.Id)
							Expect(err).ToNot(HaveOccurred())
							Eventually(commandMessages, time.Second).Should(Receive(&commandMessage))
							Expect(json.Unmarshal(commandMessage.Payload, &command)).To(Succeed())
						})

						It("accepts and publishes grade 1", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusAccepted))
							Expect(response.PersonID).To(Equal(person.Id))
							Expect(response.Grade).To(Equal(1))
							Expect(response.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
							Expect(persisted.ID).To(Equal(response.ID))
							Expect(persisted.Grade).To(Equal(1))
							Expect(persisted.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
							Expect(command.EnrollmentID).To(Equal(persisted.ID))
							Expect(command.PersonID).To(Equal(persisted.PersonID))
							Expect(command.Grade).To(Equal(1))
							Expect(command.RequestedAt).ToNot(Equal(time.Time{}))
						})
					})

					Context("with grade 12", func() {
						BeforeEach(func() {
							request = fun.EnrollmentRequest{PersonID: person.Id, Grade: 12}
							var err error
							commandMessages, err = channel.Subscribe(ctx, fun.TopicEnrollCmd)
							Expect(err).ToNot(HaveOccurred())
							req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", request)
							responseRecorder = recorder
							router.ServeHTTP(responseRecorder, req)
							Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
							persisted, err = enrollmentManager.GetEnrollment(ctx, person.Id)
							Expect(err).ToNot(HaveOccurred())
							Eventually(commandMessages, time.Second).Should(Receive(&commandMessage))
							Expect(json.Unmarshal(commandMessage.Payload, &command)).To(Succeed())
						})

						It("accepts and publishes grade 12", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusAccepted))
							Expect(response.PersonID).To(Equal(person.Id))
							Expect(response.Grade).To(Equal(12))
							Expect(response.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
							Expect(persisted.ID).To(Equal(response.ID))
							Expect(persisted.Grade).To(Equal(12))
							Expect(persisted.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
							Expect(command.EnrollmentID).To(Equal(persisted.ID))
							Expect(command.PersonID).To(Equal(persisted.PersonID))
							Expect(command.Grade).To(Equal(12))
							Expect(command.RequestedAt).ToNot(Equal(time.Time{}))
						})
					})
				})

				Context("Bad Values", func() {
					Context("with grade 0", func() {
						BeforeEach(func() {
							var err error
							commandMessages, err = channel.Subscribe(ctx, fun.TopicEnrollCmd)
							Expect(err).ToNot(HaveOccurred())
							req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", fun.EnrollmentRequest{PersonID: person.Id, Grade: 0})
							responseRecorder = recorder
							router.ServeHTTP(responseRecorder, req)
							persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, person.Id)
						})

						It("returns a Grade validation error without persistence or publication", func() {
							util.AssertError(responseRecorder, "Grade", "required")
							Expect(persisted).To(Equal(fun.Enrollment{}))
							Expect(persistedErr).To(Equal(common.ErrNotFound))
							Consistently(commandMessages, 100*time.Millisecond).ShouldNot(Receive())
						})
					})

					Context("with grade 13", func() {
						BeforeEach(func() {
							var err error
							commandMessages, err = channel.Subscribe(ctx, fun.TopicEnrollCmd)
							Expect(err).ToNot(HaveOccurred())
							req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", fun.EnrollmentRequest{PersonID: person.Id, Grade: 13})
							responseRecorder = recorder
							router.ServeHTTP(responseRecorder, req)
							persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, person.Id)
						})

						It("returns a Grade maximum validation error without persistence or publication", func() {
							util.AssertError(responseRecorder, "Grade", "max")
							Expect(persisted).To(Equal(fun.Enrollment{}))
							Expect(persistedErr).To(Equal(common.ErrNotFound))
							Consistently(commandMessages, 100*time.Millisecond).ShouldNot(Receive())
						})
					})
				})
			})
		})

		Context("Errors", func() {
			Context("malformed JSON", func() {
				BeforeEach(func() {
					var err error
					commandMessages, err = channel.Subscribe(ctx, fun.TopicEnrollCmd)
					Expect(err).ToNot(HaveOccurred())
					req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/v1/enrollments", strings.NewReader(`{"personId":`))
					req.Header.Set("Content-Type", "application/json")
					responseRecorder = httptest.NewRecorder()
					router.ServeHTTP(responseRecorder, req)
					persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, person.Id)
				})

				It("returns HTTP 400 without persistence or publication", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
					Expect(persisted).To(Equal(fun.Enrollment{}))
					Expect(persistedErr).To(Equal(common.ErrNotFound))
					Consistently(commandMessages, 100*time.Millisecond).ShouldNot(Receive())
				})
			})

			Context("unknown person ID", func() {
				BeforeEach(func() {
					request = fun.EnrollmentRequest{PersonID: "unknown-person", Grade: 4}
					var err error
					commandMessages, err = channel.Subscribe(ctx, fun.TopicEnrollCmd)
					Expect(err).ToNot(HaveOccurred())
					req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", request)
					responseRecorder = recorder
					router.ServeHTTP(responseRecorder, req)
					persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, request.PersonID)
				})

				It("returns HTTP 404 without persistence or publication", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
					Expect(persisted).To(Equal(fun.Enrollment{}))
					Expect(persistedErr).To(Equal(common.ErrNotFound))
					Consistently(commandMessages, 100*time.Millisecond).ShouldNot(Receive())
				})
			})
		})
	})

	Context("EnrollCmdV1", func() {
		Context("Happy Path", func() {
			Context("with an initiated enrollment", func() {
				var (
					allocationMsgs <-chan *message.Message
					allocationCmd  fun.AllocateSeatCmdV1
					enrollment     fun.Enrollment
					persisted      fun.Enrollment
				)

				BeforeEach(func() {
					enrollment = fun.Enrollment{
						PersonID: person.Id,
						Grade:    4,
						Status:   fun.EnrollmentStatusSeatAllocationInitiated,
					}
					Expect(enrollmentDao.Create(ctx, &enrollment)).ToNot(HaveOccurred())

					var err error
					allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
					Expect(err).ToNot(HaveOccurred())

					enrollmentHandler := handlers.NewEnrollmentMessageHandler(enrollmentManager)

					payload, marshalErr := json.Marshal(fun.EnrollCmdV1{
						EnrollmentID: enrollment.ID,
						PersonID:     enrollment.PersonID,
						Grade:        enrollment.Grade,
						Status:       enrollment.Status,
						RequestedAt:  time.Now().UTC(),
					})
					Expect(marshalErr).ToNot(HaveOccurred())
					Expect(enrollmentHandler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))).To(Succeed())
					var allocationMessage *message.Message
					Eventually(allocationMsgs, time.Second).Should(Receive(&allocationMessage))
					Expect(json.Unmarshal(allocationMessage.Payload, &allocationCmd)).To(Succeed())
					persisted, err = enrollmentManager.GetEnrollment(ctx, person.Id)
					Expect(err).ToNot(HaveOccurred())
				})

				It("publishes allocation for the initiated enrollment", func() {
					Expect(allocationCmd.EnrollmentID).To(Equal(enrollment.ID))
					Expect(allocationCmd.PersonID).To(Equal(person.Id))
					Expect(allocationCmd.Grade).To(Equal(4))
					Expect(allocationCmd.RequestedAt).ToNot(Equal(time.Time{}))
					Expect(persisted.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
				})
			})
		})
	})
})
