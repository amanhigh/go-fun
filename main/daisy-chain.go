package main

import "fmt"

func main() {
	const n = 500000

	leftmost, rightMost := buildDaisyChain(n)

	/** Start whisper on Right Most Node */
	go func(c chan int) { c <- 1 }(rightMost)

	/** Wait for Whisped to be heard on leftMost (Other End) of Chain */
	fmt.Println(<-leftmost)
}

/*
	Build Daisy Chain of Size n.
	Return Left Most & Rightmost Node
*/
func buildDaisyChain(n int) (chan int, chan int) {
	/** Start which 1 Channel Chain, Hence left==current */
	leftmost := make(chan int)
	left := leftmost
	current := leftmost

	/** Increase Chain Size by N */
	for i := 0; i < n; i++ {
		/** Build a New Node */
		current = make(chan int)

		/** Link Nodes by making them Listen in a separate go Routine */
		go whisper(left, current)

		/** Move Current Node to Nex Node */
		left = current
	}
	return leftmost, current
}

func whisper(left, right chan int) {
	left <- (1 + <-right)
}
