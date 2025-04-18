package play_test

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/etcinit/speedbump"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"gopkg.in/redis.v5"
)

var _ = Describe("SpeedBump", Ordered, Label(models.GINKGO_SLOW), func() {
	var (
		err       error
		ctx       = context.Background()
		hasher    = speedbump.PerSecondHasher{}
		testIp    = "10.10.10.10"
		container testcontainers.Container
		client    *redis.Client
	)

	BeforeAll(func() {
		container, err = util.RedisTestContainer(ctx)
		Expect(err).ToNot(HaveOccurred())

		endpoint, err := container.Endpoint(ctx, "")
		Expect(err).ToNot(HaveOccurred())

		log.Info().Str("Host", endpoint).Msg("Redis Endpoint")

		// dman set redis
		client = redis.NewClient(&redis.Options{
			Addr:     endpoint,
			Password: "",
			DB:       0,
		})
	})

	AfterAll(func() {
		log.Warn().Msg("Redis Shutting Down")
		Expect(container.Terminate(ctx)).To(Succeed())
	})

	It("should build", func() {
		Expect(client).To(Not(BeNil()))
		Expect(hasher).To(Not(BeNil()))
	})

	It("should limit", func() {
		/*
			Gin can use gin bump
			https://github.com/etcinit/speedbump/tree/master/ginbump
			RateLimitLB() also honors X-Forwarded-For
		*/

		// Here we create a limiter that will only allow 5 requests per second
		limiter := speedbump.NewLimiter(client, hasher, 5) // Create one Limiter for each rate limit in usecase
		// First 5 Request not limited
		for i := 0; i < 5; i++ {
			success, err := limiter.Attempt(testIp) // TestIp can be combined with api to do api level rate limiting.
			Expect(err).ToNot(HaveOccurred())
			Expect(success).To(BeTrue())
		}
		// Next 5 Request are limited
		for i := 0; i < 5; i++ {
			success, err := limiter.Attempt(testIp)
			Expect(err).ToNot(HaveOccurred())
			Expect(success).To(BeFalse())
		}
	})
})
