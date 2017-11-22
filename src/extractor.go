package scrapher

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/richardpanda/scrapher/src/models"
)

type Extractor interface {
	MovieIDsFromDoc(doc *goquery.Document) map[string]bool
	MovieInfo(doc *goquery.Document) (*models.Movie, error)
	StartMovieID() string
}
