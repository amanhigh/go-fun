package play_fast

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/maypok86/otter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"
	"github.com/viccon/sturdyc"
)

// testMetricsRecorder implements sturdyc.MetricsRecorder for testing
type testMetricsRecorder struct {
	onHit  func()
	onMiss func()
}

func (r *testMetricsRecorder) CacheHit()                     { r.onHit() }
func (r *testMetricsRecorder) CacheMiss()                    { r.onMiss() }
func (r *testMetricsRecorder) AsynchronousRefresh()          {}
func (r *testMetricsRecorder) SynchronousRefresh()           {}
func (r *testMetricsRecorder) MissingRecord()                {}
func (r *testMetricsRecorder) ForcedEviction()               {}
func (r *testMetricsRecorder) EntriesEvicted(_ int)          {}
func (r *testMetricsRecorder) ShardIndex(_ int)              {}
func (r *testMetricsRecorder) CacheBatchRefreshSize(_ int)   {}
func (r *testMetricsRecorder) ObserveCacheSize(_ func() int) {}

var _ sturdyc.MetricsRecorder = (*testMetricsRecorder)(nil)

var _ = Describe("Cache", func() {

	// Common test constants for all cache implementations
	const (
		testKey        = "testKey"
		testValue      = "testValue"
		updatedValue   = "updatedValue"
		nonExistentKey = "nonExistentKey"
		cacheSize      = 100
	)

	// TASK: Bigcache, Freecache
	Context("Ristretto", func() {
		var (
			cache *ristretto.Cache
			err   error
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
			Expect(err).ToNot(HaveOccurred())
			Expect(cache).To(Not(BeNil()))
		})

		Context("Basic", func() {
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

			It("1.2 should get a value", func() {
				value, found := cache.Get(testKey)
				Expect(found).To(BeTrue())
				Expect(value).To(Equal(testValue))
			})

			It("1.1 should update an existing key", func() {
				By("Updating the existing key")
				success := cache.Set(testKey, updatedValue, 1)
				Expect(success).To(BeTrue())
				cache.Wait()

				By("Verifying the updated value")
				value, found := cache.Get(testKey)
				Expect(found).To(BeTrue())
				Expect(value).To(Equal(updatedValue))
			})

			It("1.3 should handle cache miss", func() {
				_, found := cache.Get(nonExistentKey)
				Expect(found).To(BeFalse())
			})

			It("1.4 should delete a value", func() {
				cache.Del(testKey)
				_, found := cache.Get(testKey)
				Expect(found).To(BeFalse())
			})

			It("1.5 should respect TTL for items", func() {
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

		Context("Medium", func() {
			Context("Cache Bulk Operations", func() {
				const itemsToAdd = 50

				BeforeEach(func() {
					By("Adding multiple items to the cache")
					for i := range itemsToAdd {
						key := fmt.Sprintf("key%d", i)
						success := cache.Set(key, i, 1)
						Expect(success).To(BeTrue())
					}
					cache.Wait()
				})

				It("1.6 should evict items when cache is full", func() {
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
					for i := range cacheSize {
						key := fmt.Sprintf("key%d", i)
						_, found := cache.Get(key)
						if !found {
							evictedCount++
						}
					}

					Expect(evictedCount).To(BeNumerically(">", 0))
					Expect(evictedCount).To(BeNumerically("<=", cacheSize))
				})

				It("2.2 should clear all items from the cache", func() {
					By("Clearing the cache")
					cache.Clear()

					By("Verifying all items are removed")
					for i := range itemsToAdd {
						key := fmt.Sprintf("key%d", i)
						_, found := cache.Get(key)
						Expect(found).To(BeFalse())
					}
				})

				It("2.3 should maintain valid cache metrics", func() {
					By("Performing cache hits")
					hitCount := itemsToAdd / 2
					for i := range hitCount {
						_, found := cache.Get(fmt.Sprintf("key%d", i))
						Expect(found).To(BeTrue(), fmt.Sprintf("Expected key%d to be found", i))
					}

					By("Performing cache misses")
					missCount := itemsToAdd / 4
					for i := itemsToAdd; i < itemsToAdd+missCount; i++ {
						_, found := cache.Get(fmt.Sprintf("key%d", i))
						Expect(found).To(BeFalse(), fmt.Sprintf("Expected key%d to not be found", i))
					}

					By("Verifying final metrics")
					metrics := cache.Metrics
					Expect(metrics.Hits()).To(Equal(uint64(hitCount)), "Hit count should match")
					Expect(metrics.Misses()).To(Equal(uint64(missCount)), "Miss count should match")
					Expect(metrics.KeysAdded()).To(Equal(uint64(itemsToAdd)), "Keys added should match itemsToAdd")
					Expect(metrics.CostAdded()).To(Equal(uint64(itemsToAdd)), "Cost added should match itemsToAdd")
				})
			})

			Context("Negative Scenarios", func() {
				It("3.7 should handle attempts to set items with negative costs", func() {
					By("Attempting to set an item with a negative cost")
					success := cache.Set("negativeKey", "negativeValue", -1)
					Expect(success).To(BeTrue(), "Setting an item with negative cost should succeed")

					cache.Wait()

					By("Verifying that the item is in the cache")
					value, found := cache.Get("negativeKey")
					Expect(found).To(BeTrue(), "Negative cost item should be in the cache")
					Expect(value).To(Equal("negativeValue"), "Negative cost item should have the correct value")
				})

				It("3.7 should handle zero cost items correctly", func() {
					By("Setting an item with zero cost")
					success := cache.Set("zeroCostKey", "zeroCostValue", 0)
					Expect(success).To(BeTrue(), "Setting an item with zero cost should succeed")
					cache.Wait()

					By("Verifying the zero cost item can be retrieved")
					value, found := cache.Get("zeroCostKey")
					Expect(found).To(BeTrue(), "Zero cost item should be in the cache")
					Expect(value).To(Equal("zeroCostValue"), "Zero cost item should have the correct value")
				})

				It("3.7 should handle attempts to set malformed values", func() {
					malformedValue := struct {
						Data string
					}{
						Data: strings.Repeat("A", 1<<20), // 1MB of data
					}
					success := cache.Set("malformedKey", malformedValue, 1)
					Expect(success).To(BeTrue(), "Setting a large value should succeed")
					cache.Wait()

					value, found := cache.Get("malformedKey")
					Expect(found).To(BeTrue(), "Malformed value should be in the cache")
					Expect(value).To(Equal(malformedValue), "Retrieved value should match the set value")
				})

				It("1.7 should reject items with extremely large costs", func() {
					success := cache.Set("largeCostKey", "largeCostValue", 1<<30) // 1GB cost
					Expect(success).To(BeTrue(), "Setting a large cost item should succeed")

					cache.Wait()

					_, found := cache.Get("largeCostKey")
					Expect(found).To(BeFalse(), "Large cost item should not be in the cache")
				})
			})
		})

		Context("Advanced", func() {
			Context("Advanced Behavior", func() {
				It("3.7 should handle concurrent reads and writes", func() {
					const (
						numGoroutines = 100
						numOperations = 1000
					)
					var wg sync.WaitGroup
					wg.Add(numGoroutines)

					for i := range numGoroutines {
						go func(id int) {
							defer wg.Done()
							for j := range numOperations {
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

				It("3.7 should verify that the OnEvict function is invoked for evicted items", func() {
					evictedItems := []any{}

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
					Expect(err).ToNot(HaveOccurred())

					By("Filling the cache to its maximum capacity")
					for i := range uint64(cacheSize) {
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

				It("1.7 should properly handle cost-based eviction", func() {
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

			Context("Performance Benchmarks", FlakeAttempts(3), func() {
				var benchCache *ristretto.Cache
				const numOperations = 10000

				BeforeEach(func() {
					var err error
					benchCache, err = ristretto.NewCache(&ristretto.Config{
						NumCounters: 1e5,     // number of keys to track frequency of (100K).
						MaxCost:     1 << 23, // maximum cost of cache (8MB).
						BufferItems: 64,      // number of keys per Get buffer.
					})
					Expect(err).ToNot(HaveOccurred())
				})
				It("2.4 should perform set operations efficiently", FlakeAttempts(3), func() {
					experiment := gmeasure.NewExperiment("Set Operations")
					AddReportEntry(experiment.Name, experiment)

					experiment.SampleDuration("set with consistency", func(_ int) {
						key := fmt.Sprintf("key-%d", GinkgoRandomSeed())
						success := benchCache.Set(key, "value", 1)
						Expect(success).To(BeTrue())
						benchCache.Wait()
					}, gmeasure.SamplingConfig{N: numOperations})

					experiment.SampleDuration("set without consistency", func(_ int) {
						key := fmt.Sprintf("key-%d", GinkgoRandomSeed())
						success := benchCache.Set(key, "value", 1)
						Expect(success).To(BeTrue())
					}, gmeasure.SamplingConfig{N: numOperations})

					AddReportEntry("Set Operations Stats", experiment.GetStats("set with consistency"))
					AddReportEntry("Set Operations Stats", experiment.GetStats("set without consistency"))

					Expect(experiment.GetStats("set with consistency").DurationFor(gmeasure.StatMedian)).To(BeNumerically("<", 3*time.Microsecond), "Median set with consistency should be less than 3µs")
					Expect(experiment.GetStats("set without consistency").DurationFor(gmeasure.StatMedian)).To(BeNumerically("<", 1*time.Microsecond), "Median set without consistency should be less than 1µs")
				})

				It("2.4 should perform get operations efficiently", func() {
					experiment := gmeasure.NewExperiment("Get Operations")
					AddReportEntry(experiment.Name, experiment)

					// Populate cache first
					for i := range numOperations {
						key := fmt.Sprintf("key-%d", i)
						success := benchCache.Set(key, fmt.Sprintf("value-%d", i), 1)
						Expect(success).To(BeTrue())
					}
					benchCache.Wait()

					experiment.SampleDuration("get", func(_ int) {
						key := fmt.Sprintf("key-%d", GinkgoRandomSeed()%numOperations)
						_, found := benchCache.Get(key)
						Expect(found).To(BeTrue())
					}, gmeasure.SamplingConfig{N: numOperations})

					AddReportEntry("Get Operations Stats", experiment.GetStats("get"))

					Expect(experiment.GetStats("get").DurationFor(gmeasure.StatMedian)).To(BeNumerically("<", 1*time.Microsecond), "Median get should be less than 1µs")
					Expect(experiment.GetStats("get").DurationFor(gmeasure.StatMax)).To(BeNumerically("<", 20*time.Microsecond), "Max get should be less than 20µs")
				})

				It("2.4 should perform delete operations efficiently", func() {
					// TODO: #A Slow Tests are taking time ensure they run faster (1-2 sec), lighten or remove benchmarks.
					experiment := gmeasure.NewExperiment("Delete Operations")
					AddReportEntry(experiment.Name, experiment)

					// Populate cache first
					for i := range numOperations {
						key := fmt.Sprintf("key-%d", i)
						success := benchCache.Set(key, fmt.Sprintf("value-%d", i), 1)
						Expect(success).To(BeTrue())
					}
					benchCache.Wait()

					experiment.SampleDuration("delete", func(_ int) {
						key := fmt.Sprintf("key-%d", GinkgoRandomSeed()%numOperations)
						benchCache.Del(key)
					}, gmeasure.SamplingConfig{N: numOperations})

					AddReportEntry("Delete Operations Stats", experiment.GetStats("delete"))

					Expect(experiment.GetStats("delete").DurationFor(gmeasure.StatMedian)).To(BeNumerically("<", 1*time.Microsecond), "Median delete should be less than 1µs")
					Expect(experiment.GetStats("delete").DurationFor(gmeasure.StatMax)).To(BeNumerically("<", 40*time.Microsecond), "Max delete should be less than 40µs")
				})
			})
		})

	})

	Context("Otter", func() {
		var (
			otterCache otter.Cache[string, string]
			err        error
		)

		BeforeEach(func() {
			otterCache, err = otter.MustBuilder[string, string](cacheSize).
				CollectStats().
				Build()
		})

		AfterEach(func() {
			otterCache.Close()
		})

		It("should build", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(otterCache).To(Not(BeNil()))
		})

		Context("Basic", func() {
			BeforeEach(func() {
				By("Setting a value")
				ok := otterCache.Set(testKey, testValue)
				Expect(ok).To(BeTrue())
			})

			It("1.2 should get a value", func() {
				value, found := otterCache.Get(testKey)
				Expect(found).To(BeTrue())
				Expect(value).To(Equal(testValue))
			})

			It("1.1 should update an existing key", func() {
				By("Updating the existing key")
				ok := otterCache.Set(testKey, updatedValue)
				Expect(ok).To(BeTrue())

				By("Verifying the updated value")
				value, found := otterCache.Get(testKey)
				Expect(found).To(BeTrue())
				Expect(value).To(Equal(updatedValue))
			})

			It("1.3 should handle cache miss", func() {
				_, found := otterCache.Get(nonExistentKey)
				Expect(found).To(BeFalse())
			})

			It("1.4 should delete a value", func() {
				otterCache.Delete(testKey)
				_, found := otterCache.Get(testKey)
				Expect(found).To(BeFalse())
			})

			It("1.5 should respect TTL for items", func() {
				ttlKey := "ttlKey"
				ttlValue := "ttlValue"
				ttlDuration := 100 * time.Millisecond

				By("Creating TTL cache")
				ttlCache, err := otter.MustBuilder[string, string](cacheSize).
					WithTTL(ttlDuration).
					Build()
				Expect(err).ToNot(HaveOccurred())
				defer ttlCache.Close()

				By("Setting a key with TTL")
				ok := ttlCache.Set(ttlKey, ttlValue)
				Expect(ok).To(BeTrue())

				By("Verifying the key exists immediately")
				value, found := ttlCache.Get(ttlKey)
				Expect(found).To(BeTrue())
				Expect(value).To(Equal(ttlValue))

				By("Verifying the key is eventually removed")
				Eventually(func() bool {
					_, found := ttlCache.Get(ttlKey)
					return found
				}, "1.2s", "50ms").Should(BeFalse())
			})
		})

		Context("Advanced", func() {
			Context("TinyLFU Adaptive Hit Ratio", func() {
				It("3.1 should maintain high hit ratio with frequency-based access", func() {
					By("Creating cache with capacity for hot items")
					smallCache, err := otter.MustBuilder[string, int](100).
						CollectStats().
						Build()
					Expect(err).ToNot(HaveOccurred())
					defer smallCache.Close()

					By("Populating cache with hot items and warming up")
					for i := range 50 {
						smallCache.Set(fmt.Sprintf("key%d", i), i)
					}
					// Allow async admission to process
					time.Sleep(50 * time.Millisecond)

					By("Accessing hot items to build frequency")
					for range 3 {
						for i := range 10 {
							smallCache.Get(fmt.Sprintf("key%d", i))
						}
					}
					time.Sleep(50 * time.Millisecond)

					By("Measuring hit ratio on hot items")
					hits := 0
					total := 100
					for i := range total {
						key := fmt.Sprintf("key%d", i%10) // hot items only
						if _, found := smallCache.Get(key); found {
							hits++
						}
					}

					By("Verifying hit ratio")
					hitRatio := float64(hits) / float64(total)
					Expect(hitRatio).To(BeNumerically(">", 0.70),
						fmt.Sprintf("Hit ratio %.2f%% should be >70%%", hitRatio*100))
				})
			})
		})
	})

	Context("Sturdyc", func() {
		var (
			cacheClient *sturdyc.Client[string]
		)

		const (
			// Sturdyc-specific configuration
			capacity    = 1000
			numShards   = 10
			ttl         = 5 * time.Second
			evictionPct = 10
		)

		BeforeEach(func() {
			cacheClient = sturdyc.New[string](capacity, numShards, ttl, evictionPct)
		})

		It("should build", func() {
			Expect(cacheClient).To(Not(BeNil()))
		})

		Context("Basic", func() {
			BeforeEach(func() {
				cacheClient.Set(testKey, testValue)
			})

			It("1.2 should get a value", func() {
				value, found := cacheClient.Get(testKey)
				Expect(found).To(BeTrue())
				Expect(value).To(Equal(testValue))
			})

			It("1.1 should update an existing key", func() {
				By("Updating the existing key")
				cacheClient.Set(testKey, updatedValue)

				By("Verifying the updated value")
				value, found := cacheClient.Get(testKey)
				Expect(found).To(BeTrue())
				Expect(value).To(Equal(updatedValue))
			})

			It("1.3 should handle cache miss", func() {
				_, found := cacheClient.Get(nonExistentKey)
				Expect(found).To(BeFalse())
			})

			It("1.4 should delete a value", func() {
				cacheClient.Delete(testKey)
				_, found := cacheClient.Get(testKey)
				Expect(found).To(BeFalse())
			})

			It("1.5 should respect TTL expiration", func() {
				By("Creating short TTL cache")
				shortTTLCache := sturdyc.New[string](capacity, numShards, 100*time.Millisecond, evictionPct)
				shortTTLCache.Set("ttlKey", "ttlValue")

				By("Verifying the key exists immediately")
				value, found := shortTTLCache.Get("ttlKey")
				Expect(found).To(BeTrue())
				Expect(value).To(Equal("ttlValue"))

				By("Verifying the key is eventually removed")
				Eventually(func() bool {
					_, found := shortTTLCache.Get("ttlKey")
					return found
				}, "500ms", "10ms").Should(BeFalse())
			})
		})

		Context("Medium", func() {
			Context("Bulk Operations", func() {
				const itemsToAdd = 50

				BeforeEach(func() {
					By("Adding multiple items to the cache")
					for i := range itemsToAdd {
						cacheClient.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
					}
				})

				It("2.1 should report correct size", func() {
					Expect(cacheClient.Size()).To(Equal(itemsToAdd))
				})

				It("1.6 should evict items when cache is full", func() {
					By("Creating small capacity cache")
					smallCache := sturdyc.New[string](10, 2, ttl, 50)

					By("Adding more items than capacity")
					for i := range 20 {
						smallCache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
					}

					By("Verifying cache size is bounded")
					Expect(smallCache.Size()).To(BeNumerically("<=", 15))
				})
			})

			Context("Metrics", func() {
				It("2.3 should track cache hits and misses via WithMetrics", func() {
					var (
						hitCount  atomic.Int32
						missCount atomic.Int32
					)

					recorder := &testMetricsRecorder{
						onHit:  func() { hitCount.Add(1) },
						onMiss: func() { missCount.Add(1) },
					}

					metricsCache := sturdyc.New[string](capacity, numShards, ttl, evictionPct,
						sturdyc.WithMetrics(recorder),
					)

					By("Setting and getting a value (hit)")
					metricsCache.Set("key1", "value1")
					sturdyc.GetOrFetch(context.Background(), metricsCache, "key1", func(_ context.Context) (string, error) {
						return "value1", nil
					})

					By("Getting a non-existent value (miss)")
					sturdyc.GetOrFetch(context.Background(), metricsCache, "missing-key", func(_ context.Context) (string, error) {
						return "fetched", nil
					})

					By("Verifying metrics were recorded")
					Expect(hitCount.Load()).To(Equal(int32(1)))
					Expect(missCount.Load()).To(Equal(int32(1)))
				})
			})

			Context("Negative Scenarios", func() {
				It("3.7 should handle empty keys", func() {
					cacheClient.Set("", "empty-key-value")
					value, found := cacheClient.Get("")
					Expect(found).To(BeTrue())
					Expect(value).To(Equal("empty-key-value"))
				})

				It("3.7 should panic on zero TTL cache", func() {
					Expect(func() {
						sturdyc.New[string](capacity, numShards, 0, evictionPct)
					}).To(Panic())
				})
			})
		})

		Context("Advanced", func() {
			// NOT SUPPORTED vs Ristretto:
			// - Direct cost-based eviction (capacity-based only)
			// - OnEvict callbacks per item
			// - Negative cost values

			// FR-002 2.5: DistributedStorage interface requires an external storage backend
			// (e.g., Redis) and is tested in play/cron_distributed_test.go with testcontainers.
			// The DistributedStorage interface (Set, SetBatch, missing record handling) is
			// designed for integration with distributed caches and is out of scope for fast tests.

			Context("Sharding", func() {
				It("should handle concurrent writes across shards", func() {
					var wg sync.WaitGroup
					const goroutines = 50
					const opsPerGoroutine = 100
					wg.Add(goroutines)

					for i := range goroutines {
						go func(id int) {
							defer wg.Done()
							for j := range opsPerGoroutine {
								key := fmt.Sprintf("shard-%d-%d", id, j)
								cacheClient.Set(key, fmt.Sprintf("value-%d-%d", id, j))
							}
						}(i)
					}
					wg.Wait()

					By("Verifying cache is functional after concurrent writes")
					cacheClient.Set("post-concurrent", "works")
					value, found := cacheClient.Get("post-concurrent")
					Expect(found).To(BeTrue())
					Expect(value).To(Equal("works"))
				})

				It("should support non-blocking reads during writes", func() {
					var wg sync.WaitGroup
					wg.Add(2)

					go func() {
						defer wg.Done()
						for i := range 1000 {
							cacheClient.Set(fmt.Sprintf("write-%d", i), "value")
						}
					}()

					go func() {
						defer wg.Done()
						for i := range 1000 {
							cacheClient.Get(fmt.Sprintf("read-%d", i))
						}
					}()

					wg.Wait()
				})
			})

			Context("Stampede Protection", func() {
				It("3.2 should coalesce concurrent fetches for the same key", func() {
					By("Creating cache with stampede protection enabled")
					stampCache := sturdyc.New[string](capacity, numShards, ttl, evictionPct)

					var fetchCount atomic.Int32
					fetchFn := func(_ context.Context) (string, error) {
						fetchCount.Add(1)
						time.Sleep(50 * time.Millisecond) // Simulate slow fetch
						return "fetched-value", nil
					}

					By("Making concurrent requests for the same key")
					var wg sync.WaitGroup
					const concurrentRequests = 10
					wg.Add(concurrentRequests)

					for range concurrentRequests {
						go func() {
							defer wg.Done()
							value, err := sturdyc.GetOrFetch(context.Background(), stampCache, "stampede-key", fetchFn)
							Expect(err).ToNot(HaveOccurred())
							Expect(value).To(Equal("fetched-value"))
						}()
					}
					wg.Wait()

					By("Verifying only one fetch was made (stampede protection)")
					Expect(fetchCount.Load()).To(Equal(int32(1)))
				})
			})

			Context("GetOrFetch", func() {
				It("3.3 should fetch on cache miss", func() {
					fetchFn := func(_ context.Context) (string, error) {
						return "fetched-value", nil
					}

					value, err := sturdyc.GetOrFetch(context.Background(), cacheClient, "fetch-key", fetchFn)
					Expect(err).ToNot(HaveOccurred())
					Expect(value).To(Equal("fetched-value"))

					By("Verifying value is cached after fetch")
					cached, found := cacheClient.Get("fetch-key")
					Expect(found).To(BeTrue())
					Expect(cached).To(Equal("fetched-value"))
				})

				It("3.3 should return cached value without fetching", func() {
					cacheClient.Set("cached-key", "cached-value")

					var fetchCalled bool
					fetchFn := func(_ context.Context) (string, error) {
						fetchCalled = true
						return "new-value", nil
					}

					value, err := sturdyc.GetOrFetch(context.Background(), cacheClient, "cached-key", fetchFn)
					Expect(err).ToNot(HaveOccurred())
					Expect(value).To(Equal("cached-value"))
					Expect(fetchCalled).To(BeFalse())
				})

				It("3.3 should handle fetch errors", func() {
					fetchFn := func(_ context.Context) (string, error) {
						return "", fmt.Errorf("fetch failed")
					}

					_, err := sturdyc.GetOrFetch(context.Background(), cacheClient, "error-key", fetchFn)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("fetch failed"))
				})
			})

			Context("Performance Benchmarks", func() {
				const (
					numOperations   = 1000
					maxTotalLatency = 50 * time.Millisecond
				)

				It("2.4 should perform set/get operations efficiently", func() {
					benchCache := sturdyc.New[string](numOperations, numShards, ttl, evictionPct)

					setStart := time.Now()
					for i := range numOperations {
						benchCache.Set(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i))
					}
					setDuration := time.Since(setStart)

					getStart := time.Now()
					for i := range numOperations {
						benchCache.Get(fmt.Sprintf("key-%d", i))
					}
					getDuration := time.Since(getStart)

					By("Verifying combined latency stays within tight budget")
					Expect(setDuration+getDuration).To(
						BeNumerically("<", maxTotalLatency),
						fmt.Sprintf("set+get operations took %s", setDuration+getDuration),
					)
				})
			})
		})
	})
})
