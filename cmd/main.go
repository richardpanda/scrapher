package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func extractTitle(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "class" && attr.Val == "title_wrapper" {
			return strings.TrimSpace(n.FirstChild.NextSibling.FirstChild.Data)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if title := extractTitle(c); title != "" {
			return title
		}
	}

	return ""
}

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://www.imdb.com/title/tt4574334/?ref_=nv_sr_1", nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("User-Agent", "Scrapher, a friendly web scraper. Code can be found at https://github.com/richardpanda/scrapher.")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	// bodyBytes, err := ioutil.ReadFile("index.html")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	bodyString := string(bodyBytes)
	doc, err := html.Parse(strings.NewReader(bodyString))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(extractTitle(doc))
}