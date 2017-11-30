package queue

import (
	"testing"

	"github.com/richardpanda/scrapher/node"

	"github.com/stretchr/testify/assert"
)

var (
	qn1 = &node.QueueNode{Depth: 1, MovieID: "1"}
	qn2 = &node.QueueNode{Depth: 2, MovieID: "2"}
)

func TestNew(t *testing.T) {
	q := New(qn1, qn2)
	assert.Equal(t, 2, len(q.Values))
	assert.Equal(t, qn1, q.Values[0])
	assert.Equal(t, qn2, q.Values[1])
	assert.Equal(t, 2, len(q.Visited))
	assert.True(t, q.Visited["1"])
	assert.True(t, q.Visited["2"])
}

func TestAppend(t *testing.T) {
	q := New(qn1)
	assert.Equal(t, 1, len(q.Values))
	q.Append(qn2)
	assert.Equal(t, 2, len(q.Values))
	assert.Equal(t, qn2, q.Values[1])
}

func TestHasVisited(t *testing.T) {
	q := New(qn1)
	assert.True(t, q.HasVisited("1"))
	assert.False(t, q.HasVisited("2"))
}

func TestIsEmpty(t *testing.T) {
	q := New()
	assert.True(t, q.IsEmpty())
	q = New(qn1)
	assert.False(t, q.IsEmpty())
}

func TestPop(t *testing.T) {
	q := New(qn1, qn2)
	assert.Equal(t, 2, len(q.Values))
	value := q.Pop()
	assert.Equal(t, qn1, value)
	assert.Equal(t, 1, len(q.Values))
	assert.Equal(t, qn2, q.Values[0])
}

func TestSetVisited(t *testing.T) {
	q := New()
	assert.Equal(t, 0, len(q.Visited))
	q.SetVisited("test1")
	assert.Equal(t, 1, len(q.Visited))
	assert.True(t, q.Visited["test1"])
}
