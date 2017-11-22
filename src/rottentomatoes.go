package scrapher

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
	rtMovieIDRegex           = regexp.MustCompile(`/m/(\w+)/?$`)
	rtMovieRatingRegex       = regexp.MustCompile(`\s([\d.]+)\/`)
	rtMovieTitleAndYearRegex = regexp.MustCompile(`(.+) \((\d{4})\)`)
	rtNumRatingsRegex        = regexp.MustCompile(`\s([\d,]+)$`)
)

func NewRottenTomatoes(url string) *RottenTomatoes {
	return &RottenTomatoes{StartURL: url}
}

func (rt *RottenTomatoes) HTMLDocument(movieID string) (*goquery.Document, error) {
	url := "https://www.rottentomatoes.com/m/" + movieID
	return utils.FetchHTMLDocument(url)
}

func (rt *RottenTomatoes) MovieIDsFromDoc(doc *goquery.Document) map[string]bool {
	movieIDs := make(map[string]bool)
	nodes := doc.Find("a").Nodes
	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" && rtMovieIDRegex.MatchString(attr.Val) {
				if rtMovieIDRegex.MatchString(attr.Val) {
					movieID := rtMovieIDRegex.FindStringSubmatch(attr.Val)[1]
					movieIDs[movieID] = true
				}
				break
			}
		}
	}
	return movieIDs
}

func (rt *RottenTomatoes) MovieInfo(doc *goquery.Document) (*models.Movie, error) {
	url, ok := doc.Find(`meta[property="og:url"]`).First().Attr("content")
	if !ok {
		return nil, errors.New("unable to find url from rotten tomatoes")
	}

	s := strings.TrimSpace(doc.Find("h1.title.hidden-xs").First().Text())
	matches := rtMovieTitleAndYearRegex.FindStringSubmatch(s)
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
	if !rtMovieRatingRegex.MatchString(ratings) {
		msg := fmt.Sprintf("unable to find ratings from rotten tomatoes (%s)", url)
		return nil, errors.New(msg)
	}

	rating, err := strconv.ParseFloat(rtMovieRatingRegex.FindStringSubmatch(ratings)[1], 64)
	if err != nil {
		msg := fmt.Sprintf("unable to parse ratings from rotten tomatoes (%s)", url)
		return nil, errors.New(msg)
	}

	if !rtNumRatingsRegex.MatchString(ratings) {
		msg := fmt.Sprintf("unable to find number of ratings from rotten tomatoes (%s)", url)
		return nil, errors.New(msg)
	}
	numRatings, err := utils.StringToInt(rtNumRatingsRegex.FindStringSubmatch(ratings)[1])
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

func (rt *RottenTomatoes) StartMovieID() string {
	return rtMovieIDRegex.FindStringSubmatch(rt.StartURL)[1]
}
