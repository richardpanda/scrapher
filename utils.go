package scrapher

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchHTMLDocument(url string) (*goquery.Document, error) {
	resp, err := getHTTPResponse(url)

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

func StringToInt(s string) (int, error) {
	return strconv.Atoi(strings.Join(strings.Split(s, ","), ""))
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
