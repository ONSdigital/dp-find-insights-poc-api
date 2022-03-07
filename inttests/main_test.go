//go:build integration
// +build integration

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/pkg/diff"
)

// parse flags
var extra = flag.Bool("extra", false, "print full diffs on failed integration test")
var local = flag.Bool("local", false, "run tests against locally-running API")

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// parseURL returns a test's URL and possibly a second url that should return the same results.
// The second url is because /query2 should return the same as /query until switch switch over.
func parseURL(test APITest) (string, string) {
	var queryString string
	if test.query != "" {
		queryString = "?" + test.query
	}

	var base string
	if *local {
		base = test.baseURLLocal
	} else {
		base = test.baseURL
	}

	var queryURL, query2URL string
	queryURL = fmt.Sprintf(`%s/%s%s`, base, test.endpoint, queryString)
	if test.endpoint == censusEndpoint {
		query2URL = fmt.Sprintf(`%s/%s%s`, base, query2Endpoint, queryString)
	}
	return queryURL, query2URL
}

func TestAPI(t *testing.T) {
	for _, test := range Tests {

		t.Run(test.desc, func(t *testing.T) {

			url, q2url := parseURL(test)
			b, header, err := HTTPget(url)

			if err != nil {
				t.Errorf("Error getting %s: %v", url, err)
			} else {
				assertCORSHeader(header, url, t)
				assertAPIResponse(b, test, t, url)
			}

			if q2url == "" {
				return
			}

			b, header, err = HTTPget(q2url)
			if err != nil {
				t.Errorf("Error getting %s: %v", q2url, err)
			} else {
				assertCORSHeader(header, q2url, t)
				assertAPIResponse(b, test, t, q2url)
			}
		})
	}
}

// Assert CORS header allows cross-origin requests from any source (needed for web apps to use the API)
func assertCORSHeader(h map[string][]string, url string, t *testing.T) {
	cors, ok := h["Access-Control-Allow-Origin"]
	if ok {
		if !IsStringInSlice("*", cors) {
			t.Errorf("Expected '*' to be included in CORS header value for response from %s, got '%s'", url, cors)
		}
	} else {
		t.Errorf("CORS header missing in response from %s", url)
	}
}

// Assert API responses are consistent with recorded ones
func assertAPIResponse(b []byte, test APITest, t *testing.T, url string) {
	respfiles, err := MatchingRespFile(test.desc)
	if err != nil {
		t.Fail()
	}
	switch len(respfiles) {
	case 0:
		// if no recorded response is available, main.go needs to be run to capture it
		t.Errorf("No response file found for test '%s', re-run 'make update' and commit results.", test.desc)

	case 1:
		// check API response against the one on file
		respfile := respfiles[0]
		wantsha1 := RespFileSha1(respfile)
		h := sha1Hash(b)
		if h != wantsha1 {
			t.Errorf("Response from %s differed from that recorded in file %s. Run 'make testvv' or 'testvv-local' to see full diff.", url, respfile)

			// use 'go test ./... -args extra'
			// for diff

			if *extra {
				if err := diff.Text(DataPref+respfile, "", nil, string(b), os.Stdout); err != nil {
					log.Print(err)
				}
			}
		}

	default:
		// if multiple recorded responses are found, a manual check is needed - possible cause is main.go was ran, and
		// the response for a previously-recorded API query differed from previous, and so was saved as a new response file.
		t.Errorf(
			"Multiple response files found for test '%s', try manually auditing files,re-run 'make update' and commit results",
			test.desc,
		)
	}
}

func BenchmarkAllTestAPI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range Tests {
			url, _ := parseURL(test) // ignoring query2 for now
			_, _, err := HTTPget(url)
			if err != nil {
				panic(err)
			}

		}
	}
}

// TODO separate benchmarks for each API call

func BenchmarkAllConcurrentTestAPI(b *testing.B) {

	// rate limit
	const qps = 100
	rate := time.Second / qps
	ticker := time.NewTicker(rate)
	defer ticker.Stop()

	for i := 0; i < b.N; i++ {
		for _, test := range Tests {
			<-ticker.C
			url, _ := parseURL(test) // ignoring query2 for now
			go func(s string) {
				HTTPget(s)
			}(url)

		}
	}
}
