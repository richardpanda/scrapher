package main

import (
	"github.com/richardpanda/scrapher/src/scrapher"
)

func main() {
	scrapher.StartFromURL("http://www.imdb.com/title/tt0468569/?ref_=rvi_tt")
}
