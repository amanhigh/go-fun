package play_fast

import (
	"fmt"

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

		const (
			testKey        = "testKey"
			testValue      = "testValue"
			nonExistentKey = "nonExistentKey"
		)

		BeforeEach(func() {
			cache, err = ristretto.NewCache(&ristretto.Config{
				NumCounters: 20,       // number of keys to track frequency of (10M).
				MaxCost:     10 << 20, // maximum cost of cache (10 MB).
				BufferItems: 64,       // number of keys per Get buffer.
			})
		})

		It("should build", func() {
			Expect(err).To(BeNil())
			Expect(cache).To(Not(BeNil()))
		})

		Context("Basic Operations", func() {
			BeforeEach(func() {
				success := cache.Set(testKey, testValue, 1)
				Expect(success).To(BeTrue())
				cache.Wait() // Wait for value to pass through buffers
			})

			AfterEach(func() {
				cache.Del(testKey)
				_, found := cache.Get(testKey)
				Expect(found).To(BeFalse())
			})

			It("should get a value", func() {
				value, found := cache.Get(testKey)
				Expect(found).To(BeTrue())
				Expect(value).To(Equal(testValue))
			})
		})

		Context("Cache Behavior", func() {
			It("should handle cache miss", func() {
				_, found := cache.Get(nonExistentKey)
				Expect(found).To(BeFalse())
			})

			It("should evict items when cache is full", func() {
				By("Filling the cache")
				for i := 0; i < 20; i++ {
					key := fmt.Sprintf("key%d", i)
					success := cache.Set(key, i, 1)
					Expect(success).To(BeTrue())
				}

				// Add one more item to trigger eviction
				cache.Set("trigger", "value", 1)
				cache.Wait()

				By("Checking if some items were evicted")
				evictedCount := 0
				for i := 0; i < 20; i++ {
					key := fmt.Sprintf("key%d", i)
					_, found := cache.Get(key)
					if !found {
						evictedCount++
					}
				}

				Expect(evictedCount).To(BeNumerically(">", 0))
				Expect(evictedCount).To(BeNumerically("<=", 10))
			})
		})
	})
})
