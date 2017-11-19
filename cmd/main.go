package main

import (
	"fmt"
	"log"
	"os"

	"github.com/richardpanda/scrapher/src/scraper"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/richardpanda/scrapher/src/models"
	"github.com/richardpanda/scrapher/src/scraper/imdb"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("url is missing")
	}

	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := gorm.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	db.AutoMigrate(&models.Movie{})

	url := os.Args[1]
	i := imdb.New(url)
	scraper.Start(db, i)
}
