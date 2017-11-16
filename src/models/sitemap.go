package models

import (
	"encoding/xml"
	"time"
)

type Sitemap struct {
	Location     string    `xml:"loc"`
	LastModified time.Time `xml:"lastmod"`
}

type SitemapIndex struct {
	XMLName  xml.Name  `xml:"sitemapindex"`
	Sitemaps []Sitemap `xml:"sitemap"`
}
