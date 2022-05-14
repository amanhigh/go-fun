package play_test

import (
	"github.com/etcinit/speedbump"
	"github.com/go-redis/redis/v8"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("SpeedBump", func() {
	var (
		client = redis.NewClient(&redis.Options{
			Addr:     "docker:6379",
			Password: "",
			DB:       0,
		})

		hasher = speedbump.PerSecondHasher{}
	)
	It("should build", func() {
		Expect(client).To(Not(BeNil()))
		Expect(hasher).To(Not(BeNil()))
	})

})
