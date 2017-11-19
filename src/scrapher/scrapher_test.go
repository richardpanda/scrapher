package scrapher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	movieID = "tt0468569"
	url     = "http://www.imdb.com/title/tt0468569"
)

func TestFetchHTMLDocument(t *testing.T) {
	doc, err := FetchHTMLDocument(url)

	time.Sleep(5 * time.Second)

	assert.NoError(t, err)
	assert.NotNil(t, doc)
}

func TestExtractMovieInfo(t *testing.T) {
	doc, err := FetchHTMLDocument(url)
	movie, err := ExtractMovieInfo(doc)

	time.Sleep(5 * time.Second)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, movie.NumRatings)
	assert.NotEqual(t, 0, movie.Rating)
	assert.Equal(t, "The Dark Knight", movie.Title)
	assert.Equal(t, "http://www.imdb.com/title/tt0468569", movie.URL)
	assert.Equal(t, 2008, movie.Year)
}

func TestProcessURL(t *testing.T) {
	s := New(url)

	assert.Equal(t, len(s.movieIDs), 1)

	movie, err := s.ProcessURL()
	_, ok := s.visited[movieID]

	assert.NoError(t, err)
	assert.True(t, ok)
	assert.True(t, len(s.movieIDs) > 1)
	assert.Equal(t, len(s.visited), 1)
	assert.NotEqual(t, 0, movie.NumRatings)
	assert.NotEqual(t, 0, movie.Rating)
	assert.Equal(t, "The Dark Knight", movie.Title)
	assert.Equal(t, "http://www.imdb.com/title/tt0468569", movie.URL)
	assert.Equal(t, 2008, movie.Year)
}
