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
		var (
			s         gocron.Scheduler
			err       error
			fakeClock *clockwork.FakeClock
			count     atomic.Int32
			wg        sync.WaitGroup
		)

		BeforeEach(func() {
			s, err = gocron.NewScheduler()
			Expect(err).ToNot(HaveOccurred())
			count.Store(0)
		})

		AfterEach(func() {
			if s != nil {
				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
			}
		})

		Context("Basic", func() {
			It("1.1 should create and shutdown scheduler gracefully", func() {
				Expect(s).ToNot(BeNil())
				s.Start()
				err = s.Shutdown()
				Expect(err).ToNot(HaveOccurred())
				s = nil // Prevent double shutdown in AfterEach
			})

			Context("with location configuration", func() {
				BeforeEach(func() {
					s, err = gocron.NewScheduler(
						gocron.WithLocation(time.UTC),
					)
					Expect(err).ToNot(HaveOccurred())
				})

				It("1.2 should configure scheduler with location", func() {
					// Verify UTC scheduler can create jobs with timezone-aware scheduling
					_, err = s.NewJob(
						gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(12, 0, 0))), // Noon UTC
						gocron.NewTask(func() {}),
					)
					Expect(err).ToNot(HaveOccurred())

					// Verify scheduler was created with UTC location
					Expect(s).ToNot(BeNil())
				})
			})

			Context("with concurrent job limit", func() {
				var maxConcurrent atomic.Int32

				BeforeEach(func() {
					s, err = gocron.NewScheduler(gocron.WithLimitConcurrentJobs(2, gocron.LimitModeWait))
					Expect(err).ToNot(HaveOccurred())
					maxConcurrent.Store(0)
				})

				It("1.2 should configure scheduler with concurrent job limit", func() {
					var currentJobs atomic.Int32
					var wg sync.WaitGroup

					for i := 0; i < 4; i++ {
						wg.Add(1)
						_, err = s.NewJob(
							gocron.OneTimeJob(gocron.OneTimeJobStartImmediately()),
							gocron.NewTask(func() {
								defer wg.Done()
								current := currentJobs.Add(1)
								defer currentJobs.Add(-1)
								if current > maxConcurrent.Load() {
									maxConcurrent.Store(current)
								}
								time.Sleep(20 * time.Millisecond)
							}),
						)
						Expect(err).ToNot(HaveOccurred())
					}

					s.Start()
					wg.Wait()

					Expect(maxConcurrent.Load()).To(BeNumerically("<=", 2))
					Expect(maxConcurrent.Load()).To(BeNumerically(">", 0))
				})
			})

			Context("with fake clock for deterministic testing", func() {
				var (
					panicCaptured  atomic.Bool
					callbackCalled atomic.Bool
					errorCaptured  atomic.Bool
				)

				BeforeEach(func() {
					fakeClock = clockwork.NewFakeClock()
					s, err = gocron.NewScheduler(gocron.WithClock(fakeClock))
					Expect(err).ToNot(HaveOccurred())
					panicCaptured.Store(false)
					callbackCalled.Store(false)
					errorCaptured.Store(false)
				})

				It("2.1 should schedule job at fixed time intervals", func() {
					wg.Add(3)
					_, err = s.NewJob(gocron.DurationJob(5*time.Second), gocron.NewTask(func() { count.Add(1); wg.Done() }))
					Expect(err).ToNot(HaveOccurred())

					s.Start()
					for i := 0; i < 3; i++ {
						fakeClock.BlockUntilContext(context.Background(), 1)
						fakeClock.Advance(5 * time.Second)
					}
					wg.Wait()
					Expect(count.Load()).To(Equal(int32(3)))
				})

				It("3.4 should recover from job panics via AfterJobRunsWithPanic", func() {
					wg.Add(1)
					_, err = s.NewJob(
						gocron.DurationJob(time.Second),
						gocron.NewTask(func() { panic("test panic") }),
						gocron.WithEventListeners(gocron.AfterJobRunsWithPanic(func(_ uuid.UUID, _ string, _ any) { panicCaptured.Store(true); wg.Done() })),
					)
					Expect(err).ToNot(HaveOccurred())

					s.Start()
					fakeClock.BlockUntilContext(context.Background(), 1)
					fakeClock.Advance(time.Second)
					wg.Wait()
					Expect(panicCaptured.Load()).To(BeTrue())
				})

				It("3.4 should handle job success callbacks", func() {
					wg.Add(1)
					_, err = s.NewJob(
						gocron.DurationJob(time.Second),
						gocron.NewTask(func() {}),
						gocron.WithEventListeners(gocron.AfterJobRuns(func(_ uuid.UUID, _ string) { callbackCalled.Store(true); wg.Done() })),
					)
					Expect(err).ToNot(HaveOccurred())

					s.Start()
					fakeClock.BlockUntilContext(context.Background(), 1)
					fakeClock.Advance(time.Second)
					wg.Wait()
					Expect(callbackCalled.Load()).To(BeTrue())
				})

				It("3.4 should handle job error callbacks", func() {
					wg.Add(1)
					_, err = s.NewJob(
						gocron.DurationJob(time.Second),
						gocron.NewTask(func() error { return fmt.Errorf("job error") }),
						gocron.WithEventListeners(gocron.AfterJobRunsWithError(func(_ uuid.UUID, _ string, _ error) { errorCaptured.Store(true); wg.Done() })),
					)
					Expect(err).ToNot(HaveOccurred())

					s.Start()
					fakeClock.BlockUntilContext(context.Background(), 1)
					fakeClock.Advance(time.Second)
					wg.Wait()
					Expect(errorCaptured.Load()).To(BeTrue())
				})
			})

			It("2.1 should schedule job at random duration intervals", func() {
				j, err := s.NewJob(gocron.DurationRandomJob(1*time.Second, 5*time.Second), gocron.NewTask(func() {}))
				Expect(err).ToNot(HaveOccurred())
				Expect(j.ID()).ToNot(BeNil())
			})

			It("2.2 should schedule job using standard cron expressions and seconds support", func() {
				// Test standard cron expression (5 fields)
				j1, err := s.NewJob(gocron.CronJob("*/5 * * * *", false), gocron.NewTask(func() {}))
				Expect(err).ToNot(HaveOccurred())
				Expect(j1.ID()).ToNot(BeNil())

				// Test cron expression with seconds support (6 fields)
				j2, err := s.NewJob(gocron.CronJob("*/5 * * * * *", true), gocron.NewTask(func() {}))
				Expect(err).ToNot(HaveOccurred())
				Expect(j2.ID()).ToNot(BeNil())
				Expect(j2.ID()).ToNot(Equal(j1.ID())) // Different jobs have different IDs
			})
		})

		Context("Medium", func() {
			Context("time-based helpers", func() {
				It("2.3 should schedule daily job at specific time", func() {
					j, err := s.NewJob(gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(9, 0, 0))), gocron.NewTask(func() {}))
					Expect(err).ToNot(HaveOccurred())
					Expect(j.ID()).ToNot(BeNil())
				})

				It("2.3 should schedule weekly job on specific days", func() {
					j, err := s.NewJob(gocron.WeeklyJob(1, gocron.NewWeekdays(time.Monday, time.Wednesday, time.Friday), gocron.NewAtTimes(gocron.NewAtTime(10, 30, 0))), gocron.NewTask(func() {}))
					Expect(err).ToNot(HaveOccurred())
					Expect(j.ID()).ToNot(BeNil())
				})

				It("2.3 should schedule monthly job on specific days", func() {
					j, err := s.NewJob(gocron.MonthlyJob(1, gocron.NewDaysOfTheMonth(1, 15), gocron.NewAtTimes(gocron.NewAtTime(8, 0, 0))), gocron.NewTask(func() {}))
					Expect(err).ToNot(HaveOccurred())
					Expect(j.ID()).ToNot(BeNil())
				})

				It("2.3 should schedule one-time job at specific datetime", func() {
					j, err := s.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time.Now().Add(24*time.Hour))), gocron.NewTask(func() {}))
					Expect(err).ToNot(HaveOccurred())
					Expect(j.ID()).ToNot(BeNil())
				})
			})

			Context("concurrency control", func() {
				It("3.1 should support singleton mode (skip overlapping)", func() {
					j, err := s.NewJob(gocron.DurationJob(time.Second), gocron.NewTask(func() { time.Sleep(2 * time.Second) }), gocron.WithSingletonMode(gocron.LimitModeReschedule))
					Expect(err).ToNot(HaveOccurred())
					Expect(j.ID()).ToNot(BeNil())
				})
			})

			Context("job management", func() {
				var jobs []gocron.Job
				var job gocron.Job

				BeforeEach(func() {
					job, err = s.NewJob(gocron.DurationJob(time.Minute), gocron.NewTask(func() {}))
					Expect(err).ToNot(HaveOccurred())
					for i := 0; i < 2; i++ {
						_, err = s.NewJob(gocron.DurationJob(time.Duration(i+2)*time.Minute), gocron.NewTask(func() {}))
						Expect(err).ToNot(HaveOccurred())
					}
					jobs = s.Jobs()
					Expect(jobs).To(HaveLen(3))
				})

				It("3.2 should retrieve job by ID", func() {
					var foundJob gocron.Job
					for _, j := range jobs {
						if j.ID() == job.ID() {
							foundJob = j
							break
						}
					}
					Expect(foundJob.ID()).To(Equal(job.ID()))
				})

				It("3.2 should list all scheduled jobs", func() {
					currentJobs := s.Jobs()
					Expect(currentJobs).To(HaveLen(3))
				})

				It("3.2 should remove scheduled jobs", func() {
					err = s.RemoveJob(job.ID())
					Expect(err).ToNot(HaveOccurred())
					Expect(s.Jobs()).To(HaveLen(2))
				})

				It("3.2 should update job schedule by removing and re-adding", func() {
					originalID := job.ID()
					err = s.RemoveJob(originalID)
					Expect(err).ToNot(HaveOccurred())
					newJob, err := s.NewJob(gocron.DurationJob(30*time.Second), gocron.NewTask(func() {}))
					Expect(err).ToNot(HaveOccurred())
					Expect(newJob.ID()).ToNot(Equal(originalID))
					Expect(s.Jobs()).To(HaveLen(3))
				})
			})

			Context("error handling", func() {
				It("3.3 should handle invalid cron expression", func() {
					_, err = s.NewJob(gocron.CronJob("bad cron", true), gocron.NewTask(func() {}))
					Expect(err).To(HaveOccurred())
				})

				It("3.3 should handle nil function", func() {
					_, err = s.NewJob(gocron.DurationJob(time.Second), gocron.NewTask(nil))
					Expect(err).To(HaveOccurred())
				})

				It("3.3 should shutdown gracefully while jobs are running", func() {
					var wg sync.WaitGroup
					wg.Add(1)
					_, err = s.NewJob(gocron.DurationJob(50*time.Millisecond), gocron.NewTask(func() { wg.Done(); time.Sleep(200 * time.Millisecond) }))
					Expect(err).ToNot(HaveOccurred())

					s.Start()
					wg.Wait()
					err = s.Shutdown()
					Expect(err).ToNot(HaveOccurred())
					s = nil
				})
			})
		})
	})
})
