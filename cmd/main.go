package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/richardpanda/scrapher/src/models"
	"github.com/richardpanda/scrapher/src/scrapher"
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
	s := scrapher.New(url)

	for s.IsNotEmpty() {
		movie, err := s.ProcessURL()

		if err != nil {
			fmt.Println(err)
			continue
		}

		db.Create(movie)
	}
}
