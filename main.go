// Prototype to download http NWS ZFP weather product, and parse each
// forecast area into a map indexed by the ZFP ZoneID
// Example usage
//   ./nws-http sew [ lot |  aer  | alu ]
//
// Original http code was borrowe from Go Programming Language net/http

package main

import (
	"fmt"
	"nws-http/nwszfp"
	"os"
	"strings"
)

// This "database" for various regions needs improvement
const (
	NWS_AK           = "AER"
	NWS_ALUETIANS_AK = "ALU"
	NWS_WA           = "KSEW"
	NWS_IL           = "KLOT"

	Anchorage_AK        = "AKZ101"
	BristonBay_AK       = "AKZ161"
	EasternAluetians_AK = "AKZ185"

	Chicago_IL = "ILZ014"
	Wheaton_IL = "ILZ013"

	Bellevue_WA  = "WAZ556"
	Bremerton_WA = "WAZ559"
	Everett_WA   = "WAZ507"
	Seattle_WA   = "WAZ558"
)

func printZone(s []string) {
	doIndent := false
	nl := "- "
	for _, l := range s {
		if strings.HasPrefix(l, ".") {
			l = strings.Replace(l, ".", nl, 1)
			nl = "\n- "
			doIndent = true
		} else if doIndent {
			fmt.Printf("  ")
		}
		fmt.Printf("%s", l)
		if strings.HasSuffix(l, "$$") == true {
			fmt.Println()
		}
	}
}

func main() {
	for _, url := range os.Args[1:] {
		m := nwszfp.Get(url)

		// Filter returned Zone Forecast by geographic region
		// This really needs to be improved.
		// printZone(m[Seattle_WA])

		printZone(m[Bellevue_WA])

		// printZone(m[Everett_WA])
		// printZone(m[NWS_WA])
		// printZone(m[Seattle_WA])
		// printZone(m[Bellevue_WA])
		// printZone(m[Everett_WA])
		// printZone(m[Bremerton_WA])

		printZone(m[NWS_IL])
		printZone(m[Chicago_IL])
		printZone(m[Wheaton_IL])

		printZone(m[NWS_AK])
		printZone(m[Anchorage_AK])

		printZone(m[NWS_ALUETIANS_AK])
		printZone(m[EasternAluetians_AK])
	}
}
