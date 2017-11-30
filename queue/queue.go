package queue

import "github.com/richardpanda/scrapher/node"

type Queue struct {
	Values  []*node.QueueNode
	Visited map[string]bool
}

func New(queueNodes ...*node.QueueNode) *Queue {
	out := &Queue{
		Values:  []*node.QueueNode{},
		Visited: make(map[string]bool),
	}

	for _, qn := range queueNodes {
		out.Append(qn)
		out.SetVisited(qn.MovieID)
	}

	return out
}

func (q *Queue) Append(qn *node.QueueNode) {
	q.Values = append(q.Values, qn)
}

func (q *Queue) HasVisited(value string) bool {
	_, ok := q.Visited[value]
	return ok
}

func (q *Queue) IsEmpty() bool {
	return len(q.Values) == 0
}

func (q *Queue) Pop() *node.QueueNode {
	qn := q.Values[0]
	q.Values = q.Values[1:]
	return qn
}

func (q *Queue) SetVisited(value string) {
	q.Visited[value] = true
}
