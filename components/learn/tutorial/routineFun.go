package tutorial

import (
	"fmt"

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

}
