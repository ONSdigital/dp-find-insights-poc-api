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
				cors, ok := header["Access-Control-Allow-Origin"]
				if ok {
					if !IsStringInSlice("*", cors) {
						t.Errorf("Expected '*' to be included in CORS header value for response from %s, got %v", test.url, cors)
					}
				} else {
					t.Errorf("CORS header missing in response from %s", test.url)
				}
				wantfiles, err := MatchingRespFile(test.desc)
				if err != nil {
					t.Fail()
				}
				switch len(wantfiles) {
				case 0:
					t.Errorf("No response file found for test '%s', looks like you need to re-run main.go!", test.desc)
				case 1:
					wantfile := wantfiles[0]
					wantsha1 := RespFileSha1(wantfile)
					h := sha1Hash(b)
					if h != wantsha1 {
						t.Errorf("wrongly got: %s", h)

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
					t.Errorf(
						"Multiple response files found for test '%s', try manually auditing files and re-run main.go",
						test.desc,
					)
				}
			}
		})
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
