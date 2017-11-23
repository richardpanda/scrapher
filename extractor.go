package scrapher

import (
	"github.com/PuerkitoBio/goquery"
)

type Extractor interface {
	MovieIDsFromDoc(doc *goquery.Document) map[string]bool
	MovieInfo(doc *goquery.Document) (*Movie, error)
	StartMovieID() string
}
