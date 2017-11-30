package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/richardpanda/scrapher/imdb"
	"github.com/richardpanda/scrapher/movie"
)

func main() {
	const sleepDuration = time.Duration(5) * time.Second
	var (
		dbUser           = os.Getenv("DB_USER")
		dbName           = os.Getenv("DB_NAME")
		connectionString = fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	)

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.AutoMigrate(&movie.Movie{})

	ticker := time.NewTicker(sleepDuration)
	i := imdb.New("http://www.imdb.com/title/tt0468569")
	i.Init(db)

	for _ = range ticker.C {
		if ok := i.Visit(); !ok {
			fmt.Println("IMDB scraper completed!")
			return
		}
	}
}
