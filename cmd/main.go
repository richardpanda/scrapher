package main

import (
	"fmt"
	"log"
	"os"

	"github.com/richardpanda/scrapher/src/scrapher"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("url is missing")
	}

	url := os.Args[1]
	s := scrapher.New(url)

	for s.IsNotEmpty() {
		movie, err := s.ProcessURL()

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%s (%d)\n", movie.Title, movie.Year)
	}
}
