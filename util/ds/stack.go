package ds

type Stack struct {
	data []rune
}

func NewStack() Stack {
	return Stack{
		data: []rune{},
	}
}

func (self *Stack) Push(value rune) {
	self.data = append(self.data, value)
}

func (self *Stack) Pop() (value rune) {
	count := len(self.data)
	if count > 0 {
		value = self.data[count-1]
		self.data = self.data[:count-1]
	} else {
		value = '#'
	}
	return
}

func (self *Stack) IsEmpty() bool {
	return len(self.data) == 0
}
