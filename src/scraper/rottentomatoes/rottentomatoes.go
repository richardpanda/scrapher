package rottentomatoes

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/richardpanda/scrapher/src/models"
	"github.com/richardpanda/scrapher/src/utils"
)

type RottenTomatoes struct {
	MovieIDs map[string]bool
	Visited  map[string]bool
}

var (
	movieIDRegex           = regexp.MustCompile(`/m/(\w+)/?$`)
	movieRatingRegex       = regexp.MustCompile(`\s([\d.]+)\/`)
	movieTitleAndYearRegex = regexp.MustCompile(`(.+) \((\d{4})\)`)
	numRatingsRegex        = regexp.MustCompile(`\s([\d,]+)$`)
)

func New(url string) *RottenTomatoes {
	movieID := movieIDRegex.FindStringSubmatch(url)[1]

	return &RottenTomatoes{
		MovieIDs: map[string]bool{movieID: true},
		Visited:  map[string]bool{},
	}
}

func (rt *RottenTomatoes) AddURLs(doc *goquery.Document) {
	nodes := doc.Find("a").Nodes

	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" && movieIDRegex.MatchString(attr.Val) {
				movieID := movieIDRegex.FindStringSubmatch(attr.Val)[1]

				if _, ok := rt.Visited[movieID]; !ok {
					rt.MovieIDs[movieID] = true
				}

				break
			}
		}
	}
}

func (rt *RottenTomatoes) ExtractMovieInfo(doc *goquery.Document) (*models.Movie, error) {
	s := strings.TrimSpace(doc.Find("h1.title.hidden-xs").First().Text())
	matches := movieTitleAndYearRegex.FindStringSubmatch(s)

	if len(matches) < 3 {
		return nil, errors.New("unable to parse title and year from rotten tomatoes")
	}

	title := matches[1]
	year, err := strconv.Atoi(matches[2])

	if err != nil {
		return nil, err
	}

	ratings := strings.TrimSpace(doc.Find("div.audience-info.hidden-xs.superPageFontColor").First().Text())

	if !movieRatingRegex.MatchString(ratings) {
		return nil, errors.New("unable to parse rotten tomatoes rating from rotten tomatoes")
	}

	rating, err := strconv.ParseFloat(movieRatingRegex.FindStringSubmatch(ratings)[1], 64)

	if err != nil {
		return nil, err
	}

	if !numRatingsRegex.MatchString(ratings) {
		return nil, errors.New("unable to parse number of ratings from rotten tomatoes")
	}

	numRatings, err := utils.StringToInt(numRatingsRegex.FindStringSubmatch(ratings)[1])

	if err != nil {
		return nil, err
	}

	url, ok := doc.Find(`meta[property="og:url"]`).First().Attr("content")

	if !ok {
		return nil, errors.New("unable to find url from rotten tomatoes")
	}

	return &models.Movie{
		RTNumRatings: numRatings,
		RTRating:     rating,
		RTURL:        url,
		Title:        title,
		Year:         year,
	}, nil
}

func (rt *RottenTomatoes) FetchHTMLDocument(movieID string) (*goquery.Document, error) {
	url := "https://www.rottentomatoes.com/m/" + movieID
	return utils.FetchHTMLDocument(url)
}

func (rt *RottenTomatoes) IsNotEmpty() bool {
	return len(rt.MovieIDs) != 0
}

func (rt *RottenTomatoes) Pop() string {
	var movieID string

	for id := range rt.MovieIDs {
		movieID = id
		break
	}

	delete(rt.MovieIDs, movieID)

	return movieID
}

func (rt *RottenTomatoes) SetVisited(movieID string) {
	rt.Visited[movieID] = true
}
