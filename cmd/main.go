package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// client := &http.Client{}
	// req, err := http.NewRequest("GET", "http://www.imdb.com/title/tt4574334/?ref_=nv_sr_1", nil)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// req.Header.Add("User-Agent", "Scrapher, a friendly web scraper. Code can be found at https://github.com/richardpanda/scrapher.")

	// resp, err := client.Do(req)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer resp.Body.Close()

	// bodyBytes, err := ioutil.ReadAll(resp.Body)

	re := regexp.MustCompile(`(.+) \((\d{4})\)`)
	f, err := os.Open("index.html")

	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)

	if err != nil {
		log.Fatal(err)
	}

	matches := re.FindStringSubmatch(strings.TrimSpace(doc.Find("[itemprop=\"name\"]").First().Text()))
	title := matches[1]
	year := matches[2]
	rating := strings.Split(doc.Find("[itemprop=\"ratingValue\"]").First().Text(), "/")[0]
	numRatings := doc.Find("[itemprop=\"ratingCount\"]").First().Text()

	fmt.Println(title)
	fmt.Println(year)
	fmt.Println(rating)
	fmt.Println(numRatings)
}
