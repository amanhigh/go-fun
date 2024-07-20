package ds_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/models/ds"
)

var _ = FDescribe("Queue", func() {
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
})
