package ds_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/models/ds"
)

var _ = Describe("Queue", func() {
	var queue ds.Queue

	BeforeEach(func() {
		queue = ds.NewQueue()
	})

	It("should build", func() {
		Expect(queue).NotTo(BeNil())
	})

	Context("Enqueue", func() {
		BeforeEach(func() {
			queue.Enqueue(1)
			queue.Enqueue(2)
			queue.Enqueue(3)
		})

		AfterEach(func() {
			Expect(queue.Dequeue()).To(Equal(2))
			Expect(queue.Dequeue()).To(Equal(3))
			Expect(queue.Dequeue()).To(Equal(-1))
		})

		It("should peek at the front element without removing it", func() {
			Expect(queue.Peek()).To(Equal(1))
			Expect(queue.Peek()).To(Equal(1)) // Peek again to ensure it doesn't remove the element
			Expect(queue.Dequeue()).To(Equal(1))
			Expect(queue.Peek()).To(Equal(2))
		})
	})

	It("should handle empty queue operations", func() {
		// Assuming Dequeue returns -1 for empty queue
		Expect(queue.Dequeue()).To(Equal(-1))
		Expect(queue.Peek()).To(Equal(-1))
	})

	It("should maintain FIFO order with mixed operations", func() {
		queue.Enqueue(1)
		queue.Enqueue(2)
		Expect(queue.Dequeue()).To(Equal(1))
		queue.Enqueue(3)
		Expect(queue.Peek()).To(Equal(2))
		Expect(queue.Dequeue()).To(Equal(2))
		Expect(queue.Dequeue()).To(Equal(3))
	})

	It("should correctly handle dequeue when exit stack is empty", func() {
		queue.Enqueue(1)
		queue.Enqueue(2)
		queue.Enqueue(3)
	
		// Dequeue all elements, which should empty both stacks
		Expect(queue.Dequeue()).To(Equal(1))
		Expect(queue.Dequeue()).To(Equal(2))
		Expect(queue.Dequeue()).To(Equal(3))
	
		// Enqueue new elements
		queue.Enqueue(4)
		queue.Enqueue(5)
	
		// This dequeue should trigger the transfer from entry to exit stack
		Expect(queue.Dequeue()).To(Equal(4))
		Expect(queue.Dequeue()).To(Equal(5))
	})
})
