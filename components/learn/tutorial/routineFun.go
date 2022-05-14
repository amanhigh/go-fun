package tutorial

import (
	"fmt"
	"sync"
	"time"

	learn2 "github.com/amanhigh/go-fun/models/learn"

	"golang.org/x/tour/tree"
)

type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

func GoRoutineFun() {
	fmt.Println("\n\nGoRoutine Fun")
	sumFun()
	fibFun()
	treeFun()
	eventFun()
	safeRead()
}

func treeFun() {
	fmt.Println("\n\nWalk The Tree")
	fmt.Println(Same(tree.New(2), tree.New(2)))
}

func fibFun() {
	/** Fibonacci */
	fmt.Println("\n\nFibonacci")
	c := make(chan int, 10)
	go fibonacci(cap(c), c)
	for i := range c {
		// For loop can detect closed channel and stop
		fmt.Println(i)
	}
	multiChannel()
}

func sumFun() int {
	s := []int{7, 2, 8, -9, 4, 0}
	/** With Buffer 2 now will work even if no goroutine is used
	  as now two responses can be buffered hence single thread won't block.
	*/
	iChannel := make(chan int, 2)
	mid := len(s) / 2
	go sum(s[:mid], iChannel)
	go sum(s[mid:], iChannel)

	x1, x2 := <-iChannel, <-iChannel
	x3 := x1 + x2
	fmt.Printf("%v+%v=%v", x1, x2, x3)
	return x3
}

func multiChannel() {
	fmt.Println("\n\n MultiChannel Fibonacci.")
	c := make(chan int)
	quit := make(chan int)
	/** Consumer */
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()

	/** Producer */
	fibonacciMultiChannel(c, quit)
}

func fibonacciMultiChannel(c, quit chan int) {
	x, y := 0, 1
	overallTimeout := time.After(1 * time.Minute)
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		case <-time.After(2 * time.Second):
			fmt.Println("Operation Timeout. Operation won't wait more  than 2 Seconds.")
			return
		case <-overallTimeout:
			fmt.Println("It has been more than a minute since loop started. Returning")
			return
		default:
			// Run when no other case is ready
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}

	}
}

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func sum(a []int, c chan int) {
	sum := 0
	for _, x := range a {
		sum += x
	}
	c <- sum
}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	/** Inorder Traversal if Node is not null */
	if t != nil {
		Walk(t.Left, ch)
		ch <- t.Value
		Walk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	c1 := make(chan int, 5)
	c2 := make(chan int, 2)

	/** Traverse (Producers) */
	go func() {
		Walk(t1, c1)
		close(c1)
	}()
	go func() {
		Walk(t2, c2)
		close(c2)
	}()

	for y := range c1 {
		z := <-c2
		fmt.Printf("Y:%v Z:%v\n", y, z)
		if y != z {
			return false
		}
	}

	return true
}

func eventFun() {
	fmt.Println("\n\n Event Fun")
	i := 0
	intc := make(chan int, 2)

	wg := sync.WaitGroup{}
	wg.Add(2) // Starting 2 Go Routines

	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		shutdown := time.After(time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Println("Ticking (100 ms)", time.Now().UnixMilli(), i)
			case v := <-intc:
				fmt.Println("Channel Written (400 ms)", time.Now().UnixMilli(), v)
				i = v
			case <-shutdown:
				fmt.Println("Shutdown (1 Sec)", time.Now().UnixMilli(), i)
				wg.Done() //Mark Goroutine Complete
				return
			default:
				//Runs if no other event is ready
				//fmt.Println("Default", time.Now().Second())
				//time.Sleep(2 * time.Second)
			}
		}
	}()

	go func() {
		//Wait Sometime and send Channel Write
		time.Sleep(400 * time.Millisecond)
		intc <- 5
		wg.Done() //Mark Goroutine Complete
	}()

	wg.Wait()
}

func safeRead() {
	fmt.Println("\n\n Safe Read")

	safe := learn2.SafeReadWrite{1, make(chan int, 2)}

	go func() {
		//Wait Sometime and send Channel Write
		time.Sleep(400 * time.Millisecond)
		fmt.Println("Writing Now (400 ms)")
		safe.Write(5)
		safe.Close()
	}()

	for i := 0; i < 10; i++ {
		fmt.Println("Ticking (100 ms)", time.Now().UnixMilli(), safe.Read())
	}
}
