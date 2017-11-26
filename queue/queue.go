package queue

type Queue struct {
	Values  []string
	Visited map[string]bool
}

func New() *Queue {
	return &Queue{
		Values:  []string{},
		Visited: make(map[string]bool),
	}
}

func (q *Queue) Append(value string) {
	q.Values = append(q.Values, value)
}

func (q *Queue) HasVisited(value string) bool {
	_, ok := q.Visited[value]
	return ok
}

func (q *Queue) IsEmpty() bool {
	return len(q.Values) == 0
}

func (q *Queue) Pop() string {
	value := q.Values[0]
	q.Values = q.Values[1:]
	return value
}

func (q *Queue) SetVisited(value string) {
	q.Visited[value] = true
}
