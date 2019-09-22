# YellowPages extractor

The utility searches in YellowPages.com by geo terms and search terms and returns the entities together with address and geo coding information.
For details on the output format see [YellowPages return model](model/ypages.go)

## Using in your program

Fetch the package
```
go get github.com/npenkov/yellow-pages-extractor
```

### Example

```go
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

```

See [example.go](example/example.go)

Run it by calling

```
go run example/example.go
```

## Using as command line utility

```
go get github.com/npenkov/yellow-pages-extractor/yp-extractor
yp-extractor -help
```