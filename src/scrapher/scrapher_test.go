package scrapher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var url = "http://www.imdb.com/title/tt0468569"

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
