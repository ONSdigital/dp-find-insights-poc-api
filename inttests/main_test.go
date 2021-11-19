// +build integration

package main

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestAPI(t *testing.T) {

	for _, test := range Tests {

		t.Run(test.desc, func(t *testing.T) {

			b, header, err := HTTPget(test.url)

			if err != nil {
				t.Errorf("Error getting %s: %v", test.url, err)
				} else {
					assertCORSHeader(header, test, t)
					assertAPIResponse(b, test, t)
			}
		})
	}
}

// Assert CORS header allows cross-origin requests from any source (needed for web apps to use the API)
func assertCORSHeader(h map[string][]string, test APITest, t *testing.T) {
	cors, ok := h["Access-Control-Allow-Origin"]
	if ok {
		if !IsStringInSlice("*", cors) {
			t.Errorf("Expected '*' to be included in CORS header value for response from %s, got '%s'", test.url, cors)
		}
	} else {
		t.Errorf("CORS header missing in response from %s", test.url)
	}
}

// Assert API responses are consistent with recorded ones
func assertAPIResponse(b []byte, test APITest, t*testing.T) {
	respfiles, err := MatchingRespFile(test.desc)
	if err != nil {
		t.Fail()
	}
	switch len(respfiles) {
	case 0:
		// if no recorded response is available, main.go needs to be run to capture it
		t.Errorf("No response file found for test '%s', looks like you need to re-run main.go!", test.desc)

	case 1:
		// check API response against the one on file
		respfile := respfiles[0]
		wantsha1 := RespFileSha1(respfile)
		h := sha1Hash(b)
		if h != wantsha1 {
			t.Errorf("Response from %s differed from that recorded in file %s. Run 'make testvv' to see full diff.", test.url, respfile)

			// use 'go test ./... -args extra'
			// for diff

			if os.Args[len(os.Args)-1] == "extra" {

				dmp := diffmatchpatch.New()

				f, _ := os.Open(DataPref + wantsha1)

				wanted, _ := io.ReadAll(f)

				diffs := dmp.DiffMain(string(b), string(wanted), true)

				t.Log(dmp.DiffPrettyText(diffs))
			}
		}

	default:
		// if multiple recorded responses are found, a manual check is needed - possible cause is main.go was ran, and
		// the response for a previously-recorded API query differed from previous, and so was saved as a new response file.
		t.Errorf(
			"Multiple response files found for test '%s', try manually auditing files and re-run main.go",
			test.desc,
		)
	}
}

func BenchmarkAllTestAPI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range Tests {
			_, _, err := HTTPget(test.url)
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

			go func(s string) {
				HTTPget(s)
			}(test.url)

		}
	}
}
