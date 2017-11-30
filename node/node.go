package node

import (
	"github.com/PuerkitoBio/goquery"
)

type DocNode struct {
	*goquery.Document
	*QueueNode
}

type QueueNode struct {
	Depth   int
	MovieID string
}
