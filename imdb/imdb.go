package imdb

import (
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/scrapher/html/document"
	"github.com/richardpanda/scrapher/movie"
	"github.com/richardpanda/scrapher/queue"
)

type IMDB struct {
	*queue.Queue
	URLChan chan string
}

var (
	movieIDRegex  = regexp.MustCompile(`/(tt\d{7})/?`)
	movieURLRegex = regexp.MustCompile(`^/title/tt\d{7}/\?`)
)

func New(url string) *IMDB {
	movieID := movieIDRegex.FindStringSubmatch(url)[1]
	return &IMDB{
		Queue:   queue.New(movieID),
		URLChan: make(chan string),
	}
}

func (i *IMDB) Init(db *gorm.DB) {
	var (
		spreadChan  = make(chan *goquery.Document)
		appendChan  = make(chan *goquery.Document)
		extractChan = make(chan *goquery.Document)
		dbChan      = make(chan *movie.Movie)
		errChan     = make(chan error)
	)

	go getHTMLDocument(i.URLChan, spreadChan, errChan)
	go fanOut(spreadChan, appendChan, extractChan)
	go appendMovieIDs(i, appendChan)
	go extractMovie(extractChan, dbChan, errChan)
	go insertMovie(db, dbChan)
	go printError(errChan)
}

func (i *IMDB) Visit() bool {
	if i.IsEmpty() {
		return false
	}
	movieID := i.Pop()
	url := "http://www.imdb.com/title/" + movieID
	i.URLChan <- url
	return true
}

func appendMovieIDs(i *IMDB, in <-chan *goquery.Document) {
	for doc := range in {
		urls := document.ExtractURLs(doc)
		for _, url := range urls {
			if !movieURLRegex.MatchString(url) {
				continue
			}

			movieID := movieIDRegex.FindStringSubmatch(url)[1]
			if !i.HasVisited(movieID) {
				i.Append(movieID)
				i.SetVisited(movieID)
			}
		}
	}
}

func extractMovie(in <-chan *goquery.Document, out chan<- *movie.Movie, e chan<- error) {
	for doc := range in {
		m, err := movie.ExtractFromIMDB(doc)
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
