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
	Valutes []Valute `xml:"Valute"`
}

// Valute holds the fields we use from each <Valute> in the CBR feed.
type Valute struct {
	CharCode  string `xml:"CharCode"`
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
	defer func() { _ = res.Body.Close() }()
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
	currency := pflag.StringP("currency", "c", "USD", "currency code to convert, e.g. USD or EUR")
	amountStr := pflag.StringP("amount", "a", "1", "amount of currency to convert")
	date := pflag.StringP("date", "d", "", "rate date in DD.MM.YYYY format (default: latest)")
	raw := pflag.BoolP("raw", "r", false, "print only the numeric result, no decoration")
	invert := pflag.BoolP("invert", "i", false, "show the inverse rate (1 RUB = X CUR)")
	pflag.Parse()

	*currency = strings.ToUpper(*currency)

	// Parse the amount ourselves (instead of via a pflag float flag) so a bad
	// value produces a friendly message instead of pflag's strconv error.
	amount, err := strconv.ParseFloat(*amountStr, 64)
	if err != nil || amount <= 0 {
		fmt.Fprintln(os.Stderr, "Amount must be a positive number")
		os.Exit(1)
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
		if v.CharCode != *currency {
			continue
		}
		// VunitRate is the rate for a single unit (CBR's own Value/Nominal), which
		// is what we convert against.
		rate, err := parseRate(v.VunitRate)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not decode rate.")
			os.Exit(1)
		}

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
			fmt.Printf("%g RUB = %.*f %s (%s)\n", amount, precision, result, *currency, curs.Date)
		} else {
			fmt.Printf("%g %s = %.*f RUB (%s)\n", amount, *currency, precision, result, curs.Date)
		}
		return
	}

	fmt.Fprintf(os.Stderr, "Unknown currency: %s\n", *currency)
	os.Exit(1)
}
