package ds

type Stack struct {
	data []int
}

func NewStack() Stack {
	return Stack{
		data: []int{},
	}
}

func (s *Stack) Push(value int) {
	s.data = append(s.data, value)
}

func (s *Stack) Pop() (value int) {
	count := len(s.data)
	if count > 0 {
		value = s.data[count-1]
		s.data = s.data[:count-1]
	} else {
		value = -1
	}
	return
}

func (s *Stack) Peek() (value int) {
	count := len(s.data)
	if count > 0 {
		value = s.data[count-1]
	} else {
		value = -1
	}
	return
}

func (s *Stack) IsEmpty() bool {
	return len(s.data) == 0
}
