package models

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

type Movie struct {
	gorm.Model
	IMDBNumRatings sql.NullInt64
	IMDBRating     sql.NullFloat64
	IMDBURL        sql.NullString
	RTNumRatings   sql.NullInt64
	RTRating       sql.NullFloat64
	RTURL          sql.NullString
	Title          string `gorm:"unique_index"`
	Year           int
}
