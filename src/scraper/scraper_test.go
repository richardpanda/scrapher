package scraper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/richardpanda/scrapher/src/scraper/imdb"
)

func TestVisitURL(t *testing.T) {
	i := imdb.New("http://www.imdb.com/title/tt0468569")
	movie, err := visitURL(i)
	time.Sleep(5 * time.Second)

	assert.NoError(t, err)
	assert.NotZero(t, movie.NumRatings)
	assert.NotZero(t, movie.Rating)
	assert.Equal(t, "The Dark Knight", movie.Title)
	assert.Equal(t, "http://www.imdb.com/title/tt0468569", movie.URL)
	assert.Equal(t, 2008, movie.Year)
	assert.Equal(t, 1, len(i.Visited))
	assert.True(t, len(i.MovieIDs) > 1)
}
