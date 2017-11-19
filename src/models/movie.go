package models

import "github.com/jinzhu/gorm"

type Movie struct {
	gorm.Model
	NumRatings int
	Rating     float64
	Title      string `gorm:"unique_index"`
	URL        string
	Year       int
}
