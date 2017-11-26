package main

import (
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/scrapher/htmldoc"
	"github.com/richardpanda/scrapher/movie"
)

var (
	imdbMovieIDRegex  = regexp.MustCompile(`/(tt\d{7})/?`)
	imdbMovieURLRegex = regexp.MustCompile(`^/title/tt\d{7}/\?`)
	rtMovieIDRegex    = regexp.MustCompile(`/m/(\w+)/?$`)
)

func appendMovieIDsFromIMDB(a *App, in <-chan *goquery.Document) {
	for doc := range in {
		urls := htmldoc.ExtractURLs(doc)
		for _, url := range urls {
			if !imdbMovieURLRegex.MatchString(url) {
				continue
			}

			movieID := imdbMovieIDRegex.FindStringSubmatch(url)[1]
			if !a.imdb.HasVisited(movieID) {
				a.imdb.Append(movieID)
				a.imdb.SetVisited(movieID)
			}
		}
	}
}

func appendMovieIDsFromRT(a *App, in <-chan *goquery.Document) {
	for doc := range in {
		urls := htmldoc.ExtractURLs(doc)
		for _, url := range urls {
			if !rtMovieIDRegex.MatchString(url) {
				continue
			}

			movieID := rtMovieIDRegex.FindStringSubmatch(url)[1]
			if !a.rt.HasVisited(movieID) {
				a.rt.Append(movieID)
				a.rt.SetVisited(movieID)
			}
		}
	}
}

func extractMovieFromIMDB(in <-chan *goquery.Document, out chan<- *movie.Movie, e chan<- error) {
	for doc := range in {
		m, err := movie.ExtractFromIMDB(doc)
		if err != nil {
			e <- err
			continue
		}
		out <- m
	}
}

func extractMovieFromRT(in <-chan *goquery.Document, out chan<- *movie.Movie, e chan<- error) {
	for doc := range in {
		m, err := movie.ExtractFromRT(doc)
		if err != nil {
			e <- err
			continue
		}
		out <- m
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
