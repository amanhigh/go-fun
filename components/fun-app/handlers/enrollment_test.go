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

const unknownEnrollmentID = "unknown-enrollment"

var _ = Describe("Enrollments", func() {
	var (
		ctx               context.Context
		db                *gorm.DB
		dbSQL             *sql.DB
		channel           *gochannel.GoChannel
		router            *gin.Engine
		student            fun.Student
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
		Expect(db.AutoMigrate(&fun.Student{}, &fun.StudentAudit{}, &fun.Enrollment{})).To(Succeed())
		dbSQL, err = db.DB()
		Expect(err).ToNot(HaveOccurred())

		baseRepository := util.NewBaseDbRepository(db)
		tracer := otel.Tracer("fun-app-handler-test")
		studentManager := manager.NewStudentManager(dao.NewStudentDao(baseRepository), tracer)
		enrollmentDao = dao.NewEnrollmentDao(baseRepository)
		enrollmentPublisher := publisher.NewEnrollmentPublisher(publisher.NewBasePublisher(channel))
		seatManager = manager.NewSeatManager(publisher.NewSeatAllocationPublisher(publisher.NewBasePublisher(channel)))
		enrollmentManager = manager.NewEnrollmentManager(
			studentManager,
			enrollmentDao,
			enrollmentPublisher,
			seatManager,
		)
		seatHandler = handlers.NewSeatMessageHandler(seatManager, enrollmentManager)

		student, err = studentManager.CreateStudent(ctx, fun.StudentRequest{
			Name:   "REST Benchmark Student",
			Age:    10,
			Gender: "MALE",
		})
		Expect(err).ToNot(HaveOccurred())

		router = util.CreateTestGinRouter()
		router.POST("/v1/enrollments", handlers.NewEnrollmentHandler(enrollmentManager, tracer).CreateEnrollment)
		router.GET("/v1/enrollments/:studentId", handlers.NewEnrollmentHandler(enrollmentManager, tracer).GetEnrollment)
	})

	AfterEach(func() {
		Expect(dbSQL.Close()).To(Succeed())
		Expect(channel.Close()).To(Succeed())
	})

	Describe("POST /v1/enrollments", func() {
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
			Context("with existing student", func() {
				BeforeEach(func() {
					request = fun.EnrollmentRequest{StudentID: student.Id, Grade: 4}
					var err error
					commandMessages, err = channel.Subscribe(ctx, fun.TopicEnrollCmd)
					Expect(err).ToNot(HaveOccurred())

					req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", request)
					responseRecorder = recorder
					router.ServeHTTP(responseRecorder, req)
					Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
					persisted, err = enrollmentManager.GetEnrollment(ctx, student.Id)
					Expect(err).ToNot(HaveOccurred())
					Eventually(commandMessages, time.Second).Should(Receive(&commandMessage))
					Expect(json.Unmarshal(commandMessage.Payload, &command)).To(Succeed())
				})

				It("accepts, persists, and publishes the initial enrollment command", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusAccepted))
					Expect(response.ID).To(Equal(persisted.ID))
					Expect(response.StudentID).To(Equal(student.Id))
					Expect(response.Grade).To(Equal(request.Grade))
					Expect(response.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
					Expect(responseRecorder.Header().Get("Location")).To(Equal("/v1/enrollments/" + student.Id))
					Expect(persisted.StudentID).To(Equal(student.Id))
					Expect(persisted.Grade).To(Equal(request.Grade))
					Expect(persisted.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
					Expect(command.EnrollmentID).To(Equal(persisted.ID))
					Expect(command.StudentID).To(Equal(persisted.StudentID))
					Expect(command.Grade).To(Equal(persisted.Grade))
					Expect(command.RequestedAt).ToNot(Equal(time.Time{}))
				})
			})
		})

		Context("Field Validations", func() {
			Context("StudentID Field", func() {
				Context("Bad Values", func() {
					Context("with a missing studentId", func() {
						BeforeEach(func() {
							var err error
							commandMessages, err = channel.Subscribe(ctx, fun.TopicEnrollCmd)
							Expect(err).ToNot(HaveOccurred())
							req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", fun.EnrollmentRequest{Grade: 4})
							responseRecorder = recorder
							router.ServeHTTP(responseRecorder, req)
							persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, student.Id)
						})

						It("returns a required StudentID validation error without persistence or publication", func() {
							util.AssertError(responseRecorder, "StudentID", "required")
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
							request = fun.EnrollmentRequest{StudentID: student.Id, Grade: 1}
							var err error
							commandMessages, err = channel.Subscribe(ctx, fun.TopicEnrollCmd)
							Expect(err).ToNot(HaveOccurred())
							req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", request)
							responseRecorder = recorder
							router.ServeHTTP(responseRecorder, req)
							Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
							persisted, err = enrollmentManager.GetEnrollment(ctx, student.Id)
							Expect(err).ToNot(HaveOccurred())
							Eventually(commandMessages, time.Second).Should(Receive(&commandMessage))
							Expect(json.Unmarshal(commandMessage.Payload, &command)).To(Succeed())
						})

						It("accepts and publishes grade 1", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusAccepted))
							Expect(response.StudentID).To(Equal(student.Id))
							Expect(response.Grade).To(Equal(1))
							Expect(response.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
							Expect(persisted.ID).To(Equal(response.ID))
							Expect(persisted.Grade).To(Equal(1))
							Expect(persisted.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
							Expect(command.EnrollmentID).To(Equal(persisted.ID))
							Expect(command.StudentID).To(Equal(persisted.StudentID))
							Expect(command.Grade).To(Equal(1))
							Expect(command.RequestedAt).ToNot(Equal(time.Time{}))
						})
					})

					Context("with grade 12", func() {
						BeforeEach(func() {
							request = fun.EnrollmentRequest{StudentID: student.Id, Grade: 12}
							var err error
							commandMessages, err = channel.Subscribe(ctx, fun.TopicEnrollCmd)
							Expect(err).ToNot(HaveOccurred())
							req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", request)
							responseRecorder = recorder
							router.ServeHTTP(responseRecorder, req)
							Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
							persisted, err = enrollmentManager.GetEnrollment(ctx, student.Id)
							Expect(err).ToNot(HaveOccurred())
							Eventually(commandMessages, time.Second).Should(Receive(&commandMessage))
							Expect(json.Unmarshal(commandMessage.Payload, &command)).To(Succeed())
						})

						It("accepts and publishes grade 12", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusAccepted))
							Expect(response.StudentID).To(Equal(student.Id))
							Expect(response.Grade).To(Equal(12))
							Expect(response.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
							Expect(persisted.ID).To(Equal(response.ID))
							Expect(persisted.Grade).To(Equal(12))
							Expect(persisted.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
							Expect(command.EnrollmentID).To(Equal(persisted.ID))
							Expect(command.StudentID).To(Equal(persisted.StudentID))
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
							req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", fun.EnrollmentRequest{StudentID: student.Id, Grade: 0})
							responseRecorder = recorder
							router.ServeHTTP(responseRecorder, req)
							persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, student.Id)
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
							req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", fun.EnrollmentRequest{StudentID: student.Id, Grade: 13})
							responseRecorder = recorder
							router.ServeHTTP(responseRecorder, req)
							persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, student.Id)
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
					req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/v1/enrollments", strings.NewReader(`{"studentId":`))
					req.Header.Set("Content-Type", "application/json")
					responseRecorder = httptest.NewRecorder()
					router.ServeHTTP(responseRecorder, req)
					persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, student.Id)
				})

				It("returns HTTP 400 without persistence or publication", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
					Expect(persisted).To(Equal(fun.Enrollment{}))
					Expect(persistedErr).To(Equal(common.ErrNotFound))
					Consistently(commandMessages, 100*time.Millisecond).ShouldNot(Receive())
				})
			})

			Context("unknown student ID", func() {
				BeforeEach(func() {
					request = fun.EnrollmentRequest{StudentID: "unknown-student", Grade: 4}
					var err error
					commandMessages, err = channel.Subscribe(ctx, fun.TopicEnrollCmd)
					Expect(err).ToNot(HaveOccurred())
					req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/enrollments", request)
					responseRecorder = recorder
					router.ServeHTTP(responseRecorder, req)
					persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, request.StudentID)
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

	Describe("GET /v1/enrollments/:studentId", func() {
		var (
			enrollment       fun.Enrollment
			response         fun.Enrollment
			responseRecorder *httptest.ResponseRecorder
		)

		Context("Happy Path", func() {
			BeforeEach(func() {
				enrollment = fun.Enrollment{
					StudentID: student.Id,
					Grade:    4,
					Status:   fun.EnrollmentStatusSeatAllocationInitiated,
				}
				Expect(enrollmentDao.Create(ctx, &enrollment)).ToNot(HaveOccurred())
				req, recorder := util.CreateTestRequest(http.MethodGet, "/v1/enrollments/"+student.Id, nil)
				responseRecorder = recorder
				router.ServeHTTP(responseRecorder, req)
				Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
			})

			It("returns the persisted enrollment for the student", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
				Expect(response.ID).To(Equal(enrollment.ID))
				Expect(response.StudentID).To(Equal(enrollment.StudentID))
				Expect(response.Grade).To(Equal(enrollment.Grade))
				Expect(response.Status).To(Equal(enrollment.Status))
				Expect(response.CreatedAt.Equal(enrollment.CreatedAt)).To(BeTrue())
				Expect(response.UpdatedAt.Equal(enrollment.UpdatedAt)).To(BeTrue())
			})
		})

		Context("Field Validations", func() {
			Context("StudentID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						enrollment = fun.Enrollment{
							StudentID: student.Id,
							Grade:    4,
							Status:   fun.EnrollmentStatusSeatAllocationInitiated,
						}
						Expect(enrollmentDao.Create(ctx, &enrollment)).ToNot(HaveOccurred())
						req, recorder := util.CreateTestRequest(http.MethodGet, "/v1/enrollments/"+student.Id, nil)
						responseRecorder = recorder
						router.ServeHTTP(responseRecorder, req)
					})

					It("accepts an existing student ID", func() {
						Expect(responseRecorder.Code).To(Equal(http.StatusOK))
					})
				})
			})
		})

		Context("Errors", func() {
			Context("unknown enrollment for an existing student", func() {
				BeforeEach(func() {
					req, recorder := util.CreateTestRequest(http.MethodGet, "/v1/enrollments/"+student.Id, nil)
					responseRecorder = recorder
					router.ServeHTTP(responseRecorder, req)
				})

				It("returns HTTP 404", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
				})
			})

		})
	})

	Describe("EnrollCmdV1", func() {
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
				StudentID: student.Id,
				Grade:    4,
				Status:   fun.EnrollmentStatusSeatAllocationInitiated,
			}
			Expect(enrollmentDao.Create(ctx, &enrollment)).ToNot(HaveOccurred())
			command = fun.EnrollCmdV1{
				EnrollmentID: enrollment.ID,
				StudentID:     enrollment.StudentID,
				Grade:        enrollment.Grade,
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
				Expect(allocationCmd.StudentID).To(Equal(enrollment.StudentID))
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

			Context("StudentID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						var err error
						allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
						Expect(err).ToNot(HaveOccurred())
						payload, marshalErr := json.Marshal(command)
						Expect(marshalErr).ToNot(HaveOccurred())
						resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))
					})

					It("accepts the student ID and publishes allocation", func() {
						Expect(resultErr).ToNot(HaveOccurred())
						Eventually(allocationMsgs, time.Second).Should(Receive())
					})
				})

				Context("Bad Values", func() {
					BeforeEach(func() {
						command.StudentID = ""
						var err error
						allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
						Expect(err).ToNot(HaveOccurred())
						payload, marshalErr := json.Marshal(command)
						Expect(marshalErr).ToNot(HaveOccurred())
						resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), payload))
					})

					It("returns a StudentID validation error without allocation publication", func() {
						var fieldErr common.FieldHttpError
						ok := errors.As(resultErr, &fieldErr)
						Expect(ok).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("StudentID"))
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

		Context("Errors", func() {
			Context("malformed payload", func() {
				BeforeEach(func() {
					var err error
					allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
					Expect(err).ToNot(HaveOccurred())
					resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), []byte(`{"enrollmentId":`)))
				})

				It("returns HTTP 400 without allocation publication", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
					Consistently(allocationMsgs, 100*time.Millisecond).ShouldNot(Receive())
				})
			})

			Context("empty payload", func() {
				BeforeEach(func() {
					var err error
					allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
					Expect(err).ToNot(HaveOccurred())
					resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), nil))
				})

				It("returns HTTP 400 without allocation publication", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
					Consistently(allocationMsgs, 100*time.Millisecond).ShouldNot(Receive())
				})
			})

			Context("JSON null", func() {
				BeforeEach(func() {
					var err error
					allocationMsgs, err = channel.Subscribe(ctx, fun.TopicAllocateSeatCmd)
					Expect(err).ToNot(HaveOccurred())
					resultErr = handler.HandleEnrollCmd(message.NewMessage(watermill.NewUUID(), []byte("null")))
				})

				It("returns HTTP 400 without allocation publication", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
					Consistently(allocationMsgs, 100*time.Millisecond).ShouldNot(Receive())
				})
			})
		})

	})

	Describe("AllocateSeatCmdV1", func() {
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
				StudentID:     student.Id,
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
				Expect(reservedEvent.StudentID).To(Equal(command.StudentID))
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

			Context("StudentID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						payload, err := json.Marshal(command)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("accepts the student ID", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() {
						command.StudentID = ""
						payload, err := json.Marshal(command)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("returns HTTP 400 without reserved or waitlisted publication", func() {
						var fieldErr common.FieldHttpError
						Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("StudentID"))
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

		Context("Errors", func() {
			Context("malformed payload", func() {
				BeforeEach(func() {
					resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), []byte(`{"enrollmentId":`)))
				})

				It("returns HTTP 400 without publication", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
					Consistently(reservedMessages, 100*time.Millisecond).ShouldNot(Receive())
					Consistently(waitlistedMsgs, 100*time.Millisecond).ShouldNot(Receive())
				})
			})

			Context("empty payload", func() {
				BeforeEach(func() {
					resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), nil))
				})

				It("returns HTTP 400 without publication", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
					Consistently(reservedMessages, 100*time.Millisecond).ShouldNot(Receive())
					Consistently(waitlistedMsgs, 100*time.Millisecond).ShouldNot(Receive())
				})
			})

			Context("JSON null", func() {
				BeforeEach(func() {
					resultErr = seatHandler.HandleAllocateSeatCmd(message.NewMessage(watermill.NewUUID(), []byte("null")))
				})

				It("returns HTTP 400 without publication", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
					Consistently(reservedMessages, 100*time.Millisecond).ShouldNot(Receive())
					Consistently(waitlistedMsgs, 100*time.Millisecond).ShouldNot(Receive())
				})
			})
		})
	})

	Describe("SeatReservedEvtV1", func() {
		var (
			enrollment   fun.Enrollment
			event        fun.SeatReservedEvtV1
			resultErr    error
			persisted    fun.Enrollment
			persistedErr common.HttpError
		)

		BeforeEach(func() {
			enrollment = fun.Enrollment{
				StudentID: student.Id,
				Grade:    4,
				Status:   fun.EnrollmentStatusSeatAllocationInitiated,
			}
			Expect(enrollmentDao.Create(ctx, &enrollment)).ToNot(HaveOccurred())
			event = fun.SeatReservedEvtV1{
				EnrollmentEvent: fun.EnrollmentEvent{
					EnrollmentID: enrollment.ID,
					StudentID:     enrollment.StudentID,
				},
				Grade:        enrollment.Grade,
				ReservedAt:   time.Now().UTC(),
			}
		})

		Context("Happy Path", func() {
			BeforeEach(func() {
				payload, err := json.Marshal(event)
				Expect(err).ToNot(HaveOccurred())
				resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
				persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, enrollment.StudentID)
			})

			It("persists the reserved enrollment without publishing another event", func() {
				Expect(resultErr).ToNot(HaveOccurred())
				Expect(persistedErr).ToNot(HaveOccurred())
				Expect(persisted.Status).To(Equal(fun.EnrollmentStatusConfirmed))
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
					It("returns HTTP 400 without publication", func() {
						var fieldErr common.FieldHttpError
						Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("EnrollmentID"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("StudentID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						payload, err := json.Marshal(event)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("accepts the student ID", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() {
						event.StudentID = ""
						payload, err := json.Marshal(event)
						Expect(err).ToNot(HaveOccurred())
						resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
					})
					It("returns HTTP 400 without publication", func() {
						var fieldErr common.FieldHttpError
						Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("StudentID"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
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
						It("returns HTTP 400 without publication", func() {
							var fieldErr common.FieldHttpError
							Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
							Expect(fieldErr.Field()).To(Equal("Grade"))
							Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
						})
					})
					Context("with grade 13", func() {
						BeforeEach(func() {
							event.Grade = 13
							payload, err := json.Marshal(event)
							Expect(err).ToNot(HaveOccurred())
							resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
						})
						It("returns HTTP 400 without publication", func() {
							var fieldErr common.FieldHttpError
							Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
							Expect(fieldErr.Field()).To(Equal("Grade"))
							Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
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
					It("returns HTTP 400 without publication", func() {
						var fieldErr common.FieldHttpError
						Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
						Expect(fieldErr.Field()).To(Equal("ReservedAt"))
						Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					})
				})
			})
		})

		Context("Errors", func() {
			Context("malformed payload", func() {
				BeforeEach(func() {
					resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), []byte(`{"enrollmentId":`)))
				})

				It("returns HTTP 400 without publication", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
				})
			})

			Context("empty payload", func() {
				BeforeEach(func() {
					resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), nil))
				})

				It("returns HTTP 400 without publication", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
				})
			})

			Context("JSON null", func() {
				BeforeEach(func() {
					resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), []byte("null")))
				})

				It("returns HTTP 400 without publication", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
				})
			})

			BeforeEach(func() {
				event.EnrollmentID = unknownEnrollmentID
				payload, err := json.Marshal(event)
				Expect(err).ToNot(HaveOccurred())
				resultErr = seatHandler.HandleSeatReservedEvt(message.NewMessage(watermill.NewUUID(), payload))
				persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, enrollment.StudentID)
			})

			It("returns not found without changing enrollment or publishing", func() {
				Expect(resultErr).To(Equal(common.ErrNotFound))
				Expect(persistedErr).ToNot(HaveOccurred())
				Expect(persisted.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
			})
		})
	})

	Describe("SeatWaitlistedEvtV1", func() {
		var (
			enrollment   fun.Enrollment
			event        fun.SeatWaitlistedEvtV1
			resultErr    error
			persisted    fun.Enrollment
			persistedErr common.HttpError
		)

		BeforeEach(func() {
			enrollment = fun.Enrollment{StudentID: student.Id, Grade: 4, Status: fun.EnrollmentStatusSeatAllocationInitiated}
			Expect(enrollmentDao.Create(ctx, &enrollment)).ToNot(HaveOccurred())
			event = fun.SeatWaitlistedEvtV1{
				EnrollmentEvent: fun.EnrollmentEvent{
					EnrollmentID: enrollment.ID,
					StudentID:     enrollment.StudentID,
				},
				Grade: enrollment.Grade,
				Reason: "capacity reached", WaitlistedAt: time.Now().UTC(),
			}
		})

		execute := func() {
			payload, err := json.Marshal(event)
			Expect(err).ToNot(HaveOccurred())
			resultErr = seatHandler.HandleSeatWaitlistedEvt(message.NewMessage(watermill.NewUUID(), payload))
			persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, enrollment.StudentID)
		}
		assertFieldError := func(field, rule string) {
			var fieldErr common.FieldHttpError
			Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
			Expect(fieldErr.Field()).To(Equal(field))
			Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
			Expect(resultErr.Error()).To(ContainSubstring(rule))
		}

		Context("Happy Path", func() {
			Context("initial enrollment", func() {
				BeforeEach(execute)
				It("transitions INITIATED to WAITLISTED", func() {
					Expect(resultErr).ToNot(HaveOccurred())
					Expect(persistedErr).ToNot(HaveOccurred())
					Expect(persisted.Status).To(Equal(fun.EnrollmentStatusWaitlisted))
				})
			})
			Context("duplicate waitlisting", func() {
				BeforeEach(func() {
					enrollment.Status = fun.EnrollmentStatusWaitlisted
					Expect(enrollmentDao.Update(ctx, &enrollment)).ToNot(HaveOccurred())
					execute()
				})
				It("acknowledges the duplicate and preserves WAITLISTED", func() {
					Expect(resultErr).ToNot(HaveOccurred())
					Expect(persisted.Status).To(Equal(fun.EnrollmentStatusWaitlisted))
				})
			})
			Context("terminal enrollment", func() {
				Context("CONFIRMED", func() {
					BeforeEach(func() {
						enrollment.Status = fun.EnrollmentStatusConfirmed
						Expect(enrollmentDao.Update(ctx, &enrollment)).ToNot(HaveOccurred())
						execute()
					})
					It("acknowledges without overwriting CONFIRMED", func() {
						Expect(resultErr).ToNot(HaveOccurred())
						Expect(persisted.Status).To(Equal(fun.EnrollmentStatusConfirmed))
					})
				})
				Context("CANCELLED", func() {
					BeforeEach(func() {
						enrollment.Status = fun.EnrollmentStatusCancelled
						Expect(enrollmentDao.Update(ctx, &enrollment)).ToNot(HaveOccurred())
						execute()
					})
					It("acknowledges without overwriting CANCELLED", func() {
						Expect(resultErr).ToNot(HaveOccurred())
						Expect(persisted.Status).To(Equal(fun.EnrollmentStatusCancelled))
					})
				})
			})
		})

		Context("Field Validations", func() {
			Context("EnrollmentID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(execute)
					It("accepts a valid enrollment ID", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() { event.EnrollmentID = ""; execute() })
					It("returns the required validation error", func() { assertFieldError("EnrollmentID", "required") })
				})
			})
			Context("StudentID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(execute)
					It("accepts a valid student ID", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() { event.StudentID = ""; execute() })
					It("returns the required validation error", func() { assertFieldError("StudentID", "required") })
				})
			})
			Context("Grade Field", func() {
				Context("Allowed Values", func() {
					Context("grade 1", func() {
						BeforeEach(func() { event.Grade = 1; execute() })
						It("is accepted", func() { Expect(resultErr).ToNot(HaveOccurred()) })
					})
					Context("grade 12", func() {
						BeforeEach(func() { event.Grade = 12; execute() })
						It("is accepted", func() { Expect(resultErr).ToNot(HaveOccurred()) })
					})
				})
				Context("Bad Values", func() {
					Context("grade 0", func() {
						BeforeEach(func() { event.Grade = 0; execute() })
						It("returns the required validation error", func() { assertFieldError("Grade", "required") })
					})
					Context("grade 13", func() {
						BeforeEach(func() { event.Grade = 13; execute() })
						It("returns the maximum validation error", func() { assertFieldError("Grade", "max") })
					})
				})
			})
			Context("Reason Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(execute)
					It("accepts a failure reason", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() { event.Reason = ""; execute() })
					It("returns the required validation error", func() { assertFieldError("Reason", "required") })
				})
			})
			Context("WaitlistedAt Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(execute)
					It("accepts a waitlisted timestamp", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() { event.WaitlistedAt = time.Time{}; execute() })
					It("returns the required validation error", func() { assertFieldError("WaitlistedAt", "required") })
				})
			})
		})

		Context("Errors", func() {
			Context("malformed payload", func() {
				BeforeEach(func() {
					resultErr = seatHandler.HandleSeatWaitlistedEvt(message.NewMessage(watermill.NewUUID(), []byte(`{"enrollmentId":`)))
				})
				It("returns HTTP 400", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
				})
			})
			Context("empty payload", func() {
				BeforeEach(func() { resultErr = seatHandler.HandleSeatWaitlistedEvt(message.NewMessage(watermill.NewUUID(), nil)) })
				It("returns HTTP 400", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
				})
			})
			Context("JSON null", func() {
				BeforeEach(func() {
					resultErr = seatHandler.HandleSeatWaitlistedEvt(message.NewMessage(watermill.NewUUID(), []byte("null")))
				})
				It("returns HTTP 400", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
				})
			})
			Context("unknown enrollment", func() {
				BeforeEach(func() { event.EnrollmentID = unknownEnrollmentID; execute() })
				It("returns not found and leaves the enrollment unchanged", func() {
					Expect(resultErr).To(Equal(common.ErrNotFound))
					Expect(persistedErr).ToNot(HaveOccurred())
					Expect(persisted.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
				})
			})
		})
	})

	Describe("SeatAllocationFailedEvtV1", func() {
		var (
			enrollment   fun.Enrollment
			event        fun.SeatAllocationFailedEvtV1
			resultErr    error
			persisted    fun.Enrollment
			persistedErr common.HttpError
		)

		BeforeEach(func() {
			enrollment = fun.Enrollment{StudentID: student.Id, Grade: 4, Status: fun.EnrollmentStatusSeatAllocationInitiated}
			Expect(enrollmentDao.Create(ctx, &enrollment)).ToNot(HaveOccurred())
			event = fun.SeatAllocationFailedEvtV1{
				EnrollmentEvent: fun.EnrollmentEvent{
					EnrollmentID: enrollment.ID,
					StudentID:     enrollment.StudentID,
				},
				Reason:   "allocation failed",
				FailedAt: time.Now().UTC(),
			}
		})

		execute := func() {
			payload, err := json.Marshal(event)
			Expect(err).ToNot(HaveOccurred())
			resultErr = seatHandler.HandleSeatAllocationFailedEvt(message.NewMessage(watermill.NewUUID(), payload))
			persisted, persistedErr = enrollmentManager.GetEnrollment(ctx, enrollment.StudentID)
		}
		assertFieldError := func(field, rule string) {
			var fieldErr common.FieldHttpError
			Expect(errors.As(resultErr, &fieldErr)).To(BeTrue())
			Expect(fieldErr.Field()).To(Equal(field))
			Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
			Expect(resultErr.Error()).To(ContainSubstring(rule))
		}

		Context("Happy Path", func() {
			Context("initial enrollment", func() {
				BeforeEach(execute)
				It("transitions INITIATED to CANCELLED", func() {
					Expect(resultErr).ToNot(HaveOccurred())
					Expect(persistedErr).ToNot(HaveOccurred())
					Expect(persisted.Status).To(Equal(fun.EnrollmentStatusCancelled))
				})
			})
			Context("duplicate cancellation", func() {
				BeforeEach(func() {
					enrollment.Status = fun.EnrollmentStatusCancelled
					Expect(enrollmentDao.Update(ctx, &enrollment)).ToNot(HaveOccurred())
					execute()
				})
				It("acknowledges the duplicate and preserves CANCELLED", func() {
					Expect(resultErr).ToNot(HaveOccurred())
					Expect(persisted.Status).To(Equal(fun.EnrollmentStatusCancelled))
				})
			})
			Context("terminal enrollment", func() {
				Context("CONFIRMED", func() {
					BeforeEach(func() {
						enrollment.Status = fun.EnrollmentStatusConfirmed
						Expect(enrollmentDao.Update(ctx, &enrollment)).ToNot(HaveOccurred())
						execute()
					})
					It("acknowledges without overwriting CONFIRMED", func() {
						Expect(resultErr).ToNot(HaveOccurred())
						Expect(persisted.Status).To(Equal(fun.EnrollmentStatusConfirmed))
					})
				})
				Context("WAITLISTED", func() {
					BeforeEach(func() {
						enrollment.Status = fun.EnrollmentStatusWaitlisted
						Expect(enrollmentDao.Update(ctx, &enrollment)).ToNot(HaveOccurred())
						execute()
					})
					It("acknowledges without overwriting WAITLISTED", func() {
						Expect(resultErr).ToNot(HaveOccurred())
						Expect(persisted.Status).To(Equal(fun.EnrollmentStatusWaitlisted))
					})
				})
			})
		})

		Context("Field Validations", func() {
			Context("EnrollmentID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(execute)
					It("accepts a valid enrollment ID", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() { event.EnrollmentID = ""; execute() })
					It("returns the required validation error", func() { assertFieldError("EnrollmentID", "required") })
				})
			})
			Context("StudentID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(execute)
					It("accepts a valid student ID", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() { event.StudentID = ""; execute() })
					It("returns the required validation error", func() { assertFieldError("StudentID", "required") })
				})
			})
			Context("Reason Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(execute)
					It("accepts a failure reason", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() { event.Reason = ""; execute() })
					It("returns the required validation error", func() { assertFieldError("Reason", "required") })
				})
			})
			Context("FailedAt Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(execute)
					It("accepts a failure timestamp", func() { Expect(resultErr).ToNot(HaveOccurred()) })
				})
				Context("Bad Values", func() {
					BeforeEach(func() { event.FailedAt = time.Time{}; execute() })
					It("returns the required validation error", func() { assertFieldError("FailedAt", "required") })
				})
			})
		})

		Context("Errors", func() {
			Context("malformed payload", func() {
				BeforeEach(func() {
					resultErr = seatHandler.HandleSeatAllocationFailedEvt(message.NewMessage(watermill.NewUUID(), []byte(`{"enrollmentId":`)))
				})
				It("returns HTTP 400", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
				})
			})
			Context("empty payload", func() {
				BeforeEach(func() {
					resultErr = seatHandler.HandleSeatAllocationFailedEvt(message.NewMessage(watermill.NewUUID(), nil))
				})
				It("returns HTTP 400", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
				})
			})
			Context("JSON null", func() {
				BeforeEach(func() {
					resultErr = seatHandler.HandleSeatAllocationFailedEvt(message.NewMessage(watermill.NewUUID(), []byte("null")))
				})
				It("returns HTTP 400", func() {
					var httpErr common.HttpError
					Expect(errors.As(resultErr, &httpErr)).To(BeTrue())
					Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
				})
			})
			Context("unknown enrollment", func() {
				BeforeEach(func() { event.EnrollmentID = unknownEnrollmentID; execute() })
				It("returns not found and leaves the enrollment unchanged", func() {
					Expect(resultErr).To(Equal(common.ErrNotFound))
					Expect(persistedErr).ToNot(HaveOccurred())
					Expect(persisted.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
				})
			})
		})
	})
})
