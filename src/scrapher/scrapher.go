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
var movieIDRegex = regexp.MustCompile(`/(tt\d{7})/`)
var movieURLRegex = regexp.MustCompile(`^/title/tt\d{7}/\?`)

func ExtractMovieInfo(url string, doc *goquery.Document) (*models.Movie, error) {
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
		URL:        url,
		NumRatings: numRatings,
		Rating:     rating,
		Year:       year,
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

			movie, err := ExtractMovieInfo(movieLink, doc)

			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(movie)
			resp.Body.Close()
		}
	}
}

func StartFromURL(url string) {
	movieIDs := []string{movieIDRegex.FindStringSubmatch(url)[1]}
	visited := map[string]bool{}

	for len(movieIDs) > 0 {
		movieID := movieIDs[0]
		movieIDs = movieIDs[1:]

		if _, ok := visited[movieID]; ok {
			continue
		}

		visited[movieID] = true
		url := "http://www.imdb.com/title/" + movieID
		resp, err := GetHTTPResponse(url)
		time.Sleep(time.Second * 5)

		if err != nil {
			log.Fatal(err)
		}

		doc, err := goquery.NewDocumentFromResponse(resp)

		if err != nil {
			log.Fatal(err)
		}

		movie, err := ExtractMovieInfo(url, doc)

		if err != nil {
			fmt.Printf("%s\t%s\n", url, err)
			continue
		}

		fmt.Printf("%s (%d)\n", movie.Title, movie.Year)
		nodes := doc.Find("a").Nodes

		for _, node := range nodes {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					if movieURLRegex.MatchString(attr.Val) {
						movieID := movieIDRegex.FindStringSubmatch(attr.Val)[1]

						if _, ok := visited[movieID]; !ok {
							movieIDs = append(movieIDs, movieID)
						}
					}
					break
				}
			}
		}

		resp.Body.Close()
	}
}
