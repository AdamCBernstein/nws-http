package nwszfp

import (
	"fmt"
	"log"
	"strings"
)

const (
	//  issuedby=%s is filled in to narrow which region is retrieved.
	urlTemplate = "https://forecast.weather.gov/product.php?site=AKQ&issuedby=%s&product=ZFP&format=txt&version=1&glossary=0"
)

func Get(zoneId string) map[string][]string {
	var tag string
	var eof bool

	m := make(map[string][]string)

	url := fmt.Sprintf(urlTemplate, zoneId)
	hs, err := newHttpScanner(url)
	if err != nil {
		log.Fatal(err)
	}

	s := skipLines(hs, 3)
	if len(s) == 0 {
		log.Fatal("Missing expected header")
	}

	s, tag = getZoneHeader(hs)
	m[tag] = s

	// Find all "Zone Forecast Product" sections
	s, eof = returnLines(hs, "$$")
	for eof == false {
		tag := strings.Split(s[0], "-")
		m[tag[0]] = s
		s, eof = returnLines(hs, "$$")
	}
	closeHttpScanner(hs)
	return m
}
