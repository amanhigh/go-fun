package play_fast

import (
	"fmt"
	"sync"
	"time"

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
				NumCounters:        cacheSize * 10, // No. of counters (10x of MaxCost)
				MaxCost:            cacheSize,      // Maximum number of entries (Can be in any unit eg. MB)
				BufferItems:        64,             // number of keys per Get buffer.
				Metrics:            true,           // Enable metrics collection
				IgnoreInternalCost: true,           // Ignore internal cost calculation (Non Byte Costs)
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

			It("should handle cache miss", func() {
				_, found := cache.Get(nonExistentKey)
				Expect(found).To(BeFalse())
			})

			It("should respect TTL for items", func() {
				ttlKey := "ttlKey"
				ttlValue := "ttlValue"
				ttlDuration := 100 * time.Millisecond

				By("Setting a key with TTL")
				success := cache.SetWithTTL(ttlKey, ttlValue, 1, ttlDuration)
				Expect(success).To(BeTrue())
				cache.Wait()

				By("Verifying the key exists immediately")
				value, found := cache.Get(ttlKey)
				Expect(found).To(BeTrue())
				Expect(value).To(Equal(ttlValue))

				By("Verifying the key is eventually removed")
				Eventually(func() bool {
					_, found := cache.Get(ttlKey)
					return found
				}, "200ms", "10ms").Should(BeFalse())
			})
		})

		Context("Cache Bulk Operations", func() {
			const itemsToAdd = 50

			BeforeEach(func() {
				By("Adding multiple items to the cache")
				for i := 0; i < itemsToAdd; i++ {
					key := fmt.Sprintf("key%d", i)
					success := cache.Set(key, i, 1)
					Expect(success).To(BeTrue())
				}
				cache.Wait()
			})

			It("should evict items when cache is full", func() {
				By("Filling the cache to its maximum capacity")
				for i := itemsToAdd; i < cacheSize; i++ {
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
				By("Clearing the cache")
				cache.Clear()

				By("Verifying all items are removed")
				for i := 0; i < itemsToAdd; i++ {
					key := fmt.Sprintf("key%d", i)
					_, found := cache.Get(key)
					Expect(found).To(BeFalse())
				}
			})
		})

		Context("Advanced Behavior", func() {
			It("should handle concurrent reads and writes", func() {
				const (
					numGoroutines = 100
					numOperations = 1000
				)
				var wg sync.WaitGroup
				wg.Add(numGoroutines)

				for i := 0; i < numGoroutines; i++ {
					go func(id int) {
						defer wg.Done()
						for j := 0; j < numOperations; j++ {
							key := fmt.Sprintf("key%d-%d", id, j)
							value := fmt.Sprintf("value%d-%d", id, j)
							// Randomly choose between Set and Get operations
							if j%2 == 0 {
								cache.Set(key, value, 1)
							} else {
								cache.Get(key)
							}
						}
					}(i)
				}

				wg.Wait()
				cache.Wait()

				By("Verifying the cache is still functional after concurrent operations")
				testKey := "concurrentTestKey"
				testValue := "concurrentTestValue"
				success := cache.Set(testKey, testValue, 1)
				Expect(success).To(BeTrue())
				cache.Wait()

				value, found := cache.Get(testKey)
				Expect(found).To(BeTrue())
				Expect(value).To(Equal(testValue))
			})

			It("should verify that the OnEvict function is invoked for evicted items", func() {
				evictedItems := []interface{}{}

				// Create a new cache with an OnEvict function
				cacheWithEvict, err := ristretto.NewCache(&ristretto.Config{
					NumCounters:        cacheSize * 10,
					MaxCost:            int64(cacheSize),
					BufferItems:        64,
					IgnoreInternalCost: true,
					OnEvict: func(item *ristretto.Item) {
						evictedItems = append(evictedItems, item.Key)
					},
				})
				Expect(err).To(BeNil())

				By("Filling the cache to its maximum capacity")
				for i := uint64(0); i < uint64(cacheSize); i++ {
					success := cacheWithEvict.Set(i, fmt.Sprintf("value%d", i), 1)
					Expect(success).To(BeTrue())
				}
				cacheWithEvict.Wait()

				By("Adding multiple items to trigger evictions")
				numExtraItems := 5
				for i := uint64(cacheSize); i < uint64(cacheSize+numExtraItems); i++ {
					success := cacheWithEvict.Set(i, fmt.Sprintf("value%d", i), 1)
					Expect(success).To(BeTrue())
				}
				cacheWithEvict.Wait()

				By("Verifying the eviction callback")
				Expect(evictedItems).ToNot(BeEmpty())
				Expect(len(evictedItems)).To(BeNumerically(">=", numExtraItems))

				By("Checking that evicted items are within the expected range")
				for _, item := range evictedItems {
					key, ok := item.(uint64)
					Expect(ok).To(BeTrue(), "Evicted item key should be of type uint64")
					Expect(key).To(BeNumerically(">=", uint64(0)))
					Expect(key).To(BeNumerically("<", uint64(cacheSize+numExtraItems)))
				}
			})

			It("should properly handle cost-based eviction", func() {
				By("Adding items with various costs")
				cache.Set("key1", "value1", 20)
				cache.Set("key2", "value2", 30)
				cache.Set("key3", "value3", 25)
				cache.Set("key4", "value4", 15)
				cache.Wait()

				By("Verifying all items are present")
				for i := 1; i <= 4; i++ {
					key := fmt.Sprintf("key%d", i)
					_, found := cache.Get(key)
					Expect(found).To(BeTrue(), fmt.Sprintf("Expected %s to be in the cache", key))
				}

				By("Adding an item that exceeds the remaining cost")
				cache.Set("key5", "value5", 50)
				cache.Wait()

				By("Verifying that some items were evicted and some remain")
				evictedCount := 0
				remainingCount := 0
				for i := 1; i <= 5; i++ {
					key := fmt.Sprintf("key%d", i)
					_, found := cache.Get(key)
					if found {
						remainingCount++
					} else {
						evictedCount++
					}
				}
				Expect(evictedCount).To(BeNumerically(">", 0), "Expected atleast one item to be evicted")
				Expect(remainingCount).To(BeNumerically(">", 1), "Expected more than one item to remain in the cache")

				By("Verifying the total cost does not exceed the maximum")
				metrics := cache.Metrics
				Expect(metrics.CostAdded()-metrics.CostEvicted()).To(BeNumerically("<=", cacheSize), "Expected total cost to not exceed maxCost")
			})
		})

		Context("Negative Scenarios", func() {
			It("should handle attempts to set items with negative costs", func() {
				By("Attempting to set an item with a negative cost")
				success := cache.Set("negativeKey", "negativeValue", -1)
				Expect(success).To(BeTrue(), "Setting an item with negative cost should succeed")

				cache.Wait()

				By("Verifying that the item is in the cache")
				value, found := cache.Get("negativeKey")
				Expect(found).To(BeTrue(), "Negative cost item should be in the cache")
				Expect(value).To(Equal("negativeValue"), "Negative cost item should have the correct value")
			})

			It("should handle zero cost items correctly", func() {
				By("Setting an item with zero cost")
				success := cache.Set("zeroCostKey", "zeroCostValue", 0)
				Expect(success).To(BeTrue(), "Setting an item with zero cost should succeed")
				cache.Wait()

				By("Verifying the zero cost item can be retrieved")
				value, found := cache.Get("zeroCostKey")
				Expect(found).To(BeTrue(), "Zero cost item should be in the cache")
				Expect(value).To(Equal("zeroCostValue"), "Zero cost item should have the correct value")
			})
		})
	})
})
