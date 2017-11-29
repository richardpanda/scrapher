package htmldoc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var url = "http://www.imdb.com/title/tt0468569"

func TestExtractURLs(t *testing.T) {
	doc, err := Get(url)
	time.Sleep(5 * time.Second)
	assert.NoError(t, err)
	assert.NotNil(t, doc)
	urls := ExtractURLs(doc)
	assert.True(t, len(urls) > 0)
}

func TestGet(t *testing.T) {
	doc, err := Get(url)
	time.Sleep(5 * time.Second)
	assert.NoError(t, err)
	assert.NotNil(t, doc)
}
