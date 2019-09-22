package model

type Address struct {
	Country    string `json:"addressCountry"`
	Address    string `json:"streetAddress"`
	City       string `json:"addressLocality"`
	State      string `json:"addressRegion"`
	PostalCode string `json:"postalCode"`
}

type Geo struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type Contact struct {
	Name      string  `json:"name"`
	Address   Address `json:"address"`
	Geo       Geo     `json:"geo"`
	Telephone string  `json:"telephone"`
	Image     string  `json:"image"`
}