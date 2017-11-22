package scrapher

import (
	"github.com/PuerkitoBio/goquery"
)

type Fetcher interface {
	HTMLDocument(movieID string) (*goquery.Document, error)
}
