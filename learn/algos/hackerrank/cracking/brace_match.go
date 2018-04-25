package cracking

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

func MatchBrace(input string) (match bool) {
	stack := NewStack()
	for _, c := range input {
		switch c {
		case '(':
			fallthrough
		case '{':
			fallthrough
		case '[':
			stack.Push(c)
			match = true //Don't Break loop in Push
		case ')':
			match = '(' == stack.Pop()
		case '}':
			match = '{' == stack.Pop()
		case ']':
			match = '[' == stack.Pop()
		}

		//Break even if one mismatch is found
		if !match {
			break
		}
	}

	//No Mismatch found and stack is exhausted
	match = match && stack.IsEmpty()
	return
}
