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

			b, err := HTTPget(test.url)

			if err != nil {
				t.Fail()
			}

			h := sha1Hash(b)
			if h != test.wantsha1 {
				t.Errorf("wrongly got: %s", h)

				// use 'go test ./... -args extra'
				// for diff

				if os.Args[len(os.Args)-1] == "extra" {

					dmp := diffmatchpatch.New()

					f, _ := os.Open(DataPref + test.wantsha1)

					wanted, _ := io.ReadAll(f)

					diffs := dmp.DiffMain(string(b), string(wanted), true)

					t.Log(dmp.DiffPrettyText(diffs))
				}
			}
		})
	}
}

func BenchmarkAllTestAPI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range Tests {
			_, err := HTTPget(test.url)
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
