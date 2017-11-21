package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/richardpanda/scrapher/src/models"
	"github.com/richardpanda/scrapher/src/scraper"
	"github.com/richardpanda/scrapher/src/scraper/rottentomatoes"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := gorm.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	db.AutoMigrate(&models.Movie{})

	// i := imdb.New("http://www.imdb.com/title/tt0468569")
	// scraper.Start(db, i)
	rt := rottentomatoes.New("https://www.rottentomatoes.com/m/the_dark_knight")
	scraper.Start(db, rt)
}
