package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/npenkov/yellow-pages-extractor/extractor"
	"github.com/npenkov/yellow-pages-extractor/model"
)

func main() {

	var searchTerms, geoTerms string
	flag.StringVar(&searchTerms, "s", "", "Space separated search terms (e.g. \"grocery stores\")")
	flag.StringVar(&geoTerms, "g", "", "Space separated geo terms (e.g. \"atlanta ga\")")

	flag.Parse()

	if searchTerms == "" {
		log.Fatalf("Unsufficient parameters - no search terms specified")
	}
	if geoTerms == "" {
		log.Fatalf("Unsufficient parameters - no geo terms specified")
	}

	e := extractor.NewYPExtractor(searchTerms, geoTerms)

	addFunc := func(stores []model.Contact) {
		for _, s := range stores {
			var js []byte
			var err error
			if js, err = json.MarshalIndent(s, "", "  "); err != nil {
				log.Fatalf("Error marshalig object to json : %v", err)
			}
			fmt.Println(string(js))
		}
	}
	hasMore, err := e.FetchNextPage(addFunc)
	if err != nil {
		log.Fatalf("Error fetching records from YellowPages : %v", err)
	}
	for hasMore && err == nil {
		hasMore, err = e.FetchNextPage(addFunc)
		if err != nil {
			log.Fatalf("Error fetching records from YellowPages : %v", err)
		}
	}
}
