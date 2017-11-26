package mux

import "github.com/PuerkitoBio/goquery"

func FanOut(in <-chan *goquery.Document, out1, out2 chan<- *goquery.Document) {
	for doc := range in {
		out1 <- doc
		out2 <- doc
	}
}
