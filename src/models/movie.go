package models

import (
	"github.com/jinzhu/gorm"
)

type Movie struct {
	gorm.Model
	IMDBNumRatings int     `gorm:"default:NULL"`
	IMDBRating     float64 `gorm:"default:NULL"`
	IMDBURL        string  `gorm:"default:NULL"`
	RTNumRatings   int     `gorm:"default:NULL"`
	RTRating       float64 `gorm:"default:NULL"`
	RTURL          string  `gorm:"default:NULL"`
	Title          string  `gorm:"unique_index"`
	Year           int
}
