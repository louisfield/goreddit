package scraper

import (
	"io/ioutil"
	"bytes"
	"net/http"
	"golang.org/x/net/html"
	"io"
	"strings"
	"errors"
	"time"
	"log"
	"github.com/PuerkitoBio/goquery"
)

type redditPage struct {
	url string
	comment string
	arrayIndex int
}

func ParseHTML(html_string string) (string, error) {
	doc, _ := html.Parse(strings.NewReader(html_string))
	bn, err := body(doc)
	if err != nil {
        return "fail", err
    }
	body := renderBody(bn)
	return body, err
}

func renderBody(node *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, node)
	return buf.String()
}
func crawler(node *html.Node) *html.Node {
	var body *html.Node 
	if node.Type == html.ElementNode && node.Data == "body" {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		body = crawler(child)
	}
	return body
}

func body(doc *html.Node) (*html.Node, error) {
    var body *html.Node = crawler(doc)

    if body != nil {
        return body, nil
    }
    return nil, errors.New("Missing <body> in the node tree")
}
func Get(url string) (*http.Response, error) {

	resp, err := http.Get(url)

	test_err(err)

	return resp, err 
}

func ReadBody(url string) []byte {
	resp, err := Get(url)
	html, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	test_err(err)

	return html

}

func GetComments(reddit_links [5]string) []redditPage {
	finalComments := []redditPage{}
	for i := 0; (i < len(reddit_links)); {
		res, err := http.Get(reddit_links[i])
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(8 * time.Second)

		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}
	
		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		didFind := false
		doc.Find("p._1qeIAgB0cPwnLhDF9XSiJM").Each(func(index int, s *goquery.Selection) {
			didFind = true
			text := s.Text() 
			finalComments = append(finalComments, redditPage{reddit_links[i], text, i}) 
		})
		if(didFind) {
			i++ 
		}
	}
	return finalComments
}

func test_err(err error) {
	if err != nil {
		panic(err)
	}
}
