package tutorial

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/tour/tree"
)

type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

func GoRoutineFun() {
	fmt.Println("\n\nGoRoutine Fun")
	treeFun()
	eventFun()
}

func treeFun() {
	fmt.Println("\n\nWalk The Tree")
	fmt.Println(Same(tree.New(2), tree.New(2)))
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
