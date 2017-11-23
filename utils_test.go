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
