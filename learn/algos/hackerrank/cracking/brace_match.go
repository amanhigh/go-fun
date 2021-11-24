package cracking

import (
	ds "github.com/amanhigh/go-fun/apps/models/ds"
)

func MatchBrace(input string) (match bool) {
	stack := ds.NewStack()
	for _, c := range input {
		switch c {
		case '(':
			fallthrough
		case '{':
			fallthrough
		case '[':
			stack.Push(int(c))
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
