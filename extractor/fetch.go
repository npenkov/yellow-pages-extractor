package extractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/npenkov/yellow-pages-extractor/model"
	"golang.org/x/net/html"
)

const (
	DefaultYPURL = "https://www.yellowpages.com/search"
)

// Extractor defines the protocol for extracting information from yellowpages.
// When record is fetched, the supplied addFunc is called.
type Extractor interface {
	FetchNextPage(addFunc func([]model.Contact)) (bool, error)
}

type YellowPages struct {
	nextPageURL string
}

// NewYPExtractor creates new instance of page extractor, by supplying search terms and geo location terms
func NewYPExtractor(searchTerms, geoLocationTerms string) *YellowPages {
	stq := url.QueryEscape(searchTerms)
	glq := url.QueryEscape(geoLocationTerms)
	yp := &YellowPages{
		nextPageURL: fmt.Sprintf("%s?search_terms=%s&geo_location_terms=%s&page=1", DefaultYPURL, stq, glq),
	}
	return yp
}

func (this *YellowPages) FetchNextPage(addFunc func([]model.Contact)) (bool, error) {
	nextPage := this.nextPageURL
	resp, err := http.Get(nextPage)
	if err != nil {
		return false, fmt.Errorf("Error fetching from url %s : %v", nextPage, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("Error reading response from %s : %v", nextPage, err)
	}
	jsonData, np, err := extractLocationsFromHTML(bytes.NewReader(body))
	if err != nil {
		return false, fmt.Errorf("Cannot parse html : %v", err)
	}
	var contacts []model.Contact
	if err := json.Unmarshal([]byte(jsonData), &contacts); err != nil {
		return false, fmt.Errorf("Cannot parse json %s : %v", jsonData, err)
	}
	if addFunc != nil {
		addFunc(contacts)
	}
	this.nextPageURL = np

	return this.nextPageURL != "", nil
}

func extractLocationsFromHTML(r io.Reader) (string, string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", "", err
	}
	var jsonData string
	var nextPage string

	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			nextActive := false
			for _, a := range n.Attr {
				if a.Key == "rel" && a.Val == "next" {
					nextActive = true
					break
				}
			}
			if nextActive {
				for _, a := range n.Attr {
					if a.Key == "href" {
						nextPage = a.Val
						break
					}
				}
			}
		}
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, a := range n.Attr {
				if a.Key == "type" && a.Val == "application/ld+json" && strings.HasPrefix(n.FirstChild.Data, "[") {
					jsonData = n.FirstChild.Data
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return jsonData, nextPage, nil
}
