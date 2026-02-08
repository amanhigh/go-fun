package play_fast

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cron", func() {

	Context("gocron", func() {

		Context("Scheduler Initialization", func() {
			It("should create and shutdown scheduler gracefully", func() {
				s, err := gocron.NewScheduler()
				Expect(err).ToNot(HaveOccurred())
				Expect(s).ToNot(BeNil())

				s.Start()
				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should configure scheduler with location", func() {
				s, err := gocron.NewScheduler(
					gocron.WithLocation(time.UTC),
				)
				Expect(err).ToNot(HaveOccurred())
				defer s.Shutdown()
				Expect(s).ToNot(BeNil())
			})

			It("should configure scheduler with concurrent job limit", func() {
				s, err := gocron.NewScheduler(
					gocron.WithLimitConcurrentJobs(5, gocron.LimitModeWait),
				)
				Expect(err).ToNot(HaveOccurred())
				defer s.Shutdown()
				Expect(s).ToNot(BeNil())
			})
		})

		Context("Duration Jobs", func() {
			It("should schedule job at fixed time intervals", func() {
				var count atomic.Int32
				fakeClock := clockwork.NewFakeClock()

				s, err := gocron.NewScheduler(gocron.WithClock(fakeClock))
				Expect(err).ToNot(HaveOccurred())

				var wg sync.WaitGroup
				wg.Add(3)
				_, err = s.NewJob(
					gocron.DurationJob(5*time.Second),
					gocron.NewTask(func() {
						count.Add(1)
						wg.Done()
					}),
				)
				Expect(err).ToNot(HaveOccurred())

				s.Start()

				By("Advancing clock to trigger 3 executions")
				for i := 0; i < 3; i++ {
					err = fakeClock.BlockUntilContext(context.Background(), 1)
					Expect(err).ToNot(HaveOccurred())
					fakeClock.Advance(5 * time.Second)
				}

				wg.Wait()
				Expect(count.Load()).To(Equal(int32(3)))

				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should schedule job at random duration intervals", func() {
				s, err := gocron.NewScheduler()
				Expect(err).ToNot(HaveOccurred())

				j, err := s.NewJob(
					gocron.DurationRandomJob(1*time.Second, 5*time.Second),
					gocron.NewTask(func() {}),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(j.ID()).ToNot(BeNil())

				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Cron Expression Jobs", func() {
			It("should schedule job using standard cron expressions", func() {
				s, err := gocron.NewScheduler()
				Expect(err).ToNot(HaveOccurred())

				j, err := s.NewJob(
					gocron.CronJob("*/5 * * * *", false),
					gocron.NewTask(func() {}),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(j.ID()).ToNot(BeNil())

				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should schedule job with seconds support", func() {
				s, err := gocron.NewScheduler()
				Expect(err).ToNot(HaveOccurred())

				j, err := s.NewJob(
					gocron.CronJob("*/5 * * * * *", true),
					gocron.NewTask(func() {}),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(j.ID()).ToNot(BeNil())

				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Time-Based Jobs", func() {
			It("should schedule daily job at specific time", func() {
				s, err := gocron.NewScheduler()
				Expect(err).ToNot(HaveOccurred())

				j, err := s.NewJob(
					gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(9, 0, 0))),
					gocron.NewTask(func() {}),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(j.ID()).ToNot(BeNil())

				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should schedule weekly job on specific days", func() {
				s, err := gocron.NewScheduler()
				Expect(err).ToNot(HaveOccurred())

				j, err := s.NewJob(
					gocron.WeeklyJob(1,
						gocron.NewWeekdays(time.Monday, time.Wednesday, time.Friday),
						gocron.NewAtTimes(gocron.NewAtTime(10, 30, 0)),
					),
					gocron.NewTask(func() {}),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(j.ID()).ToNot(BeNil())

				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should schedule monthly job on specific days", func() {
				s, err := gocron.NewScheduler()
				Expect(err).ToNot(HaveOccurred())

				j, err := s.NewJob(
					gocron.MonthlyJob(1,
						gocron.NewDaysOfTheMonth(1, 15),
						gocron.NewAtTimes(gocron.NewAtTime(8, 0, 0)),
					),
					gocron.NewTask(func() {}),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(j.ID()).ToNot(BeNil())

				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should schedule one-time job at specific datetime", func() {
				s, err := gocron.NewScheduler()
				Expect(err).ToNot(HaveOccurred())

				futureTime := time.Now().Add(24 * time.Hour)
				j, err := s.NewJob(
					gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(futureTime)),
					gocron.NewTask(func() {}),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(j.ID()).ToNot(BeNil())

				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Job Concurrency Control", func() {
			It("should support singleton mode (skip overlapping)", func() {
				s, err := gocron.NewScheduler()
				Expect(err).ToNot(HaveOccurred())

				j, err := s.NewJob(
					gocron.DurationJob(time.Second),
					gocron.NewTask(func() {
						time.Sleep(2 * time.Second)
					}),
					gocron.WithSingletonMode(gocron.LimitModeReschedule),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(j.ID()).ToNot(BeNil())

				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Job Management", func() {
			var (
				s   gocron.Scheduler
				err error
			)

			BeforeEach(func() {
				s, err = gocron.NewScheduler()
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should retrieve job by ID", func() {
				j, err := s.NewJob(
					gocron.DurationJob(time.Minute),
					gocron.NewTask(func() {}),
				)
				Expect(err).ToNot(HaveOccurred())

				jobs := s.Jobs()
				Expect(jobs).To(HaveLen(1))
				Expect(jobs[0].ID()).To(Equal(j.ID()))
			})

			It("should remove scheduled jobs", func() {
				j, err := s.NewJob(
					gocron.DurationJob(time.Minute),
					gocron.NewTask(func() {}),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(s.Jobs()).To(HaveLen(1))

				err = s.RemoveJob(j.ID())
				Expect(err).ToNot(HaveOccurred())
				Expect(s.Jobs()).To(HaveLen(0))
			})

			It("should list all scheduled jobs", func() {
				for i := 0; i < 3; i++ {
					_, err = s.NewJob(
						gocron.DurationJob(time.Duration(i+1)*time.Minute),
						gocron.NewTask(func() {}),
					)
					Expect(err).ToNot(HaveOccurred())
				}

				Expect(s.Jobs()).To(HaveLen(3))
			})

			It("should update job schedule by removing and re-adding", func() {
				By("Creating initial job with 1 minute interval")
				j, err := s.NewJob(
					gocron.DurationJob(time.Minute),
					gocron.NewTask(func() {}),
				)
				Expect(err).ToNot(HaveOccurred())
				originalID := j.ID()

				By("Removing old job and creating new one with updated schedule")
				err = s.RemoveJob(originalID)
				Expect(err).ToNot(HaveOccurred())

				j2, err := s.NewJob(
					gocron.DurationJob(30*time.Second),
					gocron.NewTask(func() {}),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(j2.ID()).ToNot(Equal(originalID))
				Expect(s.Jobs()).To(HaveLen(1))
			})
		})

		Context("Error Scenarios", func() {
			It("should handle invalid cron expression", func() {
				s, err := gocron.NewScheduler()
				Expect(err).ToNot(HaveOccurred())
				defer s.Shutdown()

				_, err = s.NewJob(
					gocron.CronJob("bad cron", true),
					gocron.NewTask(func() {}),
				)
				Expect(err).To(HaveOccurred())
			})

			It("should handle nil function", func() {
				s, err := gocron.NewScheduler()
				Expect(err).ToNot(HaveOccurred())
				defer s.Shutdown()

				_, err = s.NewJob(
					gocron.DurationJob(time.Second),
					gocron.NewTask(nil),
				)
				Expect(err).To(HaveOccurred())
			})

			It("should shutdown gracefully while jobs are running", func() {
				var jobStarted sync.WaitGroup
				jobStarted.Add(1)

				s, err := gocron.NewScheduler()
				Expect(err).ToNot(HaveOccurred())

				_, err = s.NewJob(
					gocron.DurationJob(50*time.Millisecond),
					gocron.NewTask(func() {
						jobStarted.Done()
						time.Sleep(200 * time.Millisecond)
					}),
				)
				Expect(err).ToNot(HaveOccurred())

				s.Start()

				By("Waiting for job to start executing")
				jobStarted.Wait()

				By("Shutting down while job is running")
				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should recover from job panics via AfterJobRunsWithPanic", func() {
				var panicCaptured atomic.Bool
				fakeClock := clockwork.NewFakeClock()

				s, err := gocron.NewScheduler(gocron.WithClock(fakeClock))
				Expect(err).ToNot(HaveOccurred())

				var wg sync.WaitGroup
				wg.Add(1)
				_, err = s.NewJob(
					gocron.DurationJob(time.Second),
					gocron.NewTask(func() {
						panic("test panic")
					}),
					gocron.WithEventListeners(
						gocron.AfterJobRunsWithPanic(func(_ uuid.UUID, _ string, _ any) {
							panicCaptured.Store(true)
							wg.Done()
						}),
					),
				)
				Expect(err).ToNot(HaveOccurred())

				s.Start()

				By("Advancing clock to trigger panicking job")
				err = fakeClock.BlockUntilContext(context.Background(), 1)
				Expect(err).ToNot(HaveOccurred())
				fakeClock.Advance(time.Second)

				wg.Wait()
				Expect(panicCaptured.Load()).To(BeTrue())

				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Event Handling", func() {
			It("should handle job success callbacks", func() {
				var callbackCalled atomic.Bool
				fakeClock := clockwork.NewFakeClock()

				s, err := gocron.NewScheduler(gocron.WithClock(fakeClock))
				Expect(err).ToNot(HaveOccurred())

				var wg sync.WaitGroup
				wg.Add(1)
				_, err = s.NewJob(
					gocron.DurationJob(time.Second),
					gocron.NewTask(func() {}),
					gocron.WithEventListeners(
						gocron.AfterJobRuns(func(_ uuid.UUID, _ string) {
							callbackCalled.Store(true)
							wg.Done()
						}),
					),
				)
				Expect(err).ToNot(HaveOccurred())

				s.Start()

				By("Advancing clock to trigger job")
				err = fakeClock.BlockUntilContext(context.Background(), 1)
				Expect(err).ToNot(HaveOccurred())
				fakeClock.Advance(time.Second)

				wg.Wait()
				Expect(callbackCalled.Load()).To(BeTrue())

				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should handle job error callbacks", func() {
				var errorCaptured atomic.Bool
				fakeClock := clockwork.NewFakeClock()

				s, err := gocron.NewScheduler(gocron.WithClock(fakeClock))
				Expect(err).ToNot(HaveOccurred())

				var wg sync.WaitGroup
				wg.Add(1)
				_, err = s.NewJob(
					gocron.DurationJob(time.Second),
					gocron.NewTask(func() error {
						return fmt.Errorf("job error")
					}),
					gocron.WithEventListeners(
						gocron.AfterJobRunsWithError(func(_ uuid.UUID, _ string, _ error) {
							errorCaptured.Store(true)
							wg.Done()
						}),
					),
				)
				Expect(err).ToNot(HaveOccurred())

				s.Start()

				By("Advancing clock to trigger job")
				err = fakeClock.BlockUntilContext(context.Background(), 1)
				Expect(err).ToNot(HaveOccurred())
				fakeClock.Advance(time.Second)

				wg.Wait()
				Expect(errorCaptured.Load()).To(BeTrue())

				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Testing Support with FakeClock", func() {
			It("should advance fake time to trigger scheduled jobs deterministically", func() {
				fakeClock := clockwork.NewFakeClock()
				var count atomic.Int32

				s, err := gocron.NewScheduler(gocron.WithClock(fakeClock))
				Expect(err).ToNot(HaveOccurred())

				var wg sync.WaitGroup
				wg.Add(5)
				_, err = s.NewJob(
					gocron.DurationJob(10*time.Second),
					gocron.NewTask(func() {
						count.Add(1)
						wg.Done()
					}),
				)
				Expect(err).ToNot(HaveOccurred())

				s.Start()

				By("Advancing fake clock to trigger 5 executions")
				for i := 0; i < 5; i++ {
					err = fakeClock.BlockUntilContext(context.Background(), 1)
					Expect(err).ToNot(HaveOccurred())
					fakeClock.Advance(10 * time.Second)
				}

				wg.Wait()
				Expect(count.Load()).To(Equal(int32(5)))

				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
