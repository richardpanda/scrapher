package scraper

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/scrapher/src/models"
)

type Scraper interface {
	AddURLs(doc *goquery.Document)
	ExtractMovieInfo(doc *goquery.Document) (*models.Movie, error)
	FetchHTMLDocument(movieID string) (*goquery.Document, error)
	IsNotEmpty() bool
	Pop() string
	SetVisited(movieID string)
}

func Start(db *gorm.DB, s Scraper) {
	for s.IsNotEmpty() {
		movie, err := visitURL(s)
		time.Sleep(5 * time.Second)

		if err != nil {
			fmt.Println(err)
			continue
		}

		db.Create(movie)
	}
}

func visitURL(s Scraper) (*models.Movie, error) {
	movieID := s.Pop()
	s.SetVisited(movieID)
	doc, err := s.FetchHTMLDocument(movieID)

	if err != nil {
		return nil, err
	}

	movie, err := s.ExtractMovieInfo(doc)

	if err != nil {
		return nil, err
	}

	s.AddURLs(doc)
	return movie, nil
}
