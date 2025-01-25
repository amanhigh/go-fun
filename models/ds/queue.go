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

func (q *Queue) Enqueue(i int) {
	q.entry.Push(i)
}

func (q *Queue) Dequeue() (i int) {
	q.transfer()
	return q.exit.Pop()
}

func (q *Queue) transfer() {
	if q.exit.IsEmpty() {
		for !q.entry.IsEmpty() {
			q.exit.Push(q.entry.Pop())
		}
	}
}

func (q *Queue) Peek() (i int) {
	q.transfer()
	return q.exit.Peek()
}
