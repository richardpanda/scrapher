package main

import "github.com/richardpanda/scrapher/queue"

type App struct {
	imdb *queue.Queue
	rt   *queue.Queue
}
