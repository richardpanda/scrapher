package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/richardpanda/scrapher/imdb"
	"github.com/richardpanda/scrapher/movie"
	"github.com/richardpanda/scrapher/rottentomatoes"
)

func main() {
	const sleepDuration = time.Duration(5) * time.Second
	var (
		dbUser           = os.Getenv("DB_USER")
		dbName           = os.Getenv("DB_NAME")
		connectionString = fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
		wg               sync.WaitGroup
	)

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&movie.Movie{})

	// imdb: http://www.imdb.com/title/tt0468569
	// rt: https://www.rottentomatoes.com/m/the_dark_knight
	imdbURL := flag.String("imdb", "", "imdb start url")
	rtURL := flag.String("rt", "", "rotten tomatoes start url")
	depth := flag.Int("d", -1, "web scraper depth")
	flag.Parse()

	if *imdbURL == "" && *rtURL == "" {
		log.Fatal("specify imdb flag and/or rt flag")
	}

	if *imdbURL != "" {
		i := imdb.New(*imdbURL, *depth)
		i.Init(db)

		wg.Add(1)
		go func() {
			fmt.Println("imdb scraper has started!")
			defer wg.Done()
			defer fmt.Println("imdb scraper completed!")

			for _ = range time.NewTicker(sleepDuration).C {
				if ok := i.Visit(); !ok {
					return
				}
			}
		}()
	}

	if *rtURL != "" {
		rt := rottentomatoes.New(*rtURL, *depth)
		rt.Init(db)

		wg.Add(1)
		go func() {
			fmt.Println("rotten tomatoes scraper has started!")
			defer wg.Done()
			defer fmt.Println("rotten tomatoes scraper completed!")

			for _ = range time.NewTicker(sleepDuration).C {
				if ok := rt.Visit(); !ok {
					return
				}
			}
		}()
	}

	wg.Wait()
}
