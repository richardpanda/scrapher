package rottentomatoes

import (
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/scrapher/html/document"
	"github.com/richardpanda/scrapher/movie"
	"github.com/richardpanda/scrapher/queue"
)

type RottenTomatoes struct {
	*queue.Queue
	URLChan chan string
}

var movieIDRegex = regexp.MustCompile(`/m/(\w+)/?$`)

func New(url string) *RottenTomatoes {
	movieID := movieIDRegex.FindStringSubmatch(url)[1]
	return &RottenTomatoes{
		Queue:   queue.New(movieID),
		URLChan: make(chan string),
	}
}

func (rt *RottenTomatoes) Init(db *gorm.DB) {
	var (
		spreadChan  = make(chan *goquery.Document)
		appendChan  = make(chan *goquery.Document)
		extractChan = make(chan *goquery.Document)
		dbChan      = make(chan *movie.Movie)
		errChan     = make(chan error)
	)

	go getHTMLDocument(rt.URLChan, spreadChan, errChan)
	go fanOut(spreadChan, appendChan, extractChan)
	go appendMovieIDs(rt, appendChan)
	go extractMovie(extractChan, dbChan, errChan)
	go insertMovie(db, dbChan)
	go printError(errChan)
}

func (rt *RottenTomatoes) Visit() bool {
	if rt.IsEmpty() {
		return false
	}
	movieID := rt.Pop()
	url := "https://www.rottentomatoes.com/m/" + movieID
	rt.URLChan <- url
	return true
}

func appendMovieIDs(rt *RottenTomatoes, in <-chan *goquery.Document) {
	for doc := range in {
		urls := document.ExtractURLs(doc)
		for _, url := range urls {
			if !movieIDRegex.MatchString(url) {
				continue
			}

			movieID := movieIDRegex.FindStringSubmatch(url)[1]
			if !rt.HasVisited(movieID) {
				rt.Append(movieID)
				rt.SetVisited(movieID)
			}
		}
	}
}

func extractMovie(in <-chan *goquery.Document, out chan<- *movie.Movie, e chan<- error) {
	for doc := range in {
		m, err := movie.ExtractFromRT(doc)
		if err != nil {
			e <- err
			continue
		}
		out <- m
	}
}

func fanOut(in <-chan *goquery.Document, out1, out2 chan<- *goquery.Document) {
	for doc := range in {
		out1 <- doc
		out2 <- doc
	}
}

func getHTMLDocument(in <-chan string, out chan<- *goquery.Document, e chan<- error) {
	for url := range in {
		doc, err := document.Get(url)
		if err != nil {
			e <- err
			continue
		}
		out <- doc
	}
}

func insertMovie(db *gorm.DB, in <-chan *movie.Movie) {
	for m := range in {
		mov := &movie.Movie{}
		result := db.Where("title = ? AND year = ?", m.Title, m.Year).First(mov)
		if result.RowsAffected == 1 {
			db.Model(mov).Update(m)
			fmt.Printf("updated %s (%d)\n", m.Title, m.Year)
		} else {
			db.Create(m)
			fmt.Printf("added %s (%d)\n", m.Title, m.Year)
		}
	}
}

func printError(e <-chan error) {
	for err := range e {
		fmt.Println(err)
	}
}
