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

	depth := flag.Int("d", -1, "web scraper depth")
	flag.Parse()

	i := imdb.New("http://www.imdb.com/title/tt0468569", *depth)
	rt := rottentomatoes.New("https://www.rottentomatoes.com/m/the_dark_knight", *depth)
	i.Init(db)
	rt.Init(db)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer fmt.Println("IMDB scraper completed!")

		for _ = range time.NewTicker(sleepDuration).C {
			if ok := i.Visit(); !ok {
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer fmt.Println("Rotten Tomatoes scraper completed!")

		for _ = range time.NewTicker(sleepDuration).C {
			if ok := rt.Visit(); !ok {
				return
			}
		}
	}()

	wg.Wait()
	fmt.Println("Finished scraping IMDB and Rotten Tomatoes!")
}
