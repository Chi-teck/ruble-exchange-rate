// Command rer fetches RUB exchange rates from the Central Bank of the Russian
// Federation and converts amounts between RUB and other currencies.
package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"golang.org/x/net/html/charset"
)

const cbrURL = "https://www.cbr.ru/scripts/XML_daily.asp"
const inputDateFormat = "02.01.2006"
const cbrDateFormat = "02/01/2006"

// version is overridden at release time via -ldflags "-X main.version=...".
// When left as "dev", buildVersion falls back to Go's build info.
var version = "dev"

// errNoRates is returned when the feed has no currencies for the requested
// date (CBR's data starts 1 July 1992).
var errNoRates = errors.New("no rates available for this date")

type valCurs struct {
	Date    string   `xml:"Date,attr"`
	Valutes []valute `xml:"Valute"`
}

// valute holds the fields we use from each <Valute> in the CBR feed.
type valute struct {
	CharCode  string `xml:"CharCode"`
	VunitRate string `xml:"VunitRate"`
}

// parseRate converts a CBR-formatted number into a float64.
func parseRate(s string) (float64, error) {
	return strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
}

// fetchRates retrieves and decodes the CBR daily rates document at url.
func fetchRates(ctx context.Context, url string) (*valCurs, error) {
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

	var curs valCurs
	decoder := xml.NewDecoder(res.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&curs); err != nil {
		return nil, fmt.Errorf("failed to decode CBR response: %w", err)
	}
	if len(curs.Valutes) == 0 {
		return nil, errNoRates
	}
	return &curs, nil
}

// formatResult renders a converted value: 2-decimal cents for normal
// magnitudes, ~4 significant figures (trailing zeros trimmed) for small
// results so inverse rates don't collapse to "0.00". Assumes result > 0.
func formatResult(result float64) string {
	if result >= 1 {
		return fmt.Sprintf("%.2f", result)
	}
	precision := 3 - int(math.Floor(math.Log10(result)))
	if precision > 6 {
		precision = 6
	}
	pow := math.Pow(10, float64(precision))
	return strconv.FormatFloat(math.Round(result*pow)/pow, 'f', -1, 64)
}

// buildVersion returns the version string to display. A release build sets
// version via -ldflags. Otherwise, we fall back to Go's build info: the module
// version for `go install ...@vX`, or "dev+<commit>[.dirty]" for a local build
// from a Git checkout.
func buildVersion() string {
	if version != "dev" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return version
	}
	if v := info.Main.Version; v != "" && v != "(devel)" {
		return v
	}
	var rev, dirty string
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			if len(s.Value) >= 7 {
				rev = s.Value[:7]
			}
		case "vcs.modified":
			if s.Value == "true" {
				dirty = ".dirty"
			}
		}
	}
	if rev != "" {
		return version + "+" + rev + dirty
	}
	return version
}

func main() {
	currency := pflag.StringP("currency", "c", "USD", "currency code to convert, e.g. USD or EUR")
	amountStr := pflag.StringP("amount", "a", "1", "amount of currency to convert")
	date := pflag.StringP("date", "d", "", "rate date in DD.MM.YYYY format (default: latest)")
	raw := pflag.BoolP("raw", "r", false, "print only the numeric result, no decoration")
	invert := pflag.BoolP("invert", "i", false, "show the inverse rate (1 RUB = X CUR)")
	showVersion := pflag.Bool("version", false, "print version and exit")
	pflag.Parse()

	if *showVersion {
		fmt.Println("rer " + buildVersion())
		return
	}

	*currency = strings.ToUpper(*currency)

	// Parse the amount ourselves (instead of via a pflag float flag) so a bad
	// value produces a friendly message instead of pflag's strconv error.
	amount, err := strconv.ParseFloat(*amountStr, 64)
	if err != nil || amount <= 0 {
		fmt.Fprintln(os.Stderr, "Amount must be a positive number.")
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
		if errors.Is(err, errNoRates) {
			fmt.Fprintln(os.Stderr, "No exchange rate data available for that date.")
		} else {
			fmt.Fprintln(os.Stderr, "Could not fetch exchange rates from CBR.")
		}
		os.Exit(1)
	}

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
		// Format the amount as plain decimal too, so large values don't switch to
		// scientific notation (as %g would, e.g. 1e+06).
		amountFmt := strconv.FormatFloat(amount, 'f', -1, 64)
		resultFmt := formatResult(result)
		if *raw {
			fmt.Printf("%.4f\n", result)
		} else if *invert {
			fmt.Printf("%s RUB = %s %s (%s)\n", amountFmt, resultFmt, *currency, curs.Date)
		} else {
			fmt.Printf("%s %s = %s RUB (%s)\n", amountFmt, *currency, resultFmt, curs.Date)
		}
		return
	}

	codes := make([]string, len(curs.Valutes))
	for i, v := range curs.Valutes {
		codes[i] = v.CharCode
	}
	fmt.Fprintf(os.Stderr, "Unknown currency: %s.\nAvailable: %s\n", *currency, strings.Join(codes, " "))
	os.Exit(1)
}
