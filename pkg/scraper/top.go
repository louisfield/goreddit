package scraper

import (
	"strings"
	"regexp"
)


func GetTopUrls(url string) [5]string {
	body, _ := get(url)
	r, _ := regexp.Compile(`(http|ftp|https):\/\/([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:\/~+#-]*[\w@?^=%&\/~+#-])`)
	match := r.FindAllString(body, -1)
	counter := 0
	reddit_links := [5]string{}
	for i := 0; (i < len(match) -1) ; {
		if strings.Contains(match[i], "comments") {
			reddit_links[counter] = match[i]
			counter++
			if counter == 5 {
				break
			}
		}
		i++
	}
	if(reddit_links == [5]string{}) {
		reddit_links = GetTopUrls(url)
	}
	return reddit_links
}

func get(url string) (string, error) {
	resp := ReadBody(url)

	return ParseHTML(string(resp[:]))
}