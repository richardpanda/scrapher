package rottentomatoes

import (
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"github.com/richardpanda/scrapher/html/document"
	"github.com/richardpanda/scrapher/movie"
	"github.com/richardpanda/scrapher/node"
	"github.com/richardpanda/scrapher/queue"
)

type RottenTomatoes struct {
	*queue.Queue
	CheckDepth    bool
	QueueNodeChan chan *node.QueueNode
}

var movieIDRegex = regexp.MustCompile(`/m/(\w+)/?$`)

func New(url string, depth int) *RottenTomatoes {
	movieID := movieIDRegex.FindStringSubmatch(url)[1]
	qn := &node.QueueNode{
		Depth:   depth,
		MovieID: movieID,
	}
	return &RottenTomatoes{
		Queue:         queue.New(qn),
		CheckDepth:    depth >= 0,
		QueueNodeChan: make(chan *node.QueueNode),
	}
}

func (rt *RottenTomatoes) Init(db *gorm.DB) {
	var (
		spreadChan  = make(chan *node.DocNode)
		appendChan  = make(chan *node.DocNode)
		extractChan = make(chan *goquery.Document)
		dbChan      = make(chan *movie.Movie)
		errChan     = make(chan error)
	)

	go getHTMLDocument(rt.QueueNodeChan, spreadChan, errChan)
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
	qn := rt.Pop()
	rt.QueueNodeChan <- qn
	return true
}

func appendMovieIDs(rt *RottenTomatoes, in <-chan *node.DocNode) {
	for dn := range in {
		if rt.CheckDepth && dn.Depth-1 < 0 {
			continue
		}

		urls := document.ExtractURLs(dn.Document)
		for _, url := range urls {
			if !movieIDRegex.MatchString(url) {
				continue
			}

			movieID := movieIDRegex.FindStringSubmatch(url)[1]
			if !rt.HasVisited(movieID) {
				qn := &node.QueueNode{Depth: dn.Depth - 1, MovieID: movieID}
				rt.Append(qn)
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

func fanOut(in <-chan *node.DocNode, out1 chan<- *node.DocNode, out2 chan<- *goquery.Document) {
	for dn := range in {
		out1 <- dn
		out2 <- dn.Document
	}
}

func getHTMLDocument(in <-chan *node.QueueNode, out chan<- *node.DocNode, e chan<- error) {
	for qn := range in {
		url := "https://www.rottentomatoes.com/m/" + qn.MovieID
		doc, err := document.Get(url)
		if err != nil {
			e <- err
			continue
		}
		out <- &node.DocNode{Document: doc, QueueNode: qn}
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
