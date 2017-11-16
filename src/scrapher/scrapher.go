package scrapher

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/richardpanda/scrapher/src/models"
	"github.com/richardpanda/scrapher/src/utils"
)

var re = regexp.MustCompile(`(.+) \((\d{4})\)`)

func ExtractMovieInfo(doc *goquery.Document) (*models.Movie, error) {
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

	return &models.Movie{
		Title:      title,
		Year:       year,
		Rating:     rating,
		NumRatings: numRatings,
	}, nil
}

func ExtractMovieLinks(url string) ([]string, error) {
	resp, err := GetHTTPResponse(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	u := models.URLSet{}
	err = xml.Unmarshal(b, &u)

	if err != nil {
		return nil, err
	}

	links := make([]string, len(u.URLs))

	for i, url := range u.URLs {
		links[i] = url.Location
	}

	return links, nil
}

func ExtractSitemapLinks() ([]string, error) {
	url := "http://www.imdb.com/sitemap/index.xml.gz"
	resp, err := GetHTTPResponse(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	s := models.SitemapIndex{}
	err = xml.Unmarshal(b, &s)

	if err != nil {
		return nil, err
	}

	links := []string{}

	for _, sitemap := range s.Sitemaps {
		if strings.HasPrefix(sitemap.Location, "http://www.imdb.com/sitemap/title") {
			links = append(links, sitemap.Location)
		}
	}

	return links, nil
}

func GetHTTPResponse(url string) (*http.Response, error) {
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

func StartFromSitemap() {
	sitemapLinks, err := ExtractSitemapLinks()

	if err != nil {
		log.Fatal(err)
	}

	for _, sitemapLink := range sitemapLinks {
		time.Sleep(time.Second * 5)
		movieLinks, err := ExtractMovieLinks(sitemapLink)

		if err != nil {
			log.Fatal(err)
		}

		for _, movieLink := range movieLinks {
			time.Sleep(time.Second * 5)
			resp, err := GetHTTPResponse(movieLink)

			if err != nil {
				log.Fatal(err)
			}

			doc, err := goquery.NewDocumentFromResponse(resp)

			if err != nil {
				log.Fatal(err)
			}

			movie, err := ExtractMovieInfo(doc)

			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(movie)
			resp.Body.Close()
		}
	}
}
