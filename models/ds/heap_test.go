package ds_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/models/ds"
)

var _ = Describe("Heap", func() {
	Context("MinHeap", func() {
		var minHeap ds.Heap

		BeforeEach(func() {
			minHeap = ds.NewMinHeap()
		})

		It("should build a min heap", func() {
			minHeap.Add(5)
			minHeap.Add(3)
			minHeap.Add(8)
			minHeap.Add(1)
			minHeap.Add(4)

			Expect(minHeap.Peek()).To(Equal(1))
		})
	})

	Context("MaxHeap", func() {
		var maxHeap ds.Heap

		BeforeEach(func() {
			maxHeap = ds.NewMaxHeap()
		})

		It("should build a max heap", func() {
			maxHeap.Add(5)
			maxHeap.Add(3)
			maxHeap.Add(8)
			maxHeap.Add(1)
			maxHeap.Add(4)

			Expect(maxHeap.Peek()).To(Equal(8))
		})
	})
})
