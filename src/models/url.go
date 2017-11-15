package models

import "encoding/xml"

type URL struct {
	Location string `xml:"loc"`
}

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}
