package htmldoc

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func ExtractURLs(doc *goquery.Document) []string {
	out := []string{}
	nodes := doc.Find("a").Nodes
	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key != "href" {
				continue
			}
			out = append(out, attr.Val)
		}
	}
	return out
}

func Get(url string) (*goquery.Document, error) {
	resp, err := getHTTPResponse(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return doc, err
}

// func Get(in <-chan string, out chan<- *goquery.Document, e chan<- error) {
// 	for url := range in {
// resp, err := getHTTPResponse(url)
// if err != nil {
// 	e <- err
// 	continue
// }
// defer resp.Body.Close()

// doc, err := goquery.NewDocumentFromResponse(resp)
// if err != nil {
// 	e <- err
// 	continue
// }

// out <- doc
// }
// }

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
