package movie

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
	Title          string  `gorm:"unique_index:idx_title_year"`
	Year           int     `gorm:"unique_index:idx_title_year"`
}

var movieTitleAndYearRegex = regexp.MustCompile(`(.+) \((\d{4})\)`)

func ExtractFromDoc(doc *goquery.Document) (*Movie, error) {
	id, ok := doc.Find("meta[property=\"pageId\"]").First().Attr("content")
	if !ok {
		return nil, errors.New("cannot find movie id from imdb")
	}

	url := "http://www.imdb.com/title/" + id
	str := doc.Find("[itemprop=\"name\"]").First().Text()

	matches := movieTitleAndYearRegex.FindStringSubmatch(strings.TrimSpace(str))
	if len(matches) < 3 {
		msg := fmt.Sprintf("unable to parse title and year from imdb (%s)", url)
		return nil, errors.New(msg)
	}

	title := matches[1]
	year, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, err
	}

	rating, err := strconv.ParseFloat(strings.Split(doc.Find("[itemprop=\"ratingValue\"]").First().Text(), "/")[0], 64)
	if err != nil {
		msg := fmt.Sprintf("unable to parse rating from imdb (%s)", url)
		return nil, errors.New(msg)
	}

	numRatings, err := stringToInt(doc.Find("[itemprop=\"ratingCount\"]").First().Text())
	if err != nil {
		msg := fmt.Sprintf("unable to parse number of ratings from imdb (%s)", url)
		return nil, errors.New(msg)
	}

	return &Movie{
		IMDBNumRatings: numRatings,
		IMDBRating:     rating,
		IMDBURL:        url,
		Title:          title,
		Year:           year,
	}, nil
}

func stringToInt(s string) (int, error) {
	return strconv.Atoi(strings.Replace(s, ",", "", -1))
}
