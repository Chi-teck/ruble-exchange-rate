// Command rer fetches RUB exchange rates from the Central Bank of the Russian
// Federation and converts amounts between RUB and other currencies.
package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"golang.org/x/net/html/charset"
)

const cbrURL = "https://www.cbr.ru/scripts/XML_daily.asp"
const inputDateFormat = "02.01.2006"
const cbrDateFormat = "02/01/2006"

type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID        string `xml:"ID,attr"`
	NumCode   string `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Nominal   int    `xml:"Nominal"`
	Name      string `xml:"Name"`
	Value     string `xml:"Value"`
	VunitRate string `xml:"VunitRate"`
}

// parseRate converts a CBR-formatted number into a float64.
func parseRate(s string) (float64, error) {
	return strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
}

// fetchRates retrieves and decodes the CBR daily rates document at url.
func fetchRates(ctx context.Context, url string) (*ValCurs, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	// CBR returns 403 for Go's default User-Agent ("Go-http-client/..."),
	// so set an explicit one.
	req.Header.Add("User-Agent", "RER")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to CBR failed: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %q of CBR response", res.Status)
	}

	var curs ValCurs
	decoder := xml.NewDecoder(res.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&curs); err != nil {
		return nil, fmt.Errorf("failed to decode CBR response: %w", err)
	}
	return &curs, nil
}

//goland:noinspection GoUnhandledErrorResult
func main() {
	date := pflag.String("date", "", "rate date in DD.MM.YYYY format (default: latest)")
	raw := pflag.Bool("raw", false, "print only the numeric result, no decoration")
	invert := pflag.Bool("invert", false, "show the inverse rate (1 RUB = X CUR)")
	pflag.Parse()

	args := pflag.Args()
	amount := 1.0
	currency := "USD"
	if len(args) > 0 {
		a, err := strconv.ParseFloat(args[0], 64)
		if err != nil || a <= 0 {
			fmt.Fprintln(os.Stderr, "Amount must be a positive number")
			os.Exit(1)
		}
		amount = a
	}
	if len(args) > 1 {
		currency = strings.ToUpper(args[1])
	}

	url := cbrURL
	if *date != "" {
		d, err := time.Parse(inputDateFormat, *date)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Date must be in DD.MM.YYYY format.")
			os.Exit(1)
		}
		url += "?date_req=" + d.Format(cbrDateFormat)
	}

	curs, err := fetchRates(context.Background(), url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Currency given: convert amount into RUB.
	for _, v := range curs.Valutes {
		if v.CharCode != currency {
			continue
		}
		// The per-unit rate is Value/Nominal. Older documents (e.g. historical --date queries)
		// may omit VunitRate, but Value and Nominal are always present.
		value, err := parseRate(v.Value)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not decode rate.")
			os.Exit(1)
		}
		if v.Nominal <= 0 {
			fmt.Fprintln(os.Stderr, "Invalid nominal in CBR response.")
			os.Exit(1)
		}
		rate := value / float64(v.Nominal)

		var result float64
		if *invert {
			result = amount / rate
		} else {
			result = amount * rate
		}
		// Show more decimals for small results (e.g. inverse rates) so they
		// don't round down to a single significant digit.
		precision := 2
		if result < 0.5 {
			precision = 3
		}
		if *raw {
			fmt.Printf("%.4f\n", result)
		} else if *invert {
			fmt.Printf("%g RUB = %.*f %s (%s)\n", amount, precision, result, currency, curs.Date)
		} else {
			fmt.Printf("%g %s = %.*f RUB (%s)\n", amount, currency, precision, result, curs.Date)
		}
		return
	}

	fmt.Fprintf(os.Stderr, "Unknown currency: %s\n", currency)
	os.Exit(1)
}
