package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/richardpanda/scrapher/channel/mux"
	"github.com/richardpanda/scrapher/movie"
	"github.com/richardpanda/scrapher/queue"
)

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
		imdbURL     = "http://www.imdb.com/title/tt0468569"
		imdbMovieID = imdbMovieIDRegex.FindStringSubmatch(imdbURL)[1]
		rtURL       = "https://www.rottentomatoes.com/m/the_dark_knight"
		rtMovieID   = rtMovieIDRegex.FindStringSubmatch(rtURL)[1]
		app         = &App{
			imdb: queue.New(imdbMovieID),
			rt:   queue.New(rtMovieID),
		}
		t               = time.NewTicker(sleepDuration)
		urlChan         = make(chan string)
		spreadChan      = make(chan *goquery.Document)
		imdbAppendChan  = make(chan *goquery.Document)
		imdbExtractChan = make(chan *goquery.Document)
		rtAppendChan    = make(chan *goquery.Document)
		rtExtractChan   = make(chan *goquery.Document)
		dbChan          = make(chan *movie.Movie)
		errChan         = make(chan error)
		m               = &mux.Mux{
			Router: map[string][]chan *goquery.Document{
				"www.imdb.com":           []chan *goquery.Document{imdbAppendChan, imdbExtractChan},
				"www.rottentomatoes.com": []chan *goquery.Document{rtAppendChan, rtExtractChan},
			},
		}
	)

	go getHTMLDocument(urlChan, spreadChan, errChan)

	go m.FanOut(spreadChan)

	go appendMovieIDsFromIMDB(app, imdbAppendChan)
	go extractMovieFromIMDB(imdbExtractChan, dbChan, errChan)

	go appendMovieIDsFromRT(app, rtAppendChan)
	go extractMovieFromRT(rtExtractChan, dbChan, errChan)

	go insertMovie(db, dbChan)

	go printError(errChan)

	for _ = range t.C {
		if app.imdb.IsEmpty() && app.rt.IsEmpty() {
			return
		}

		if !app.imdb.IsEmpty() {
			movieID := app.imdb.Pop()
			url := "http://www.imdb.com/title/" + movieID
			urlChan <- url
		}

		if !app.rt.IsEmpty() {
			movieID := app.rt.Pop()
			url := "https://www.rottentomatoes.com/m/" + movieID
			urlChan <- url
		}
	}
}
