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
			updatedValue   = "updatedValue"
			nonExistentKey = "nonExistentKey"
			cacheSize      = 100
		)

		BeforeEach(func() {
			cache, err = ristretto.NewCache(&ristretto.Config{
				NumCounters: cacheSize * 10, // No. of counters (10x of MaxCost)
				MaxCost:     cacheSize,      // Maximum number of entries (Can be in any unit eg. MB)
				BufferItems: 64,             // number of keys per Get buffer.
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

			It("should update an existing key", func() {
				By("Updating the existing key")
				success := cache.Set(testKey, updatedValue, 1)
				Expect(success).To(BeTrue())
				cache.Wait()

				By("Verifying the updated value")
				value, found := cache.Get(testKey)
				Expect(found).To(BeTrue())
				Expect(value).To(Equal(updatedValue))
			})
		})

		Context("Cache Behavior", func() {
			It("should handle cache miss", func() {
				_, found := cache.Get(nonExistentKey)
				Expect(found).To(BeFalse())
			})

			It("should evict items when cache is full", func() {
				By("Filling the cache")
				for i := 0; i < cacheSize; i++ {
					key := fmt.Sprintf("key%d", i)
					success := cache.Set(key, i, 1)
					Expect(success).To(BeTrue())
				}

				// Add one more item to trigger eviction
				cache.Set("trigger", "value", 1)
				cache.Wait()

				By("Checking if some items were evicted")
				evictedCount := 0
				for i := 0; i < cacheSize; i++ {
					key := fmt.Sprintf("key%d", i)
					_, found := cache.Get(key)
					if !found {
						evictedCount++
					}
				}

				Expect(evictedCount).To(BeNumerically(">", 0))
				Expect(evictedCount).To(BeNumerically("<=", cacheSize))
			})

			It("should clear all items from the cache", func() {
				By("Adding items to the cache")
				for i := 0; i < 10; i++ {
					key := fmt.Sprintf("key%d", i)
					success := cache.Set(key, i, 1)
					Expect(success).To(BeTrue())
				}
				cache.Wait()

				By("Clearing the cache")
				cache.Clear()

				By("Verifying all items are removed")
				for i := 0; i < 10; i++ {
					key := fmt.Sprintf("key%d", i)
					_, found := cache.Get(key)
					Expect(found).To(BeFalse())
				}
			})
		})
	})
})
