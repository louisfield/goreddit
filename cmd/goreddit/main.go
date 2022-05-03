package main

import (
	"github.com/louisfield/goreddit/pkg/scraper"
	"log"
)

type redditPage struct {
	url string
	comment string
	arrayIndex int
}

func main() {
	reddit_links := scraper.GetTopUrls("https://www.reddit.com/r/leagueoflegends/top/")
	finalComments := scraper.GetComments(reddit_links)
	log.Println(finalComments)
}
