package rottentomatoes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	movieID = "the_dark_knight"
	url     = "https://www.rottentomatoes.com/m/the_dark_knight"
)

func TestNew(t *testing.T) {
	rt := New(url)

	assert.Equal(t, 1, len(rt.MovieIDs))
	assert.Equal(t, 0, len(rt.Visited))
	assert.True(t, rt.MovieIDs["the_dark_knight"])
}

func TestAddURLs(t *testing.T) {
	rt := New(url)
	doc, err := rt.FetchHTMLDocument(movieID)
	time.Sleep(5 * time.Second)

	assert.NoError(t, err)

	rt.AddURLs(doc)

	assert.True(t, len(rt.MovieIDs) > 1)
}

func TestExtractMovieInfo(t *testing.T) {
	rt := New(url)
	doc, err := rt.FetchHTMLDocument(movieID)
	time.Sleep(5 * time.Second)

	assert.NoError(t, err)

	movie, err := rt.ExtractMovieInfo(doc)

	assert.NoError(t, err)
	assert.True(t, movie.RTNumRatings > 0)
	assert.True(t, movie.RTRating > 0)
	assert.Equal(t, "https://www.rottentomatoes.com/m/the_dark_knight/", movie.RTURL)
	assert.Equal(t, "The Dark Knight", movie.Title)
	assert.Equal(t, 2008, movie.Year)
}

func TestFetchHTMLDocument(t *testing.T) {
	rt := New(url)
	doc, err := rt.FetchHTMLDocument(movieID)
	time.Sleep(5 * time.Second)

	assert.NoError(t, err)
	assert.NotNil(t, doc)
}

func TestIsNotEmpty(t *testing.T) {
	rt := New(url)

	assert.True(t, rt.IsNotEmpty())
}

func TestPop(t *testing.T) {
	rt := New(url)
	movieID := rt.Pop()

	assert.Equal(t, "the_dark_knight", movieID)
	assert.Equal(t, 0, len(rt.MovieIDs))
}

func TestSetVisited(t *testing.T) {
	rt := New(url)
	rt.SetVisited(movieID)

	assert.True(t, rt.Visited[movieID])
}
