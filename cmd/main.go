package main

import (
	"fmt"
	"log"

	"github.com/richardpanda/scrapher/src/scrapher"
)

func main() {
	links, err := scrapher.ExtractMovieLinks("http://www.imdb.com/sitemap/title-474.xml.gz")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(links)
	fmt.Println(len(links))
}
