package scrapher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	i := NewIMDB("http://www.imdb.com/title/tt0468569")
	s := New(nil, i)

	assert.Equal(t, 1, len(s.MovieIDs))
}

func TestPop(t *testing.T) {
	i := NewIMDB("http://www.imdb.com/title/tt0468569")
	s := New(nil, i)
	movieID := s.pop()

	assert.Equal(t, 0, len(s.MovieIDs))
	assert.Equal(t, "tt0468569", movieID)
}

func TestVisitURL(t *testing.T) {
	time.Sleep(5 * time.Second)

	i := NewIMDB("http://www.imdb.com/title/tt0468569")
	s := New(nil, i)
	movie, err := s.visitURL()

	assert.NoError(t, err)
	assert.True(t, len(s.MovieIDs) > 0)
	assert.Equal(t, 1, len(s.Visited))
	assert.True(t, movie.IMDBNumRatings > 0)
	assert.True(t, movie.IMDBRating > 0)
	assert.Equal(t, "http://www.imdb.com/title/tt0468569", movie.IMDBURL)
	assert.Equal(t, "The Dark Knight", movie.Title)
	assert.Equal(t, 2008, movie.Year)
}
