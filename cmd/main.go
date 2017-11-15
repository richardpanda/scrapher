package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/richardpanda/scrapher/src/scrapher"
)

func main() {
	resp, err := scrapher.GetMoviePageResponse("http://www.imdb.com/title/tt0468569/?ref_=nv_sr_1")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)

	if err != nil {
		log.Fatal(err)
	}

	m := scrapher.ExtractMovieInfo(doc)
	fmt.Println(m)
}
