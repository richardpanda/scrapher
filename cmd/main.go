package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/richardpanda/scrapher/channel/mux"
	"github.com/richardpanda/scrapher/htmldoc"
	"github.com/richardpanda/scrapher/movie"
	"github.com/richardpanda/scrapher/queue"
)

var (
	movieIDRegex  = regexp.MustCompile(`/(tt\d{7})/?`)
	movieURLRegex = regexp.MustCompile(`^/title/tt\d{7}/\?`)
)

func appendMovieIDs(q *queue.Queue, in <-chan *goquery.Document) {
	for doc := range in {
		urls := htmldoc.ExtractURLs(doc)
		for _, url := range urls {
			if !movieURLRegex.MatchString(url) {
				continue
			}

			movieID := movieIDRegex.FindStringSubmatch(url)[1]
			if !q.HasVisited(movieID) {
				q.Append(movieID)
				q.SetVisited(movieID)
			}
		}
	}
}

func extractMovie(in <-chan *goquery.Document, out chan<- *movie.Movie, e chan<- error) {
	for doc := range in {
		m, err := movie.ExtractFromDoc(doc)
		if err != nil {
			e <- err
			continue
		}
		out <- m
	}
}

func getHTMLDocument(in <-chan string, out chan<- *goquery.Document, e chan<- error) {
	for url := range in {
		doc, err := htmldoc.Get(url)
		if err != nil {
			e <- err
			continue
		}
		out <- doc
	}
}

func insertMovie(db *gorm.DB, in <-chan *movie.Movie) {
	for m := range in {
		db.Create(m)
		fmt.Printf("added %s (%d)\n", m.Title, m.Year)
	}
}

func printError(e <-chan error) {
	for err := range e {
		fmt.Println(err)
	}
}

func main() {
	const sleepDuration = time.Duration(5) * time.Second
	var (
		dbUser           = os.Getenv("DB_USER")
		dbName           = os.Getenv("DB_NAME")
		connectionString = fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	)

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.AutoMigrate(&movie.Movie{})

	var (
		url         = "http://www.imdb.com/title/tt0468569"
		movieID     = movieIDRegex.FindStringSubmatch(url)[1]
		q           = queue.New()
		t           = time.NewTicker(sleepDuration)
		urlChan     = make(chan string)
		spreadChan  = make(chan *goquery.Document)
		appendChan  = make(chan *goquery.Document)
		extractChan = make(chan *goquery.Document)
		dbChan      = make(chan *movie.Movie)
		errChan     = make(chan error)
	)

	q.Append(movieID)
	q.SetVisited(movieID)

	go getHTMLDocument(urlChan, spreadChan, errChan)
	go mux.FanOut(spreadChan, appendChan, extractChan)
	go appendMovieIDs(q, appendChan)
	go extractMovie(extractChan, dbChan, errChan)
	go insertMovie(db, dbChan)
	go printError(errChan)

	for _ = range t.C {
		if q.IsEmpty() {
			return
		}

		movieID := q.Pop()
		url := "http://www.imdb.com/title/" + movieID
		urlChan <- url
	}
}
