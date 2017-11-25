package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/richardpanda/scrapher"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	movieTitleAndYearRegex = regexp.MustCompile(`(.+) \((\d{4})\)`)
	movieIDRegex           = regexp.MustCompile(`/(tt\d{7})/?`)
	movieURLRegex          = regexp.MustCompile(`^/title/tt\d{7}/\?`)
)

type movie struct {
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

func appendMovieIDs(queue *[]string, visited map[string]bool, in <-chan *goquery.Document) {
	for doc := range in {
		nodes := doc.Find("a").Nodes
		for _, node := range nodes {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					if movieURLRegex.MatchString(attr.Val) {
						movieID := movieIDRegex.FindStringSubmatch(attr.Val)[1]
						if _, ok := visited[movieID]; !ok {
							*queue = append(*queue, movieID)
							visited[movieID] = true
						}
					}
					break
				}
			}
		}
	}
}

func extractMovieInfo(in <-chan *goquery.Document, out chan<- *movie, e chan<- error) {
	for doc := range in {
		id, ok := doc.Find("meta[property=\"pageId\"]").First().Attr("content")
		if !ok {
			e <- errors.New("cannot find movie id from imdb")
			continue
		}

		url := "http://www.imdb.com/title/" + id
		str := doc.Find("[itemprop=\"name\"]").First().Text()

		matches := movieTitleAndYearRegex.FindStringSubmatch(strings.TrimSpace(str))
		if len(matches) < 3 {
			msg := fmt.Sprintf("unable to parse title and year from imdb (%s)", url)
			e <- errors.New(msg)
			continue
		}

		title := matches[1]
		year, err := strconv.Atoi(matches[2])
		if err != nil {
			e <- err
			continue
		}

		rating, err := strconv.ParseFloat(strings.Split(doc.Find("[itemprop=\"ratingValue\"]").First().Text(), "/")[0], 64)
		if err != nil {
			msg := fmt.Sprintf("unable to parse rating from imdb (%s)", url)
			e <- errors.New(msg)
			continue
		}

		numRatings, err := stringToInt(doc.Find("[itemprop=\"ratingCount\"]").First().Text())
		if err != nil {
			msg := fmt.Sprintf("unable to parse number of ratings from imdb (%s)", url)
			e <- errors.New(msg)
			continue
		}

		out <- &movie{
			IMDBNumRatings: numRatings,
			IMDBRating:     rating,
			IMDBURL:        url,
			Title:          title,
			Year:           year,
		}
	}
}

func fanOut(in <-chan *goquery.Document, out1, out2 chan<- *goquery.Document) {
	for doc := range in {
		out1 <- doc
		out2 <- doc
	}
}

func fetchHTMLDocument(in <-chan string, out chan<- *goquery.Document, e chan<- error) {
	for url := range in {
		resp, err := getHTTPResponse(url)
		if err != nil {
			e <- err
			continue
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			e <- err
			continue
		}

		out <- doc
	}
}

func getHTTPResponse(url string) (*http.Response, error) {
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

func insertMovie(db *gorm.DB, in <-chan *movie) {
	for movie := range in {
		db.Create(movie)
		fmt.Printf("added %s (%d)\n", movie.Title, movie.Year)
	}
}

func printError(e <-chan error) {
	for err := range e {
		fmt.Println(err)
	}
}

func stringToInt(s string) (int, error) {
	return strconv.Atoi(strings.Replace(s, ",", "", -1))
}

func main() {
	const sleepDuration = time.Duration(5) * time.Second
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.AutoMigrate(&scrapher.Movie{})

	url := "http://www.imdb.com/title/tt0468569"
	movieID := movieIDRegex.FindStringSubmatch(url)[1]
	queue := []string{movieID}
	visited := map[string]bool{movieID: true}

	t := time.NewTicker(sleepDuration)
	fetcher := make(chan string)
	spreader := make(chan *goquery.Document)
	appender := make(chan *goquery.Document)
	extractor := make(chan *goquery.Document)
	database := make(chan *movie)
	e := make(chan error)

	go fetchHTMLDocument(fetcher, spreader, e)
	go fanOut(spreader, appender, extractor)
	go appendMovieIDs(&queue, visited, appender)
	go extractMovieInfo(extractor, database, e)
	go insertMovie(db, database)
	go printError(e)

	for _ = range t.C {
		if len(queue) == 0 {
			return
		}

		movieID := queue[0]
		queue = queue[1:]
		url := "http://www.imdb.com/title/" + movieID
		fetcher <- url
	}
}
