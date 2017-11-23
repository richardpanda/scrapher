package scrapher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewIMDB(t *testing.T) {
	url := "http://www.imdb.com/title/tt0468569"
	i := NewIMDB(url)

	assert.Equal(t, url, i.StartURL)
}

func TestIMDBHTMLDocument(t *testing.T) {
	time.Sleep(5 * time.Second)

	i := NewIMDB("")
	doc, err := i.HTMLDocument("tt0468569")

	assert.NoError(t, err)
	assert.NotNil(t, doc)
}

func TestIMDBMovieIDsFromDoc(t *testing.T) {
	time.Sleep(5 * time.Second)

	i := NewIMDB("")
	doc, err := i.HTMLDocument("tt0468569")
	movieIDs := i.MovieIDsFromDoc(doc)

	assert.NoError(t, err)
	assert.True(t, len(movieIDs) > 0)
}

func TestIMDBMovieInfo(t *testing.T) {
	time.Sleep(5 * time.Second)

	i := NewIMDB("")
	doc, err := i.HTMLDocument("tt0468569")

	assert.NoError(t, err)

	movie, err := i.MovieInfo(doc)

	assert.NoError(t, err)
	assert.True(t, movie.IMDBNumRatings > 0)
	assert.True(t, movie.IMDBRating > 0)
	assert.Equal(t, "http://www.imdb.com/title/tt0468569", movie.IMDBURL)
	assert.Equal(t, "The Dark Knight", movie.Title)
	assert.Equal(t, 2008, movie.Year)
}

func TestIMDBStartMovieID(t *testing.T) {
	i := NewIMDB("http://www.imdb.com/title/tt0468569")
	movieID := i.StartMovieID()

	assert.Equal(t, "tt0468569", movieID)
}
