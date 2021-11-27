package ds

type Stack struct {
	data []int
}

func NewStack() Stack {
	return Stack{
		data: []int{},
	}
}

func (self *Stack) Push(value int) {
	self.data = append(self.data, value)
}

func (self *Stack) Pop() (value int) {
	count := len(self.data)
	if count > 0 {
		value = self.data[count-1]
		self.data = self.data[:count-1]
	} else {
		value = -1
	}
	return
}

func (self *Stack) Peek() (value int) {
	count := len(self.data)
	if count > 0 {
		value = self.data[count-1]
	} else {
		value = -1
	}
	return
}

func (self *Stack) IsEmpty() bool {
	return len(self.data) == 0
}
