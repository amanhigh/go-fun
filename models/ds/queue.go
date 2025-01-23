package ds

type Queue struct {
	entry Stack
	exit  Stack
}

func NewQueue() Queue {
	return Queue{
		entry: NewStack(),
		exit:  NewStack(),
	}
}

func (self *Queue) Enqueue(i int) {
	self.entry.Push(i)
}

func (self *Queue) Dequeue() (i int) {
	self.transfer()
	return self.exit.Pop()
}

func (self *Queue) transfer() {
	if self.exit.IsEmpty() {
		for !self.entry.IsEmpty() {
			self.exit.Push(self.entry.Pop())
		}
	}
}

func (self *Queue) Peek() (i int) {
	self.transfer()
	return self.exit.Peek()
}
