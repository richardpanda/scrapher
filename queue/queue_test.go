package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	q := New("test1", "test2")
	assert.Equal(t, 2, len(q.Values))
	assert.Equal(t, "test1", q.Values[0])
	assert.Equal(t, "test2", q.Values[1])
	assert.Equal(t, 2, len(q.Visited))
	assert.True(t, q.Visited["test1"])
	assert.True(t, q.Visited["test2"])
}

func TestAppend(t *testing.T) {
	q := New("test1")
	assert.Equal(t, 1, len(q.Values))
	q.Append("test2")
	assert.Equal(t, 2, len(q.Values))
	assert.Equal(t, "test2", q.Values[1])
}

func TestHasVisited(t *testing.T) {
	q := New("test1")
	assert.True(t, q.HasVisited("test1"))
	assert.False(t, q.HasVisited("test2"))
}

func TestIsEmpty(t *testing.T) {
	q := New()
	assert.True(t, q.IsEmpty())
	q = New("test1")
	assert.False(t, q.IsEmpty())
}

func TestPop(t *testing.T) {
	q := New("test1", "test2")
	assert.Equal(t, 2, len(q.Values))
	value := q.Pop()
	assert.Equal(t, "test1", value)
	assert.Equal(t, 1, len(q.Values))
	assert.Equal(t, "test2", q.Values[0])
}

func TestSetVisited(t *testing.T) {
	q := New()
	assert.Equal(t, 0, len(q.Visited))
	q.SetVisited("test1")
	assert.Equal(t, 1, len(q.Visited))
	assert.True(t, q.Visited["test1"])
}
