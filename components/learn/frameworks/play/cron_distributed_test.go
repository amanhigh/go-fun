package play_test

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	redislock "github.com/go-co-op/gocron-redis-lock/v2"
	"github.com/go-co-op/gocron/v2"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
)

// pgAdvisoryLock implements gocron.Lock using PostgreSQL advisory locks
type pgAdvisoryLock struct {
	conn    *sql.Conn
	lockKey int64
}

func (l *pgAdvisoryLock) Unlock(ctx context.Context) error {
	_, err := l.conn.ExecContext(ctx, "SELECT pg_advisory_unlock($1)", l.lockKey)
	closeErr := l.conn.Close()
	if err != nil {
		return err
	}
	return closeErr
}

// pgAdvisoryLocker implements gocron.Locker using PostgreSQL session-level advisory locks
type pgAdvisoryLocker struct {
	db *sql.DB
}

func newPgAdvisoryLocker(db *sql.DB) gocron.Locker {
	return &pgAdvisoryLocker{db: db}
}

func (l *pgAdvisoryLocker) Lock(ctx context.Context, key string) (gocron.Lock, error) {
	// Hash the key to get a consistent positive int64 for advisory lock
	h := fnv.New64a()
	h.Write([]byte(key))
	lockKey := int64(h.Sum64() & math.MaxInt64)

	// Use a dedicated connection for session-level advisory lock
	conn, err := l.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	var acquired bool
	err = conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1)", lockKey).Scan(&acquired)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to acquire advisory lock: %w", err)
	}
	if !acquired {
		conn.Close()
		return nil, fmt.Errorf("advisory lock not acquired for key: %s", key)
	}

	return &pgAdvisoryLock{conn: conn, lockKey: lockKey}, nil
}

var _ gocron.Locker = (*pgAdvisoryLocker)(nil)
var _ gocron.Lock = (*pgAdvisoryLock)(nil)

var _ = Describe("CronDistributed", Ordered, Label(models.GINKGO_SLOW), func() {
	var (
		ctx            = context.Background()
		redisContainer testcontainers.Container
		redisClient    *redis.Client
		pgContainer    testcontainers.Container
		pgDB           *sql.DB
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

		By("Starting PostgreSQL container")
		pgContainer, err = util.PostgresTestContainer(ctx)
		Expect(err).ToNot(HaveOccurred())

		pgHost, err := pgContainer.PortEndpoint(ctx, "5432/tcp", "")
		Expect(err).ToNot(HaveOccurred())
		log.Info().Str("Host", pgHost).Msg("PostgreSQL Endpoint")

		dsn := fmt.Sprintf("postgres://test:test@%s/testdb?sslmode=disable", pgHost)
		pgDB, err = sql.Open("postgres", dsn)
		Expect(err).ToNot(HaveOccurred())

		By("Verifying PostgreSQL connection")
		err = pgDB.PingContext(ctx)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		log.Warn().Msg("Redis Shutting Down")
		if redisClient != nil {
			redisClient.Close()
		}
		err = redisContainer.Terminate(ctx)
		Expect(err).ToNot(HaveOccurred())

		log.Warn().Msg("PostgreSQL Shutting Down")
		if pgDB != nil {
			pgDB.Close()
		}
		err = pgContainer.Terminate(ctx)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should connect to Redis", func() {
		Expect(redisClient).ToNot(BeNil())
		result, err := redisClient.Ping(ctx).Result()
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal("PONG"))
	})

	It("should connect to PostgreSQL", func() {
		Expect(pgDB).ToNot(BeNil())
		err := pgDB.PingContext(ctx)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("PostgreSQL Distributed Locking", func() {
		It("should acquire and release advisory lock", func() {
			By("Creating PostgreSQL advisory locker")
			locker := newPgAdvisoryLocker(pgDB)
			Expect(locker).ToNot(BeNil())

			By("Acquiring lock")
			lock, err := locker.Lock(ctx, "test-job")
			Expect(err).ToNot(HaveOccurred())
			Expect(lock).ToNot(BeNil())

			By("Releasing lock")
			err = lock.Unlock(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should prevent concurrent lock acquisition for same key", func() {
			By("Creating two lockers sharing the same DB pool")
			locker1 := newPgAdvisoryLocker(pgDB)
			locker2 := newPgAdvisoryLocker(pgDB)

			By("First locker acquires lock")
			lock1, err := locker1.Lock(ctx, "exclusive-job")
			Expect(err).ToNot(HaveOccurred())

			By("Second locker fails to acquire same lock")
			_, err = locker2.Lock(ctx, "exclusive-job")
			Expect(err).To(HaveOccurred())

			By("Releasing first lock")
			err = lock1.Unlock(ctx)
			Expect(err).ToNot(HaveOccurred())

			By("Second locker can now acquire the lock")
			lock2, err := locker2.Lock(ctx, "exclusive-job")
			Expect(err).ToNot(HaveOccurred())
			err = lock2.Unlock(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should verify only one scheduler executes jobs with PostgreSQL lock", func() {
			By("Creating PostgreSQL locker")
			locker := newPgAdvisoryLocker(pgDB)

			var totalExecutions atomic.Int32
			var wg sync.WaitGroup
			wg.Add(3)

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

		It("should handle failover when leader PostgreSQL session ends", func() {
			By("Creating PostgreSQL locker")
			locker := newPgAdvisoryLocker(pgDB)

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
			wg.Add(3)

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
