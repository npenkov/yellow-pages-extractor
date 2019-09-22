package main

import (
	"fmt"
	"log"

	"github.com/npenkov/yellow-pages-extractor/extractor"
	"github.com/npenkov/yellow-pages-extractor/model"
)

func main() {
	e := extractor.NewYPExtractor("grocery stores", "atlanta ga")

	addFunc := func(stores []model.Contact) {
		for _, s := range stores {
			fmt.Printf("Store: %s Address: %s Coordinates: %f, %f\n", s.Name, s.Address.Address, s.Geo.Latitude, s.Geo.Longitude)
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
