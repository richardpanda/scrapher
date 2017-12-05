# Scrapher

[![Build Status](https://travis-ci.org/richardpanda/scrapher.svg?branch=master)](https://travis-ci.org/richardpanda/scrapher)

A web scraper for IMDB and Rotten Tomatoes.

## Quick Start

Download and install dependencies
```bash
$ go get ./...
```

Set DB_NAME and DB_USER environment variables for Postgres connection
```bash
export DB_NAME="scrapher_dev"
export DB_USER="user"
```

Execute main.go to start scraping!
```bash
go run cmd/main.go -imdb http://www.imdb.com/title/tt0468569 -rt https://www.rottentomatoes.com/m/the_dark_knight
```

## Usage
### -imdb *url*
IMDB start URL
```bash
go run cmd/main.go -imdb http://www.imdb.com/title/tt0468569
```

### -rt *url*
Rotten Tomatoes start URL
```bash
go run cmd/main.go -rt https://www.rottentomatoes.com/m/the_dark_knight
```

### -d *depth* (optional)
How deep the web scraper should crawl
```bash
# The scraper will visit the initial URL and all the URLs inside the initial URL.
go run cmd/main.go -imdb http://www.imdb.com/title/tt0468569 -d 1

# If depth is not specified, the scraper will keep scraping until there is no more URLs to visit.
go run cmd/main.go -imdb http://www.imdb.com/title/tt0468569
```
