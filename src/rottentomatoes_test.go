package scrapher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRottenTomatoes(t *testing.T) {
	url := "https://www.rottentomatoes.com/m/the_dark_knight"
	rt := NewRottenTomatoes(url)

	assert.Equal(t, url, rt.StartURL)
}

func TestRTHTMLDocument(t *testing.T) {
	time.Sleep(5 * time.Second)

	rt := NewRottenTomatoes("")
	doc, err := rt.HTMLDocument("the_dark_knight")

	assert.NoError(t, err)
	assert.NotNil(t, doc)
}

func TestRTMovieIDsFromDoc(t *testing.T) {
	time.Sleep(5 * time.Second)

	rt := NewRottenTomatoes("")
	doc, err := rt.HTMLDocument("the_dark_knight")
	movieIDs := rt.MovieIDsFromDoc(doc)

	assert.NoError(t, err)
	assert.True(t, len(movieIDs) > 0)
}

func TestRTMovieInfo(t *testing.T) {
	time.Sleep(5 * time.Second)

	rt := NewRottenTomatoes("")
	doc, err := rt.HTMLDocument("the_dark_knight")

	assert.NoError(t, err)

	movie, err := rt.MovieInfo(doc)

	assert.NoError(t, err)
	assert.True(t, movie.RTNumRatings > 0)
	assert.True(t, movie.RTRating > 0)
	assert.Equal(t, "https://www.rottentomatoes.com/m/the_dark_knight/", movie.RTURL)
	assert.Equal(t, "The Dark Knight", movie.Title)
	assert.Equal(t, 2008, movie.Year)
}

func TestRTStartMovieID(t *testing.T) {
	i := NewIMDB("http://www.imdb.com/title/tt0468569")
	movieID := i.StartMovieID()

	assert.Equal(t, "tt0468569", movieID)
}
