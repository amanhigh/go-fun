package cracking

import . "github.com/amanhigh/go-fun/util/ds"

/**
https://www.youtube.com/watch?v=VmogG01IjYc
*/
type MedianFinder struct {
	lowers *Heap
	uppers *Heap
}

func NewMedianFinder() MedianFinder {
	maxHeap := NewMaxHeap()
	minHeap := NewMinHeap()
	return MedianFinder{
		lowers: &maxHeap,
		uppers: &minHeap,
	}
}

func (self *MedianFinder) getSmallerHeap() *Heap {
	if self.lowers.Size() <= self.uppers.Size() {
		return self.lowers
	} else {
		return self.uppers
	}
}

func (self *MedianFinder) getBiggerHeap() *Heap {
	if self.lowers.Size() > self.uppers.Size() {
		return self.lowers
	} else {
		return self.uppers
	}
}

func (self *MedianFinder) Add(i int) {
	/* If First Half is Empty or Number is less than max of first half */
	if self.lowers.Size() == 0 || i < self.lowers.Peek() {
		self.lowers.Add(i)
	} else {
		self.uppers.Add(i)
	}
	self.rebalance()
}

func (self *MedianFinder) rebalance() {
	bigger := self.getBiggerHeap()
	smaller := self.getSmallerHeap()
	if bigger.Size()-smaller.Size() > 1 {
		smaller.Add(bigger.Poll())
	}
	//fmt.Println(self.lowers, self.uppers)
}

func (self *MedianFinder) GetMedian() (result float64) {
	bigger := self.getBiggerHeap()
	smaller := self.getSmallerHeap()
	if bigger.Size() == smaller.Size() {
		return float64(smaller.Peek()+bigger.Peek()) / 2
	} else {
		return float64(bigger.Peek())
	}
}
