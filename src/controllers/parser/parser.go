package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/richardpanda/scrapher/src/models"
)

type Parser interface {
	AddURLs(doc *goquery.Document, urls, visited map[string]bool)
	ExtractMovieID(url string) string
	ExtractMovieInfo(doc *goquery.Document) (*models.Movie, error)
	FetchHTMLDocument(movieID string) (*goquery.Document, error)
	GetStartMovieID() string
}
