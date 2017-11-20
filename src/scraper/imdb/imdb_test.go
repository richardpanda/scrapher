package imdb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	movieID = "tt0468569"
	url     = "http://www.imdb.com/title/tt0468569"
)

func TestNew(t *testing.T) {
	i := New(url)

	assert.Equal(t, 1, len(i.MovieIDs))
	assert.Equal(t, 0, len(i.Visited))
	assert.True(t, i.MovieIDs[movieID])
}

func TestAddURLs(t *testing.T) {
	i := New(url)
	doc, err := i.FetchHTMLDocument(movieID)
	time.Sleep(5 * time.Second)

	assert.NoError(t, err)

	i.AddURLs(doc)

	assert.True(t, len(i.MovieIDs) > 1)
}

func TestExtractMovieInfo(t *testing.T) {
	i := New(url)
	doc, err := i.FetchHTMLDocument(movieID)
	time.Sleep(5 * time.Second)

	assert.NoError(t, err)

	movie, err := i.ExtractMovieInfo(doc)

	assert.NoError(t, err)
	assert.True(t, movie.IMDBNumRatings > 0)
	assert.True(t, movie.IMDBRating > 0)
	assert.Equal(t, "http://www.imdb.com/title/tt0468569", movie.IMDBURL)
	assert.Equal(t, "The Dark Knight", movie.Title)
	assert.Equal(t, 2008, movie.Year)
}

func TestIsNotEmpty(t *testing.T) {
	i := New(url)

	assert.True(t, i.IsNotEmpty())
}

func TestFetchHTMLDocument(t *testing.T) {
	i := New(url)
	doc, err := i.FetchHTMLDocument(movieID)
	time.Sleep(5 * time.Second)

	assert.NoError(t, err)
	assert.NotNil(t, doc)
}

func TestPop(t *testing.T) {
	i := New(url)
	mid := i.Pop()

	assert.Equal(t, "tt0468569", mid)
	assert.Equal(t, 0, len(i.MovieIDs))
}

func TestSetVisited(t *testing.T) {
	i := New(url)
	i.SetVisited(movieID)

	assert.True(t, i.Visited[movieID])
}
