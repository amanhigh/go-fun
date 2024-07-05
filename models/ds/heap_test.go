package ds_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/models/ds"
)

var _ = Describe("Heap", func() {
	Context("Min", func() {
		var minHeap ds.Heap

		BeforeEach(func() {
			minHeap = ds.NewMinHeap()
		})

		Context("Peek", func() {
			It("should peek the smallest element", func() {
				minHeap.Add(5)
				minHeap.Add(3)
				minHeap.Add(8)
				minHeap.Add(1)
				minHeap.Add(4)

				Expect(minHeap.Peek()).To(Equal(1))
			})
		})

		Context("Poll", func() {
			It("should poll the smallest element", func() {
				minHeap.Add(5)
				minHeap.Add(3)
				minHeap.Add(8)
				minHeap.Add(1)
				minHeap.Add(4)

				Expect(minHeap.Poll()).To(Equal(1))
				Expect(minHeap.Peek()).To(Equal(3))
			})
		})
	})

	Context("Max", func() {
		var maxHeap ds.Heap

		BeforeEach(func() {
			maxHeap = ds.NewMaxHeap()
		})

		Context("Peek", func() {
			It("should peek the largest element", func() {
				maxHeap.Add(5)
				maxHeap.Add(3)
				maxHeap.Add(8)
				maxHeap.Add(1)
				maxHeap.Add(4)

				Expect(maxHeap.Peek()).To(Equal(8))
			})
		})

		Context("Poll", func() {
			It("should poll the largest element", func() {
				maxHeap.Add(5)
				maxHeap.Add(3)
				maxHeap.Add(8)
				maxHeap.Add(1)
				maxHeap.Add(4)

				Expect(maxHeap.Poll()).To(Equal(8))
				Expect(maxHeap.Peek()).To(Equal(5))
			})
		})
	})
})
