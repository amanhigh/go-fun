package tutorial_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DaisyChain", func() {
	const n = 500000

	var (
		leftMost  chan int
		rightMost chan int
	)

	BeforeEach(func() {
		leftMost, rightMost = buildDaisyChain(n)
	})

	It("should build", func() {
		Expect(leftMost).To(Not(BeNil()))
		Expect(rightMost).To(Not(BeNil()))
	})

	It("should whisper", func() {
		/**
			LeftMost (Listen) <- * <- * <- * <- * <- RightMost (Speak)
			Start whisper on Right Most Node
			If not using Go Routine it would lead into deadlock as channel is unbuffered
		 **/
		go func(c chan int) { c <- 1 }(rightMost)

		/** Wait for Whisper to be heard on leftMost (Other End) of Chain */
		Expect(<-leftMost).To(Equal(n + 1))
	})

})

/*
	Build Daisy Chain of Size n.
	Return Left Most & Rightmost Node

	Channel is unbuffered that it will block read until write happens
*/
func buildDaisyChain(n int) (leftmost chan int, current chan int) {
	/** Start which 1 Channel Chain, Hence left==current */
	leftmost = make(chan int)
	left := leftmost

	/** Increase Chain Size by N */
	for i := 0; i < n; i++ {
		/** Build a New Node */
		current = make(chan int)

		/**
		Link Nodes by making them Listen in a separate go Routine
		This will spawn routines equal to daisy chain size.
		*/
		go whisper(left, current)

		/** Move Current Node to Next Node */
		left = current
	}
	return leftmost, current
}

/**
	Read from Right Node, add 1
	write to Left Node.
 **/
func whisper(left, right chan int) {
	left <- (1 + <-right)
}
