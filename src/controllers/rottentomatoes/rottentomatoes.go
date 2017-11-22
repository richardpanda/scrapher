package rottentomatoes

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

type RottenTomatoes struct {
	StartURL string
}

var (
	movieIDRegex           = regexp.MustCompile(`/m/(\w+)/?$`)
	movieRatingRegex       = regexp.MustCompile(`\s([\d.]+)\/`)
	movieTitleAndYearRegex = regexp.MustCompile(`(.+) \((\d{4})\)`)
	numRatingsRegex        = regexp.MustCompile(`\s([\d,]+)$`)
)

func New(url string) *RottenTomatoes {
	return &RottenTomatoes{StartURL: url}
}

func (rt *RottenTomatoes) AddURLs(doc *goquery.Document, urls, visited map[string]bool) {
	nodes := doc.Find("a").Nodes

	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" && movieIDRegex.MatchString(attr.Val) {
				movieID := movieIDRegex.FindStringSubmatch(attr.Val)[1]

				if _, ok := visited[movieID]; !ok {
					urls[movieID] = true
				}

				break
			}
		}
	}
}

func (rt *RottenTomatoes) ExtractMovieID(url string) string {
	return movieIDRegex.FindStringSubmatch(url)[1]
}

func (rt *RottenTomatoes) ExtractMovieInfo(doc *goquery.Document) (*models.Movie, error) {
	url, ok := doc.Find(`meta[property="og:url"]`).First().Attr("content")

	if !ok {
		return nil, errors.New("unable to find url from rotten tomatoes")
	}

	s := strings.TrimSpace(doc.Find("h1.title.hidden-xs").First().Text())
	matches := movieTitleAndYearRegex.FindStringSubmatch(s)

	if len(matches) < 3 {
		msg := fmt.Sprintf("unable to parse title and year from rotten tomatoes (%s)", url)
		return nil, errors.New(msg)
	}

	title := matches[1]
	year, err := strconv.Atoi(matches[2])

	if err != nil {
		return nil, err
	}

	ratings := strings.TrimSpace(doc.Find("div.audience-info.hidden-xs.superPageFontColor").First().Text())

	if !movieRatingRegex.MatchString(ratings) {
		msg := fmt.Sprintf("unable to find ratings from rotten tomatoes (%s)", url)
		return nil, errors.New(msg)
	}

	rating, err := strconv.ParseFloat(movieRatingRegex.FindStringSubmatch(ratings)[1], 64)

	if err != nil {
		msg := fmt.Sprintf("unable to parse ratings from rotten tomatoes (%s)", url)
		return nil, errors.New(msg)
	}

	if !numRatingsRegex.MatchString(ratings) {
		msg := fmt.Sprintf("unable to find number of ratings from rotten tomatoes (%s)", url)
		return nil, errors.New(msg)
	}

	numRatings, err := utils.StringToInt(numRatingsRegex.FindStringSubmatch(ratings)[1])

	if err != nil {
		msg := fmt.Sprintf("unable to parse number of ratings from rotten tomatoes (%s)", url)
		return nil, errors.New(msg)
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

func (rt *RottenTomatoes) GetStartMovieID() string {
	return rt.ExtractMovieID(rt.StartURL)
}
