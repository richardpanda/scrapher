package mux

import (
	"github.com/PuerkitoBio/goquery"
)

type Mux struct {
	Router map[string][]chan *goquery.Document
}

func (m *Mux) FanOut(in <-chan *goquery.Document) {
	for doc := range in {
		chs := m.Router[doc.Url.Hostname()]
		for _, ch := range chs {
			ch <- doc
		}
	}
}
