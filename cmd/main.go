package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/richardpanda/scrapher/src"
	"github.com/richardpanda/scrapher/src/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

	i := scrapher.NewIMDB("http://www.imdb.com/title/tt0468569")
	rt := scrapher.NewRottenTomatoes("https://www.rottentomatoes.com/m/the_dark_knight")
	efs := []scrapher.ExtractFetcher{i, rt}

	fmt.Println("scraping")
	defer fmt.Println("done scraping")

	for _, ef := range efs {
		wg.Add(1)
		go func(ef scrapher.ExtractFetcher) {
			defer wg.Done()

			switch ef.(type) {
			case *scrapher.IMDB:
				fmt.Println("imdb scraper started")
				defer fmt.Println("imdb scraper finished")
			case *scrapher.RottenTomatoes:
				fmt.Println("rotten tomatoes scraper started")
				defer fmt.Println("rotten tomatoes scraper finished")
			default:
				log.Fatal("unknown scraper type")
			}

			s := scrapher.New(db, ef)
			s.Start()
		}(ef)
	}

	wg.Wait()
}
