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

type Scrapher struct {
	movieIDs map[string]bool
	visited  map[string]bool
}

var (
	re            = regexp.MustCompile(`(.+) \((\d{4})\)`)
	movieIDRegex  = regexp.MustCompile(`/(tt\d{7})/?`)
	movieURLRegex = regexp.MustCompile(`^/title/tt\d{7}/\?`)
)

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

func ExtractMovieURLs(url string) ([]string, error) {
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

	urls := make([]string, len(u.URLs))

	for i, url := range u.URLs {
		urls[i] = url.Location
	}

	return urls, nil
}

func ExtractSitemapURLs() ([]string, error) {
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

	urls := []string{}

	for _, sitemap := range s.Sitemaps {
		if strings.HasPrefix(sitemap.Location, "http://www.imdb.com/sitemap/title") {
			urls = append(urls, sitemap.Location)
		}
	}

	return urls, nil
}

func FetchHTMLDocument(url string) (*goquery.Document, error) {
	resp, err := GetHTTPResponse(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)

	if err != nil {
		return nil, err
	}

	return doc, nil
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

func (s *Scrapher) IsNotEmpty() bool {
	return len(s.movieIDs) != 0
}

func New(url string) *Scrapher {
	movieID := movieIDRegex.FindStringSubmatch(url)[1]

	return &Scrapher{
		movieIDs: map[string]bool{movieID: true},
		visited:  map[string]bool{},
	}
}

func (s *Scrapher) ProcessURL() (*models.Movie, error) {
	var movieID string

	for id, _ := range s.movieIDs {
		movieID = id
		break
	}

	delete(s.movieIDs, movieID)

	s.visited[movieID] = true
	url := "http://www.imdb.com/title/" + movieID
	doc, err := FetchHTMLDocument(url)
	time.Sleep(time.Second * 5)

	if err != nil {
		return nil, err
	}

	movie, err := ExtractMovieInfo(doc)

	if err != nil {
		return nil, err
	}

	nodes := doc.Find("a").Nodes

	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				if movieURLRegex.MatchString(attr.Val) {
					movieID := movieIDRegex.FindStringSubmatch(attr.Val)[1]

					if _, ok := s.visited[movieID]; !ok {
						s.movieIDs[movieID] = true
					}
				}
				break
			}
		}
	}

	return movie, nil
}

func StartFromSitemap() {
	sitemapURLs, err := ExtractSitemapURLs()

	if err != nil {
		log.Fatal(err)
	}

	for _, sitemapURL := range sitemapURLs {
		time.Sleep(time.Second * 5)
		movieURLs, err := ExtractMovieURLs(sitemapURL)

		if err != nil {
			log.Fatal(err)
		}

		for _, movieURL := range movieURLs {
			time.Sleep(time.Second * 5)
			resp, err := GetHTTPResponse(movieURL)

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
