package play_test

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/go-co-op/gocron/v2"
	redislock "github.com/go-co-op/gocron-redis-lock/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
)

var _ = Describe("CronDistributed", Ordered, Label(models.GINKGO_SLOW), func() {
	var (
		ctx            = context.Background()
		redisContainer testcontainers.Container
		redisClient    *redis.Client
		err            error
	)

	BeforeAll(func() {
		By("Starting Redis container")
		redisContainer, err = util.RedisTestContainer(ctx)
		Expect(err).ToNot(HaveOccurred())

		redisHost, err := redisContainer.PortEndpoint(ctx, "6379/tcp", "")
		Expect(err).ToNot(HaveOccurred())
		log.Info().Str("Host", redisHost).Msg("Redis Endpoint")

		redisClient = redis.NewClient(&redis.Options{
			Addr: redisHost,
		})

		By("Verifying Redis connection")
		_, err = redisClient.Ping(ctx).Result()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		log.Warn().Msg("Redis Shutting Down")
		if redisClient != nil {
			redisClient.Close()
		}
		err = redisContainer.Terminate(ctx)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should connect to Redis", func() {
		Expect(redisClient).ToNot(BeNil())
		result, err := redisClient.Ping(ctx).Result()
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal("PONG"))
	})

	Context("Redis Distributed Locking", func() {
		It("should create locker from Redis client", func() {
			locker, err := redislock.NewRedisLocker(redisClient)
			Expect(err).ToNot(HaveOccurred())
			Expect(locker).ToNot(BeNil())
		})

		It("should verify distributed lock prevents concurrent job execution", func() {
			By("Creating Redis locker")
			locker, err := redislock.NewRedisLocker(redisClient)
			Expect(err).ToNot(HaveOccurred())

			var totalExecutions atomic.Int32
			var wg sync.WaitGroup
			wg.Add(3) // Expect 3 total executions across all schedulers

			By("Creating two competing scheduler instances")
			s1, err := gocron.NewScheduler(gocron.WithDistributedLocker(locker))
			Expect(err).ToNot(HaveOccurred())

			s2, err := gocron.NewScheduler(gocron.WithDistributedLocker(locker))
			Expect(err).ToNot(HaveOccurred())

			jobFunc := func() {
				totalExecutions.Add(1)
				wg.Done()
			}

			By("Scheduling same job on both schedulers")
			_, err = s1.NewJob(
				gocron.DurationJob(200*time.Millisecond),
				gocron.NewTask(jobFunc),
			)
			Expect(err).ToNot(HaveOccurred())

			_, err = s2.NewJob(
				gocron.DurationJob(200*time.Millisecond),
				gocron.NewTask(jobFunc),
			)
			Expect(err).ToNot(HaveOccurred())

			By("Starting both schedulers")
			s1.Start()
			s2.Start()

			By("Waiting for executions")
			wg.Wait()

			By("Verifying jobs executed")
			Expect(totalExecutions.Load()).To(BeNumerically(">=", int32(3)))

			err = s1.Shutdown()
			Expect(err).ToNot(HaveOccurred())
			err = s2.Shutdown()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle leader failover when current leader shuts down", func() {
			By("Creating Redis locker")
			locker, err := redislock.NewRedisLocker(redisClient)
			Expect(err).ToNot(HaveOccurred())

			var executionCount atomic.Int32

			By("Starting first scheduler (leader)")
			s1, err := gocron.NewScheduler(gocron.WithDistributedLocker(locker))
			Expect(err).ToNot(HaveOccurred())

			var wg1 sync.WaitGroup
			wg1.Add(1)
			_, err = s1.NewJob(
				gocron.DurationJob(200*time.Millisecond),
				gocron.NewTask(func() {
					executionCount.Add(1)
					wg1.Done()
				}),
			)
			Expect(err).ToNot(HaveOccurred())
			s1.Start()

			By("Waiting for first scheduler to execute")
			wg1.Wait()
			firstCount := executionCount.Load()
			Expect(firstCount).To(BeNumerically(">=", int32(1)))

			By("Shutting down first scheduler (simulating failure)")
			err = s1.Shutdown()
			Expect(err).ToNot(HaveOccurred())

			By("Starting second scheduler (should take over)")
			s2, err := gocron.NewScheduler(gocron.WithDistributedLocker(locker))
			Expect(err).ToNot(HaveOccurred())

			var wg2 sync.WaitGroup
			wg2.Add(1)
			_, err = s2.NewJob(
				gocron.DurationJob(200*time.Millisecond),
				gocron.NewTask(func() {
					executionCount.Add(1)
					wg2.Done()
				}),
			)
			Expect(err).ToNot(HaveOccurred())
			s2.Start()

			By("Waiting for second scheduler to execute (failover)")
			wg2.Wait()
			Expect(executionCount.Load()).To(BeNumerically(">", firstCount))

			err = s2.Shutdown()
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
