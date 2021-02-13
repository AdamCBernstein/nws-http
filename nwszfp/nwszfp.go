// Prototype to download http NWS ZFP weather product, and parse each
// forecast area into a map indexed by the ZFP ZoneID
// Example useage
// ./ksew-http sew
// 
// This will fill in the below URL template "issuedby=%s"
//    https://forecast.weather.gov/product.php?site=AKQ&issuedby=%s&product=ZFP&format=txt&version=1&glossary=0'

package nwszfp

import (
	"bufio"
	"net/http"
	"strconv"
	"strings"
)

type httpScanner struct {
	scanner  *bufio.Scanner
	httpResp *http.Response
	prevLine string
}

func hasNumPrefix(s string, i int) bool {
	if i < 1 || len(s) < i {
		return false
	}
	_, err := strconv.Atoi(s)
	return err == nil
}

func skipLines(hs *httpScanner, l int) []string {
	for hs.scanner.Scan() {
		t := hs.scanner.Text()
		if hasNumPrefix(t, l) == true {
			r := strings.SplitN(t, "", 1)
			return r
		}
	}
	return strings.SplitN("", "", 1)
}

func getZoneHeader(hs *httpScanner) ([]string, string) {
	var lines []string
	var tags []string

	for hs.scanner.Scan() {
		l := hs.scanner.Text()
		if strings.HasSuffix(l, "-") {
			hs.prevLine = l
			break
		}
		lines = append(lines, l)
		lines = append(lines, "\n")
	}
	tags = strings.Split(lines[0], " ")

	// Get rid of blank lines
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] == "" {
			lines = lines[0:i]
		} else {
			break
		}
	}
	return lines, tags[1]
}

func eatBlankLines(hs *httpScanner) {
	var l string
	for hs.scanner.Scan() {
		l = hs.scanner.Text()
		if len(l) == 0 {
			continue
		} else if l[0] == '\n' {
			continue
		} else {
			break
		}
	}
	hs.prevLine = l
}

func returnLines(hs *httpScanner, s string) ([]string, bool) {
	var l string
	var lines []string

	// Scan until next 's' token line is found
	if len(hs.prevLine) > 0 {
		lines = append(lines, hs.prevLine)
		hs.prevLine = ""
	}

	for hs.scanner.Scan() {
		l = hs.scanner.Text()
		lines = append(lines, l)
		lines = append(lines, "\n")
		if strings.HasPrefix(l, s) == true {
			// Dispose of empty lines following terminal string
			eatBlankLines(hs)

			// Return accumulated slice of lines
			return lines, false
		}
	}
	return lines, true
}

func newHttpScanner(url string) (hs *httpScanner, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	hs = new(httpScanner)
	hs.httpResp = resp
	hs.scanner = bufio.NewScanner(resp.Body)

	return hs, nil
}

func closeHttpScanner(hs *httpScanner) {
	hs.httpResp.Body.Close()
}
