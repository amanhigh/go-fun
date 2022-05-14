package play_test

import (
	"github.com/amanhigh/go-fun/models"
	"github.com/etcinit/speedbump"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/redis.v5"
)

var _ = Describe("SpeedBump", Label(models.GINKGO_SETUP), func() {
	var (
		// dman set redis
		client = redis.NewClient(&redis.Options{
			Addr:     "docker:6379",
			Password: "",
			DB:       0,
		})
		hasher = speedbump.PerSecondHasher{}
		testIp = "127.0.0.1"
	)

	It("should build", func() {
		Expect(client).To(Not(BeNil()))
		Expect(hasher).To(Not(BeNil()))
	})

	It("should limit", func() {
		// Here we create a limiter that will only allow 5 requests per second
		limiter := speedbump.NewLimiter(client, hasher, 5) //Create one Limiter for each rate limit in usecase
		// First 5 Request not limited
		for i := 0; i < 5; i++ {
			success, err := limiter.Attempt(testIp) //TestIp can be combined with api to do api level rate limiting.
			Expect(err).To(BeNil())
			Expect(success).To(BeTrue())
		}
		// Next 5 Request are limited
		for i := 0; i < 5; i++ {
			success, err := limiter.Attempt(testIp)
			Expect(err).To(BeNil())
			Expect(success).To(BeFalse())
		}
	})
})
