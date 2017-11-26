package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/richardpanda/scrapher/htmldoc"
)

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
