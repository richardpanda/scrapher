package imdb

import (
	"testing"
	"time"

	"github.com/richardpanda/scrapher/src/utils"

	"github.com/stretchr/testify/assert"
)

var (
	movieID = "tt0468569"
	url     = "http://www.imdb.com/title/tt0468569"
)

func TestExtractMovieInfo(t *testing.T) {
	doc, err := utils.FetchHTMLDocument(url)
	movie, err := extractMovieInfo(doc)

	time.Sleep(5 * time.Second)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, movie.NumRatings)
	assert.NotEqual(t, 0, movie.Rating)
	assert.Equal(t, "The Dark Knight", movie.Title)
	assert.Equal(t, "http://www.imdb.com/title/tt0468569", movie.URL)
	assert.Equal(t, 2008, movie.Year)
}

func TestProcessURL(t *testing.T) {
	i := New(url)

	assert.Equal(t, len(i.movieIDs), 1)

	movie, err := i.ProcessURL()
	_, ok := i.visited[movieID]

	assert.NoError(t, err)
	assert.True(t, ok)
	assert.True(t, len(i.movieIDs) > 1)
	assert.Equal(t, len(i.visited), 1)
	assert.NotEqual(t, 0, movie.NumRatings)
	assert.NotEqual(t, 0, movie.Rating)
	assert.Equal(t, "The Dark Knight", movie.Title)
	assert.Equal(t, "http://www.imdb.com/title/tt0468569", movie.URL)
	assert.Equal(t, 2008, movie.Year)
}
