//nolint:dupl
package handlers_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
		seatManager       manager.SeatManagerInterface
		seatHandler       handlers.SeatMessageHandler
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
		seatManager = manager.NewSeatManager(publisher.NewSeatAllocationPublisher(publisher.NewBasePublisher(channel)))
		enrollmentManager = manager.NewEnrollmentManager(
			personManager,
			enrollmentDao,
			enrollmentPublisher,
			seatManager,
		)
		seatHandler = handlers.NewSeatMessageHandler(seatManager, enrollmentManager)

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
		var (
			allocationCmd  fun.AllocateSeatCmdV1
			allocationMsgs <-chan *message.Message
			command        fun.EnrollCmdV1
			enrollment     fun.Enrollment
			handler        handlers.EnrollmentMessageHandler
			resultErr      error
		)

		BeforeEach(func() {
			enrollment = fun.Enrollment{
				PersonID: person.Id,
				Grade:    4,
				Status:   fun.EnrollmentStatusSeatAllocationInitiated,
			}
			Expect(enrollmentDao.Create(ctx, &enrollment)).ToNot(HaveOccurred())
			command = fun.EnrollCmdV1{
				EnrollmentID: enrollment.ID,
				PersonID:     enrollment.PersonID,
				Grade:        enrollment.Grade,
				Status:       enrollment.Status,
				RequestedAt:  time.Now().UTC(),
			}
			handler = handlers.NewEnrollmentMessageHandler(enrollmentManager)
		})

		Context("Happy Path", func() {
			BeforeEach(func() {
				var err error
				allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
				Expect(err).ToNot(HaveOccurred())
				payload, marshalErr := json.Marshal(command)
				Expect(marshalErr).ToNot(HaveOccurred())
				resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))
				var allocationMessage *message.Message
				Eventually(allocationMsgs, time.Second).Should(Receive(&allocationMessage))
				Expect(json.Unmarshal(allocationMessage.Payload, &allocationCmd)).To(Succeed())
			})

			It("publishes allocation while preserving enrollment identity and grade", func() {
				Expect(resultErr).ToNot(HaveOccurred())
				Expect(allocationCmd.EnrollmentID).To(Equal(enrollment.ID))
				Expect(allocationCmd.PersonID).To(Equal(enrollment.PersonID))
				Expect(allocationCmd.Grade).To(Equal(enrollment.Grade))
				Expect(allocationCmd.RequestedAt).ToNot(Equal(time.Time{}))
			})
		})

		Context("Field Validations", func() {
			Context("EnrollmentID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						var err error
						allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
						Expect(err).ToNot(HaveOccurred())
						payload, marshalErr := json.Marshal(command)
						Expect(marshalErr).ToNot(HaveOccurred())
						resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))
					})

					It("accepts the enrollment ID and publishes allocation", func() {
						Expect(resultErr).ToNot(HaveOccurred())
						Eventually(allocationMsgs, time.Second).Should(Receive())
					})
				})

				Context("Bad Values", func() {
					BeforeEach(func() {
						command.EnrollmentID = ""
						var err error
						allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
						Expect(err).ToNot(HaveOccurred())
						payload, marshalErr := json.Marshal(command)
						Expect(marshalErr).ToNot(HaveOccurred())
						resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))
					})

					It("returns an EnrollmentID validation error without allocation publication", func() {
						var fieldErr common.FieldHttpError
						ok := errors.As(resultErr, &fieldErr)
						Expect(ok).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("EnrollmentID"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
						Consistently(allocationMsgs, 100*time.Millisecond).ShouldNot(Receive())
					})
				})
			})

			Context("PersonID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						var err error
						allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
						Expect(err).ToNot(HaveOccurred())
						payload, marshalErr := json.Marshal(command)
						Expect(marshalErr).ToNot(HaveOccurred())
						resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))
					})

					It("accepts the person ID and publishes allocation", func() {
						Expect(resultErr).ToNot(HaveOccurred())
						Eventually(allocationMsgs, time.Second).Should(Receive())
					})
				})

				Context("Bad Values", func() {
					BeforeEach(func() {
						command.PersonID = ""
						var err error
						allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
						Expect(err).ToNot(HaveOccurred())
						payload, marshalErr := json.Marshal(command)
						Expect(marshalErr).ToNot(HaveOccurred())
						resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))
					})

					It("returns a PersonID validation error without allocation publication", func() {
						var fieldErr common.FieldHttpError
						ok := errors.As(resultErr, &fieldErr)
						Expect(ok).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("PersonID"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
						Consistently(allocationMsgs, 100*time.Millisecond).ShouldNot(Receive())
					})
				})
			})

			Context("Grade Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						var err error
						allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
						Expect(err).ToNot(HaveOccurred())
						payload, marshalErr := json.Marshal(command)
						Expect(marshalErr).ToNot(HaveOccurred())
						resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))
					})

					It("accepts the grade and publishes allocation", func() {
						Expect(resultErr).ToNot(HaveOccurred())
						Eventually(allocationMsgs, time.Second).Should(Receive())
					})
				})

				Context("Bad Values", func() {
					BeforeEach(func() {
						command.Grade = 0
						var err error
						allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
						Expect(err).ToNot(HaveOccurred())
						payload, marshalErr := json.Marshal(command)
						Expect(marshalErr).ToNot(HaveOccurred())
						resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))
					})

					It("returns a Grade validation error without allocation publication", func() {
						var fieldErr common.FieldHttpError
						ok := errors.As(resultErr, &fieldErr)
						Expect(ok).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("Grade"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
						Consistently(allocationMsgs, 100*time.Millisecond).ShouldNot(Receive())
					})
				})
			})

			Context("Status Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						var err error
						allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
						Expect(err).ToNot(HaveOccurred())
						payload, marshalErr := json.Marshal(command)
						Expect(marshalErr).ToNot(HaveOccurred())
						resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))
					})

					It("accepts the status and publishes allocation", func() {
						Expect(resultErr).ToNot(HaveOccurred())
						Eventually(allocationMsgs, time.Second).Should(Receive())
					})
				})

				Context("Bad Values", func() {
					BeforeEach(func() {
						command.Status = "INVALID_STATUS"
						var err error
						allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
						Expect(err).ToNot(HaveOccurred())
						payload, marshalErr := json.Marshal(command)
						Expect(marshalErr).ToNot(HaveOccurred())
						resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))
					})

					It("returns a Status validation error without allocation publication", func() {
						var fieldErr common.FieldHttpError
						ok := errors.As(resultErr, &fieldErr)
						Expect(ok).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("Status"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
						Consistently(allocationMsgs, 100*time.Millisecond).ShouldNot(Receive())
					})
				})
			})

			Context("RequestedAt Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						var err error
						allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
						Expect(err).ToNot(HaveOccurred())
						payload, marshalErr := json.Marshal(command)
						Expect(marshalErr).ToNot(HaveOccurred())
						resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))
					})

					It("accepts the requested timestamp and publishes allocation", func() {
						Expect(resultErr).ToNot(HaveOccurred())
						Eventually(allocationMsgs, time.Second).Should(Receive())
					})
				})

				Context("Bad Values", func() {
					BeforeEach(func() {
						command.RequestedAt = time.Time{}
						var err error
						allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
						Expect(err).ToNot(HaveOccurred())
						payload, marshalErr := json.Marshal(command)
						Expect(marshalErr).ToNot(HaveOccurred())
						resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))
					})

					It("returns a RequestedAt validation error without allocation publication", func() {
						var fieldErr common.FieldHttpError
						ok := errors.As(resultErr, &fieldErr)
						Expect(ok).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("RequestedAt"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
						Consistently(allocationMsgs, 100*time.Millisecond).ShouldNot(Receive())
					})
				})
			})
		})
	})

	Context("AllocateSeatCmdV1", func() {
		var (
			command          fun.AllocateSeatCmdV1
			reservedEvent    fun.SeatReservedEvtV1
			reservedMessages <-chan *message.Message
			waitlistedMsgs   <-chan *message.Message
			resultErr        error
		)

		BeforeEach(func() {
			command = fun.AllocateSeatCmdV1{
				EnrollmentID: "enrollment-1",
				PersonID:     person.Id,
				Grade:        4,
				RequestedAt:  time.Now().UTC(),
			}
			var err error
			reservedMessages, err = channel.Subscribe(ctx, fun.TopicSeatReservedEvt)
			Expect(err).ToNot(HaveOccurred())
			waitlistedMsgs, err = channel.Subscribe(ctx, fun.TopicSeatWaitlistedEvt)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Happy Path", func() {
			BeforeEach(func() {
				payload, err := json.Marshal(command)
				Expect(err).ToNot(HaveOccurred())
				resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), payload))
				var reservedMessage *message.Message
				Eventually(reservedMessages, time.Second).Should(Receive(&reservedMessage))
				Expect(json.Unmarshal(reservedMessage.Payload, &reservedEvent)).To(Succeed())
			})

			It("publishes a complete reserved event preserving identity and grade", func() {
				Expect(resultErr).ToNot(HaveOccurred())
				Expect(reservedEvent.EnrollmentID).To(Equal(command.EnrollmentID))
				Expect(reservedEvent.PersonID).To(Equal(command.PersonID))
				Expect(reservedEvent.Grade).To(Equal(command.Grade))
				Expect(reservedEvent.ReservedAt).ToNot(Equal(time.Time{}))
			})
		})

		Context("Field Validations", func() {
			Context("EnrollmentID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						payload, err := json.Marshal(command)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("accepts the enrollment ID", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() {
						command.EnrollmentID = ""
						payload, err := json.Marshal(command)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("returns HTTP 400 without reserved or waitlisted publication", func() {
						var fieldErr common.FieldHttpError
						Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("EnrollmentID"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
						Consistently(reservedMessages, 100*time.Millisecond).ShouldNot(Receive())
						Consistently(waitlistedMsgs, 100*time.Millisecond).ShouldNot(Receive())
					})
				})
			})

			Context("PersonID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						payload, err := json.Marshal(command)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("accepts the person ID", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() {
						command.PersonID = ""
						payload, err := json.Marshal(command)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("returns HTTP 400 without reserved or waitlisted publication", func() {
						var fieldErr common.FieldHttpError
						Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("PersonID"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
						Consistently(reservedMessages, 100*time.Millisecond).ShouldNot(Receive())
						Consistently(waitlistedMsgs, 100*time.Millisecond).ShouldNot(Receive())
					})
				})
			})

			Context("Grade Field", func() {
				Context("Allowed Values", func() {
					Context("below reservation threshold", func() {
						BeforeEach(func() {
							command.Grade = 4
							payload, err := json.Marshal(command)
							Expect(err).ToNot(HaveOccurred())
							resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), payload))
						})
						It("publishes a reserved event", func() {
							Expect(resultErr).ToNot(HaveOccurred())
							Eventually(reservedMessages, time.Second).Should(Receive())
						})
					})
					Context("at the waitlisting threshold", func() {
						BeforeEach(func() {
							command.Grade = 5
							payload, err := json.Marshal(command)
							Expect(err).ToNot(HaveOccurred())
							resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), payload))
						})
						It("publishes a waitlisted event", func() {
							Expect(resultErr).ToNot(HaveOccurred())
							Eventually(waitlistedMsgs, time.Second).Should(Receive())
						})
					})
				})
				Context("Bad Values", func() {
					Context("with grade 0", func() {
						BeforeEach(func() {
							command.Grade = 0
							payload, err := json.Marshal(command)
							Expect(err).ToNot(HaveOccurred())
							resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), payload))
						})
						It("returns HTTP 400 without reserved or waitlisted publication", func() {
							var fieldErr common.FieldHttpError
							Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
							Expect(fieldErr.Field()).To(Equal("Grade"))
							Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
							Consistently(reservedMessages, 100*time.Millisecond).ShouldNot(Receive())
							Consistently(waitlistedMsgs, 100*time.Millisecond).ShouldNot(Receive())
						})
					})
					Context("with grade 13", func() {
						BeforeEach(func() {
							command.Grade = 13
							payload, err := json.Marshal(command)
							Expect(err).ToNot(HaveOccurred())
							resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), payload))
						})
						It("returns HTTP 400 without reserved or waitlisted publication", func() {
							var fieldErr common.FieldHttpError
							Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
							Expect(fieldErr.Field()).To(Equal("Grade"))
							Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
							Consistently(reservedMessages, 100*time.Millisecond).ShouldNot(Receive())
							Consistently(waitlistedMsgs, 100*time.Millisecond).ShouldNot(Receive())
						})
					})
				})
			})

			Context("RequestedAt Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						payload, err := json.Marshal(command)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("accepts the requested timestamp", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() {
						command.RequestedAt = time.Time{}
						payload, err := json.Marshal(command)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("returns HTTP 400 without reserved or waitlisted publication", func() {
						var fieldErr common.FieldHttpError
						Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("RequestedAt"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
						Consistently(reservedMessages, 100*time.Millisecond).ShouldNot(Receive())
						Consistently(waitlistedMsgs, 100*time.Millisecond).ShouldNot(Receive())
					})
				})
			})
		})
	})

	Context("SeatReservedEvtV1", func() {
		var (
			enrollment       fun.Enrollment
			event            fun.SeatReservedEvtV1
			confirmation     fun.EnrollmentConfirmedEvtV1
			confirmationMsgs <-chan *message.Message
			resultErr        error
			persisted        fun.Enrollment
			persistedErr     common.HttpError
		)

		BeforeEach(func() {
			enrollment = fun.Enrollment{
				PersonID: person.Id,
				Grade:    4,
				Status:   fun.EnrollmentStatusSeatAllocationInitiated,
			}
			Expect(enrollmentDao.Create(ctx, &enrollment)).ToNot(HaveOccurred())
			event = fun.SeatReservedEvtV1{
				EnrollmentID: enrollment.ID,
				PersonID:     enrollment.PersonID,
				Grade:        enrollment.Grade,
				ReservedAt:   time.Now().UTC(),
			}
			var err error
			confirmationMsgs, err = channel.Subscribe(ctx, fun.TopicEnrollmentConfirmedEvt)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Happy Path", func() {
			BeforeEach(func() {
				payload, err := json.Marshal(event)
				Expect(err).ToNot(HaveOccurred())
				resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
				persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, enrollment.PersonID)
				var confirmationMessage *message.Message
				Eventually(confirmationMsgs, time.Second).Should(Receive(&confirmationMessage))
				Expect(json.Unmarshal(confirmationMessage.Payload, &confirmation)).To(Succeed())
			})

			It("persists confirmation and publishes a complete enrollment event", func() {
				Expect(resultErr).ToNot(HaveOccurred())
				Expect(persistedErr).ToNot(HaveOccurred())
				Expect(persisted.Status).To(Equal(fun.EnrollmentStatusConfirmed))
				Expect(confirmation.EnrollmentID).To(Equal(event.EnrollmentID))
				Expect(confirmation.PersonID).To(Equal(event.PersonID))
				Expect(confirmation.ConfirmedAt).ToNot(Equal(time.Time{}))
			})
		})

		Context("Field Validations", func() {
			Context("EnrollmentID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						payload, err := json.Marshal(event)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("accepts the enrollment ID", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() {
						event.EnrollmentID = ""
						payload, err := json.Marshal(event)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("returns HTTP 400 without confirmation publication", func() {
						var fieldErr common.FieldHttpError
						Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("EnrollmentID"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
						Consistently(confirmationMsgs, 100*time.Millisecond).ShouldNot(Receive())
					})
				})
			})

			Context("PersonID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						payload, err := json.Marshal(event)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("accepts the person ID", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() {
						event.PersonID = ""
						payload, err := json.Marshal(event)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("returns HTTP 400 without confirmation publication", func() {
						var fieldErr common.FieldHttpError
						Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("PersonID"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
						Consistently(confirmationMsgs, 100*time.Millisecond).ShouldNot(Receive())
					})
				})
			})

			Context("Grade Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						payload, err := json.Marshal(event)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("accepts the grade", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					Context("with grade 0", func() {
						BeforeEach(func() {
							event.Grade = 0
							payload, err := json.Marshal(event)
							Expect(err).ToNot(HaveOccurred())
							resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
						})
						It("returns HTTP 400 without confirmation publication", func() {
							var fieldErr common.FieldHttpError
							Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
							Expect(fieldErr.Field()).To(Equal("Grade"))
							Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
							Consistently(confirmationMsgs, 100*time.Millisecond).ShouldNot(Receive())
						})
					})
					Context("with grade 13", func() {
						BeforeEach(func() {
							event.Grade = 13
							payload, err := json.Marshal(event)
							Expect(err).ToNot(HaveOccurred())
							resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
						})
						It("returns HTTP 400 without confirmation publication", func() {
							var fieldErr common.FieldHttpError
							Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
							Expect(fieldErr.Field()).To(Equal("Grade"))
							Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
							Consistently(confirmationMsgs, 100*time.Millisecond).ShouldNot(Receive())
						})
					})
				})
			})

			Context("ReservedAt Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						payload, err := json.Marshal(event)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("accepts the reserved timestamp", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() {
						event.ReservedAt = time.Time{}
						payload, err := json.Marshal(event)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("returns HTTP 400 without confirmation publication", func() {
						var fieldErr common.FieldHttpError
						Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("ReservedAt"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
						Consistently(confirmationMsgs, 100*time.Millisecond).ShouldNot(Receive())
					})
				})
			})
		})

		Context("Errors", func() {
			BeforeEach(func() {
				event.EnrollmentID = "unknown-enrollment"
				payload, err := json.Marshal(event)
				Expect(err).ToNot(HaveOccurred())
				resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
				persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, enrollment.PersonID)
			})

			It("returns not found without changing enrollment or publishing confirmation", func() {
				Expect(resultErr).To(Equal(common.ErrNotFound))
				Expect(persistedErr).ToNot(HaveOccurred())
				Expect(persisted.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
				Consistently(confirmationMsgs, 100*time.Millisecond).ShouldNot(Receive())
			})
		})
	})
})
