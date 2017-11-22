package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/richardpanda/scrapher/src/controllers/imdb"
	"github.com/richardpanda/scrapher/src/controllers/parser"
	"github.com/richardpanda/scrapher/src/controllers/rottentomatoes"
	"github.com/richardpanda/scrapher/src/controllers/scraper"
	"github.com/richardpanda/scrapher/src/models"
)

func main() {
	var wg sync.WaitGroup
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := gorm.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	db.AutoMigrate(&models.Movie{})

	i := imdb.New("http://www.imdb.com/title/tt0468569")
	rt := rottentomatoes.New("https://www.rottentomatoes.com/m/the_dark_knight")
	parsers := []parser.Parser{i, rt}

	fmt.Println("scraping")
	defer fmt.Println("done scraping")

	for _, p := range parsers {
		wg.Add(1)
		go func(p parser.Parser) {
			defer wg.Done()

			switch p.(type) {
			case *imdb.IMDB:
				fmt.Println("imdb scraper started")
				defer fmt.Println("imdb scraper finished")
			case *rottentomatoes.RottenTomatoes:
				fmt.Println("rotten tomatoes scraper started")
				defer fmt.Println("rotten tomatoes scraper finished")
			default:
				log.Fatal("unknown scraper type")
			}

			s := scraper.New(db, p)
			s.Start()
		}(p)
	}

	wg.Wait()
}
