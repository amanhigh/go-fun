package tutorial_test

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("RoutineFun", func() {

	Context("Mutex", func() {
		It("should protect multi threads", func() {

			key := "somekey"
			c := SafeCounter{v: make(map[string]int)}
			for i := 0; i < 500; i++ {
				go c.Inc(key)
			}
			time.Sleep(time.Millisecond * 200)
			Expect(c.Value(key)).To(Equal(500))
		})

	})

	Context("Channels", func() {
		It("can sum safely", func() {
			ints := []int{7, 2, 8, -9, 4, 0}
			/** With Buffer 2 now will work even if no goroutine is used
			  as now two responses can be buffered hence single thread won't block.
			*/
			iChannel := make(chan int, 2)
			mid := len(ints) / 2
			go sumOnChannel(ints[:mid], iChannel)
			go sumOnChannel(ints[mid:], iChannel)

			secondHalfSum, firstHalfSum := <-iChannel, <-iChannel
			Expect(firstHalfSum).To(Equal(17))
			Expect(secondHalfSum).To(Equal(-5))
		})
	})

})

/* Mutex */
// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	v   map[string]int
	mux sync.Mutex
}

// Inc increments the counter for the given key.
func (c *SafeCounter) Inc(key string) {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key]++ //Notice we are not using c.Value which will do deadlock
	c.mux.Unlock()
}

// Value returns the current value of the counter for the given key.
func (c *SafeCounter) Value(key string) int {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mux.Unlock()
	return c.v[key]
}

/* Channels */

/**
	Sum and send result only channel
	so its threadsafe.
**/
func sumOnChannel(a []int, c chan int) {
	sum := 0
	for _, x := range a {
		sum += x
	}
	c <- sum
}
