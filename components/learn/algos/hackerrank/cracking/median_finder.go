package cracking

import (
	"github.com/amanhigh/go-fun/models/ds"
)

/*
*
https://www.youtube.com/watch?v=VmogG01IjYc
*/
type MedianFinder struct {
	lowers *ds.Heap
	uppers *ds.Heap
}

func NewMedianFinder() MedianFinder {
	maxHeap := ds.NewMaxHeap()
	minHeap := ds.NewMinHeap()
	return MedianFinder{
		lowers: &maxHeap,
		uppers: &minHeap,
	}
}

func (mf *MedianFinder) getSmallerHeap() *ds.Heap {
	if mf.lowers.Size() <= mf.uppers.Size() {
		return mf.lowers
	}
	return mf.uppers
}

func (mf *MedianFinder) getBiggerHeap() *ds.Heap {
	if mf.lowers.Size() > mf.uppers.Size() {
		return mf.lowers
	}
	return mf.uppers
}

func (mf *MedianFinder) Add(i int) {
	/* If First Half is Empty or Number is less than max of first half */
	if mf.lowers.Size() == 0 || i < mf.lowers.Peek() {
		mf.lowers.Add(i)
	} else {
		mf.uppers.Add(i)
	}
	mf.rebalance()
}

func (mf *MedianFinder) rebalance() {
	bigger := mf.getBiggerHeap()
	smaller := mf.getSmallerHeap()
	if bigger.Size()-smaller.Size() > 1 {
		smaller.Add(bigger.Poll())
	}
	// fmt.Println(mf.lowers, mf.uppers)
}

func (mf *MedianFinder) GetMedian() (result float64) {
	bigger := mf.getBiggerHeap()
	smaller := mf.getSmallerHeap()
	if bigger.Size() == smaller.Size() {
		return float64(smaller.Peek()+bigger.Peek()) / 2
	}
	return float64(bigger.Peek())
}
