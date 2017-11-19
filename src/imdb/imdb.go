package imdb

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/richardpanda/scrapher/src/models"
	"github.com/richardpanda/scrapher/src/utils"
)

type IMDB struct {
	movieIDs map[string]bool
	visited  map[string]bool
}

var (
	re            = regexp.MustCompile(`(.+) \((\d{4})\)`)
	movieIDRegex  = regexp.MustCompile(`/(tt\d{7})/?`)
	movieURLRegex = regexp.MustCompile(`^/title/tt\d{7}/\?`)
)

func New(url string) *IMDB {
	movieID := movieIDRegex.FindStringSubmatch(url)[1]

	return &IMDB{
		movieIDs: map[string]bool{movieID: true},
		visited:  map[string]bool{},
	}
}

func (i *IMDB) IsNotEmpty() bool {
	return len(i.movieIDs) != 0
}

func (i *IMDB) ProcessURL() (*models.Movie, error) {
	var movieID string

	for id := range i.movieIDs {
		movieID = id
		break
	}

	delete(i.movieIDs, movieID)

	i.visited[movieID] = true
	url := "http://www.imdb.com/title/" + movieID
	doc, err := utils.FetchHTMLDocument(url)
	time.Sleep(time.Second * 5)

	if err != nil {
		return nil, err
	}

	movie, err := extractMovieInfo(doc)

	if err != nil {
		return nil, err
	}

	nodes := doc.Find("a").Nodes

	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				if movieURLRegex.MatchString(attr.Val) {
					movieID := movieIDRegex.FindStringSubmatch(attr.Val)[1]

					if _, ok := i.visited[movieID]; !ok {
						i.movieIDs[movieID] = true
					}
				}
				break
			}
		}
	}

	return movie, nil
}

func extractMovieInfo(doc *goquery.Document) (*models.Movie, error) {
	matches := re.FindStringSubmatch(strings.TrimSpace(doc.Find("[itemprop=\"name\"]").First().Text()))

	if len(matches) < 3 {
		return nil, errors.New("unable to parse title and year")
	}

	title := matches[1]
	year, err := strconv.Atoi(matches[2])

	if err != nil {
		return nil, err
	}

	rating, err := strconv.ParseFloat(strings.Split(doc.Find("[itemprop=\"ratingValue\"]").First().Text(), "/")[0], 64)

	if err != nil {
		return nil, err
	}

	numRatings, err := utils.StringToInt(doc.Find("[itemprop=\"ratingCount\"]").First().Text())

	if err != nil {
		return nil, err
	}

	id, ok := doc.Find("meta[property=\"pageId\"]").First().Attr("content")

	if !ok {
		return nil, errors.New("cannot find movie id")
	}

	url := "http://www.imdb.com/title/" + id

	return &models.Movie{
		Title:      title,
		URL:        url,
		NumRatings: numRatings,
		Rating:     rating,
		Year:       year,
	}, nil
}
