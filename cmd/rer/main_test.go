package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// sampleXML is a trimmed XML_daily.asp document. Only ASCII is used so it can be
// served as-is under the declared windows-1251 charset (1251 is ASCII-compatible
// for these bytes), keeping the test free of an encoding dependency.
const sampleXML = `<?xml version="1.0" encoding="windows-1251"?>` +
	`<ValCurs Date="13.05.2025" name="Foreign Currency Market">` +
	`<Valute ID="R01235"><NumCode>840</NumCode><CharCode>USD</CharCode><Nominal>1</Nominal><Name>US Dollar</Name><Value>80,8883</Value><VunitRate>80,8883</VunitRate></Valute>` +
	`<Valute ID="R01060"><NumCode>051</NumCode><CharCode>AMD</CharCode><Nominal>100</Nominal><Name>Armenian Dram</Name><Value>20,8014</Value><VunitRate>0,208014</VunitRate></Valute>` +
	`</ValCurs>`

func TestFetchRates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/xml; charset=windows-1251")
		_, _ = w.Write([]byte(sampleXML))
	}))
	defer srv.Close()

	curs, err := fetchRates(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("fetchRates: %v", err)
	}
	if curs.Date != "13.05.2025" {
		t.Errorf("Date = %q, want 13.05.2025", curs.Date)
	}
	if len(curs.Valutes) != 2 {
		t.Fatalf("got %d valutes, want 2", len(curs.Valutes))
	}

	// AMD has a nominal of 100; we convert against VunitRate, the per-unit rate.
	var amd valute
	for _, v := range curs.Valutes {
		if v.CharCode == "AMD" {
			amd = v
		}
	}
	rate, err := parseRate(amd.VunitRate)
	if err != nil {
		t.Fatalf("parseRate: %v", err)
	}
	if want := 0.208014; rate != want {
		t.Errorf("AMD VunitRate = %v, want %v", rate, want)
	}
}

// TestFetchRatesNoData checks that an empty document — CBR's response for dates
// before its data begins (1 July 1992) — yields errNoRates, not an empty result.
func TestFetchRatesNoData(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/xml; charset=windows-1251")
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="windows-1251"?><ValCurs></ValCurs>`))
	}))
	defer srv.Close()

	if _, err := fetchRates(context.Background(), srv.URL); !errors.Is(err, errNoRates) {
		t.Errorf("got %v, want errNoRates", err)
	}
}

// TestBuildVersion covers the release path: when version is set via -ldflags
// (i.e. not "dev"), buildVersion returns it verbatim. The build-info fallback
// is environment-dependent and left unasserted.
func TestBuildVersion(t *testing.T) {
	saved := version
	defer func() { version = saved }()

	version = "1.2.3"
	if got := buildVersion(); got != "1.2.3" {
		t.Errorf("buildVersion() = %q, want %q", got, "1.2.3")
	}
}

func TestFormatResult(t *testing.T) {
	for _, tc := range []struct {
		in   float64
		want string
	}{
		{8503.05, "8503.05"}, // normal magnitude: 2-decimal cents
		{366, "366.00"},      // whole result keeps cents
		{0.1995, "0.1995"},   // small: ~4 significant figures
		{0.0136, "0.0136"},   // inverse rate
		{0.0004, "0.0004"},   // would have collapsed to 0.000 before
		{0.00009, "0.00009"}, // plain decimal, never scientific
	} {
		if got := formatResult(tc.in); got != tc.want {
			t.Errorf("formatResult(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
