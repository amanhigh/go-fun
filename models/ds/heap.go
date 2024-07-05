package ds

/*
*
Heap - https://www.youtube.com/watch?v=t0Cq6tVNRBA
*/
type Heap struct {
	data []int
	up   func(parent, node int) bool
	down func(left, right int) bool
}

func NewMinHeap() Heap {
	return Heap{
		data: []int{},
		up: func(parent, current int) bool {
			return parent > current
		},
		down: func(left, right int) bool {
			return left < right
		},
	}
}

func NewMaxHeap() Heap {
	return Heap{
		data: []int{},
		up: func(parent, current int) bool {
			return parent < current
		},
		down: func(left, right int) bool {
			return left > right
		},
	}
}

// 0-Parent, 1-Left Child, 2-Right Child
func (self *Heap) getLeftChildIndex(i int) int {
	return 2*i + 1
}

func (self *Heap) getRightChildIndex(i int) int {
	return 2*i + 2
}

func (self *Heap) getParentIndex(i int) int {
	return (i - 1) / 2
}

func (self *Heap) Size() int {
	return len(self.data)
}

func (self *Heap) hasLeft(i int) bool {
	return self.getLeftChildIndex(i) < self.Size()
}

func (self *Heap) hasRight(i int) bool {
	return self.getRightChildIndex(i) < self.Size()
}

func (self *Heap) hasParent(i int) bool {
	return self.getParentIndex(i) >= 0
}

func (self *Heap) left(i int) int {
	return self.data[self.getLeftChildIndex(i)]
}

func (self *Heap) right(i int) int {
	return self.data[self.getRightChildIndex(i)]
}

func (self *Heap) parent(i int) int {
	return self.data[self.getParentIndex(i)]
}

func (self *Heap) swap(i, j int) {
	self.data[i], self.data[j] = self.data[j], self.data[i]
}

func (self *Heap) Add(value int) {
	/* Add new element at end */
	self.data = append(self.data, value)
	/* Heapify Up to ensure parent child order is mantained */
	self.heapifyUp()
}

func (self *Heap) Poll() (value int) {
	if self.Size() > 0 {
		/* Return root & swap root with last element */
		value, self.data[0] = self.data[0], self.data[self.Size()-1]
		/* Shrink size as last element is now at root location */
		self.data = self.data[:self.Size()-1]
		/* Heapify down to ensure parent child order is maintained */
		self.heapifyDown()
	}
	return
}

func (self *Heap) Peek() (value int) {
	if self.Size() > 0 {
		value = self.data[0]
	}
	return
}

func (self *Heap) heapifyUp() {
	//Start from bottom till root, swapping until parent is out of order w.r.t child
	for i := self.Size() - 1; self.hasParent(i) && self.up(self.parent(i), self.data[i]); i = self.getParentIndex(i) {
		self.swap(self.getParentIndex(i), i)
	}
}

func (self *Heap) heapifyDown() {
	//Start from Root replacing node with smaller of left & right child
	for i := 0; self.hasLeft(i); {
		s := self.getLeftChildIndex(i)
		if self.hasRight(i) && self.down(self.right(i), self.left(i)) {
			s = self.getRightChildIndex(i)
		}

		//Current Node less than small child heap is ordered
		if self.down(self.data[i], self.data[s]) {
			break
		} else {
			self.swap(i, s)
			//Traverse towards smaller child to check further child heap
			i = s
		}
	}
}
