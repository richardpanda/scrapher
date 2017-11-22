package imdb

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/richardpanda/scrapher/src/models"
	"github.com/richardpanda/scrapher/src/utils"
)

type IMDB struct {
	StartURL string
}

var (
	movieTitleAndYearRegex = regexp.MustCompile(`(.+) \((\d{4})\)`)
	movieIDRegex           = regexp.MustCompile(`/(tt\d{7})/?`)
	movieURLRegex          = regexp.MustCompile(`^/title/tt\d{7}/\?`)
)

func New(url string) *IMDB {
	return &IMDB{StartURL: url}
}

func (i *IMDB) AddURLs(doc *goquery.Document, urls, visited map[string]bool) {
	nodes := doc.Find("a").Nodes

	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				if movieURLRegex.MatchString(attr.Val) {
					movieID := movieIDRegex.FindStringSubmatch(attr.Val)[1]

					if _, ok := visited[movieID]; !ok {
						urls[movieID] = true
					}
				}
				break
			}
		}
	}
}

func (i *IMDB) ExtractMovieID(url string) string {
	return movieIDRegex.FindStringSubmatch(url)[1]
}

func (i *IMDB) ExtractMovieInfo(doc *goquery.Document) (*models.Movie, error) {
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

	numRatings, err := utils.StringToInt(doc.Find("[itemprop=\"ratingCount\"]").First().Text())

	if err != nil {
		msg := fmt.Sprintf("unable to parse number of ratings from imdb (%s)", url)
		return nil, errors.New(msg)
	}

	return &models.Movie{
		IMDBNumRatings: numRatings,
		IMDBRating:     rating,
		IMDBURL:        url,
		Title:          title,
		Year:           year,
	}, nil
}

func (i *IMDB) FetchHTMLDocument(movieID string) (*goquery.Document, error) {
	url := "http://www.imdb.com/title/" + movieID

	return utils.FetchHTMLDocument(url)
}

func (i *IMDB) GetStartMovieID() string {
	return i.ExtractMovieID(i.StartURL)
}
