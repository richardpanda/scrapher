package scrapher

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/richardpanda/scrapher/src/models"
	"github.com/richardpanda/scrapher/src/utils"
)

var re = regexp.MustCompile(`(.+) \((\d{4})\)`)

func GetMoviePageResponse(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Scrapher, a friendly web scraper. Code can be found at https://github.com/richardpanda/scrapher.")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, err
}

func ExtractMovieInfo(doc *goquery.Document) models.Movie {
	matches := re.FindStringSubmatch(strings.TrimSpace(doc.Find("[itemprop=\"name\"]").First().Text()))

	title := matches[1]
	year, _ := strconv.Atoi(matches[2])
	rating, _ := strconv.ParseFloat(strings.Split(doc.Find("[itemprop=\"ratingValue\"]").First().Text(), "/")[0], 64)
	numRatings, _ := utils.StringToInt(doc.Find("[itemprop=\"ratingCount\"]").First().Text())

	return models.Movie{
		Title:      title,
		Year:       year,
		Rating:     rating,
		NumRatings: numRatings,
	}
}
