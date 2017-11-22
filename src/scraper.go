package scrapher

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/richardpanda/scrapher/src/models"
)

type Scraper struct {
	DB *gorm.DB
	ExtractFetcher
	MovieIDs map[string]bool
	Visited  map[string]bool
}

func New(db *gorm.DB, ef ExtractFetcher) *Scraper {
	movieID := ef.StartMovieID()
	return &Scraper{
		DB:             db,
		ExtractFetcher: ef,
		MovieIDs:       map[string]bool{movieID: true},
		Visited:        make(map[string]bool),
	}
}

func (s *Scraper) Start() {
	for len(s.MovieIDs) != 0 {
		movie, err := s.visitURL()
		time.Sleep(5 * time.Second)
		if err != nil {
			fmt.Println(err)
			continue
		}

		m := &models.Movie{}
		result := s.DB.Where("title = ? AND year = ?", movie.Title, movie.Year).First(m)
		if result.RowsAffected == 1 {
			s.DB.Model(m).Update(movie)
			fmt.Printf("updated %s (%d)\n", movie.Title, movie.Year)
		} else {
			s.DB.Create(movie)
			fmt.Printf("added %s (%d)\n", movie.Title, movie.Year)
		}
	}
}

func (s *Scraper) pop() string {
	var movieID string
	for id := range s.MovieIDs {
		movieID = id
		break
	}
	delete(s.MovieIDs, movieID)
	return movieID
}

func (s *Scraper) visitURL() (*models.Movie, error) {
	movieID := s.pop()
	s.Visited[movieID] = true

	doc, err := s.HTMLDocument(movieID)
	if err != nil {
		return nil, err
	}

	movie, err := s.MovieInfo(doc)
	if err != nil {
		return nil, err
	}

	movieIDs := s.MovieIDsFromDoc(doc)
	for movieID := range movieIDs {
		if _, ok := s.Visited[movieID]; !ok {
			s.MovieIDs[movieID] = true
		}
	}

	return movie, nil
}
