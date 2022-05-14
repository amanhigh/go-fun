package tutorial_test

import (
	"fmt"
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

		It("can compute fibonacci", func() {
			var i int
			c := make(chan int, 10)
			go fibonacci(cap(c), c)
			for i = range c {
				// For loop can detect closed channel and stop
				// fmt.Println(i)
			}
			Expect(i).To(Equal(34))
			Eventually(c).Should(BeClosed())
		})

		It("can compute fibonacci parallely", func() {
			defer GinkgoRecover()
			c := make(chan int)    // Channel to Get Fibonacci Result
			quit := make(chan int) // Channel to Signal Quit

			/** Producer (Keeps producing result until asked to quit) */
			go fibonacciMultiChannel(c, quit)

			/** Consumer */
			for i := 0; i < 10; i++ {
				<-c //Read Result so next fibonacci can be computed
			}

			Expect(<-c).To(Equal(55))

			quit <- 0 //Ask Producter to quit after reading required results.
			Eventually(c).Should(BeClosed())
		})
	})

	Context("WaitGroup", func() {
		It("should help wait", func() {
			wg := sync.WaitGroup{}
			wg.Add(2) // Starting 2 Go Routines

			go func() {
				/* Can Mark Routine start inside if routine count is not known */
				// wg.Add(1)

				//Business Logic Goes Here

				wg.Done() //Mark Job Done
			}()
			go func() {
				//Business Logic Goes Here
				wg.Done() //Mark Job Done
			}()

			wg.Wait() // Wait For Both Go Rouines are Done.
			Expect(true).To(BeTrue())
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

/* Fibonacci */
func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func fibonacciMultiChannel(c, quit chan int) {
	x, y := 0, 1
	overallTimeout := time.After(1 * time.Minute)
	for {
		select {
		case c <- x: //Write Fib to Result Channel
			x, y = y, x+y // Compute Fib
		case <-quit:
			close(c) //Close Result Channel
			return
		case <-time.After(2 * time.Second):
			fmt.Println("Operation Timeout. Operation won't wait more  than 2 Seconds.")
			return
		case <-overallTimeout:
			fmt.Println("It has been more than a minute since loop started. Returning")
			return
		default:
			// Run when no other case is ready
			time.Sleep(5 * time.Millisecond)
		}

	}
}
