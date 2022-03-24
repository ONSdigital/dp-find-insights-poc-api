//go:build integration
// +build integration

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

// makeURL returns a test's URL and possibly a second url that should return the same results.
// The second url is because /query2 should return the same as /query until switch switch over.
func makeURL(endpoint, query string) (string, string) {
	var queryString string
	if query != "" {
		queryString = "?" + query
	}

	var base string
	if *local {
		base = baseURLLocal
	} else {
		base = baseURL
	}

	var queryURL, query2URL string
	queryURL = fmt.Sprintf(`%s/%s%s`, base, endpoint, queryString)
	if endpoint == censusEndpoint {
		query2URL = fmt.Sprintf(`%s/%s%s`, base, query2Endpoint, queryString)
	}
	return queryURL, query2URL
}

func TestOPTIONS(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url, _ := makeURL(geoEndpoint, "")
	req, err := http.NewRequestWithContext(ctx, "OPTIONS", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	headers := resp.Header.Values("Access-Control-Allow-Headers")
	for _, h := range headers {
		for _, tok := range strings.Split(h, ",") {
			if strings.EqualFold(tok, "Cache-Control") {
				return
			}
		}
	}
	t.Fatal("expected Cache-Control in Access-Control-Allow-Headers")
}

func TestAPI(t *testing.T) {
	for _, test := range Tests {

		t.Run(test.desc, func(t *testing.T) {

			url, q2url := makeURL(test.endpoint, test.query)
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
			url, _ := makeURL(test.endpoint, test.query) // ignoring query2 for now
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
			url, _ := makeURL(test.endpoint, test.query) // ignoring query2 for now
			go func(s string) {
				HTTPget(s)
			}(url)

		}
	}
}

// Verify API responds with non-200 when there is a problem with the query.
func Test_API_Error_Detection(t *testing.T) {
	var (
		msf      = "-3.280,54.911"
		guinness = "-6.286840,53.341738"
		askja    = "-16.73041,65.03272"
		tenerife = "-16.577575,28.047316"
		porsche  = "9.152403,48.834302"
		santa    = "25.84677,66.54384"
		galway   = "-9.054191,53.272517"
		eiffel   = "2.294532,48.8582722"
		twilight = "181,0" //
	)
	var tests = map[string]struct {
		endpoint string
		query    string
	}{
		"non-numeric bbox coordinate": {
			censusEndpoint,
			"bbox=x,y",
		},
		"incomplete bbox coordinates": {
			censusEndpoint,
			fmt.Sprintf("bbox=%s,0", msf),
		},
		"bbox not two points": {
			censusEndpoint,
			fmt.Sprintf("bbox=%s,%s,%s", msf, guinness, galway),
		},
		"bbox in the twilight zone": {
			censusEndpoint,
			fmt.Sprintf("bbox=%s,%s", twilight, porsche),
		},
		"bbox not in UK": {
			censusEndpoint,
			fmt.Sprintf("bbox=%s,%s", santa, eiffel),
		},

		"missing circle location": {
			censusEndpoint,
			"radius=1000",
		},
		"non-numeric circle location": {
			censusEndpoint,
			"location=x,y&radius=1000",
		},
		"incomplete circle coordinates": {
			censusEndpoint,
			"location=1&radius=1000",
		},
		"too many circle coordinates": {
			censusEndpoint,
			fmt.Sprintf("location=%s,%s&radius=1000", msf, guinness),
		},
		"circle in twilight zone": {
			censusEndpoint,
			fmt.Sprintf("location=%s&radius=1000", twilight),
		},
		"circle not in UK": {
			censusEndpoint,
			fmt.Sprintf("location=%s&radius=1000", tenerife),
		},
		"small circle radius": {
			censusEndpoint,
			fmt.Sprintf("location=%s&radius=0", msf),
		},
		"huge circle radius": {
			censusEndpoint,
			fmt.Sprintf("location=%s&radius=1000001", msf),
		},

		"non-numeric polygon coordinate": {
			censusEndpoint,
			fmt.Sprintf("polygon=x,y,%s,%s,%s,x,y", askja, msf, santa),
		},
		"not enough polygon coordinates": {
			censusEndpoint,
			fmt.Sprintf("polygon=%s,%s,%s", eiffel, porsche, eiffel),
		},
		"first and last polygon coordinates don't match": {
			censusEndpoint,
			fmt.Sprintf("polygon=%s,%s,%s,%s", galway, tenerife, porsche, santa),
		},
		"polygon coordinate in the twilight zone": {
			censusEndpoint,
			fmt.Sprintf("polygon=%s,%s,%s,%s", galway, twilight, porsche, galway),
		},
		"polygon not in UK": {
			censusEndpoint,
			fmt.Sprintf("polygon=%s,%s,%s,%s", tenerife, porsche, eiffel, tenerife),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			url, q2url := makeURL(test.endpoint, test.query)
			var err error
			if _, _, err = HTTPget(url); err == nil {
				t.Errorf("%s: expected API to return error", url)
			}
			t.Logf("%s: %s", url, err)
			if q2url == "" {
				return
			}
			if _, _, err = HTTPget(q2url); err == nil {
				t.Errorf("%s: expected API to return error", url)
			}
			t.Logf("%s: %s", q2url, err)
		})
	}
}
