package play_fast

import (
	"github.com/dgraph-io/ristretto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Cache", func() {
	Context("Ristretto", func() {
		var (
			cache *ristretto.Cache
			err   error
		)

		BeforeEach(func() {
			cache, err = ristretto.NewCache(&ristretto.Config{
				NumCounters: 1e7,     // number of keys to track frequency of (10M).
				MaxCost:     1 << 30, // maximum cost of cache (1GB).
				BufferItems: 64,      // number of keys per Get buffer.
			})
		})

		It("should build", func() {
			Expect(err).To(BeNil())
			Expect(cache).To(Not(BeNil()))
		})
	})
})
