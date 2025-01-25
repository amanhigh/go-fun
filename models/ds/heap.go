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
func (h *Heap) getLeftChildIndex(i int) int {
	return 2*i + 1
}

func (h *Heap) getRightChildIndex(i int) int {
	return 2*i + 2
}

func (h *Heap) getParentIndex(i int) int {
	return (i - 1) / 2
}

func (h *Heap) Size() int {
	return len(h.data)
}

func (h *Heap) hasLeft(i int) bool {
	return h.getLeftChildIndex(i) < h.Size()
}

func (h *Heap) hasRight(i int) bool {
	return h.getRightChildIndex(i) < h.Size()
}

func (h *Heap) hasParent(i int) bool {
	return h.getParentIndex(i) >= 0
}

func (h *Heap) left(i int) int {
	return h.data[h.getLeftChildIndex(i)]
}

func (h *Heap) right(i int) int {
	return h.data[h.getRightChildIndex(i)]
}

func (h *Heap) parent(i int) int {
	return h.data[h.getParentIndex(i)]
}

func (h *Heap) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *Heap) Add(value int) {
	/* Add new element at end */
	h.data = append(h.data, value)
	/* Heapify Up to ensure parent child order is mantained */
	h.heapifyUp()
}

func (h *Heap) Poll() (value int) {
	if h.Size() > 0 {
		/* Return root & swap root with last element */
		value, h.data[0] = h.data[0], h.data[h.Size()-1]
		/* Shrink size as last element is now at root location */
		h.data = h.data[:h.Size()-1]
		/* Heapify down to ensure parent child order is maintained */
		h.heapifyDown()
	}
	return
}

func (h *Heap) Peek() (value int) {
	if h.Size() > 0 {
		value = h.data[0]
	}
	return
}

func (h *Heap) heapifyUp() {
	//Start from bottom till root, swapping until parent is out of order w.r.t child
	for i := h.Size() - 1; h.hasParent(i) && h.up(h.parent(i), h.data[i]); i = h.getParentIndex(i) {
		h.swap(h.getParentIndex(i), i)
	}
}

func (h *Heap) heapifyDown() {
	//Start from Root replacing node with smaller of left & right child
	for i := 0; h.hasLeft(i); {
		s := h.getLeftChildIndex(i)
		if h.hasRight(i) && h.down(h.right(i), h.left(i)) {
			s = h.getRightChildIndex(i)
		}

		//Current Node less than small child heap is ordered
		if h.down(h.data[i], h.data[s]) {
			break
		} else {
			h.swap(i, s)
			//Traverse towards smaller child to check further child heap
			i = s
		}
	}
}
