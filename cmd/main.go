package main

import (
	"fmt"

	"github.com/richardpanda/scrapher/src/scrapher"
)

func main() {
	s := scrapher.New("http://www.imdb.com/title/tt0468569")

	for s.IsNotEmpty() {
		movie, err := s.ProcessURL()

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%s (%d)\n", movie.Title, movie.Year)
	}
}
