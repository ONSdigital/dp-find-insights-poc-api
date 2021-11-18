// +build integration

package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var Tests = []struct {
	desc     string
	url      string
}{
	{"no params", `https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew`},
	{"cols param", `https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=geography_code,total,_1`},
	{"rows param", `https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?rows=geography_code:E01000001,E01000003...E01000006`},
	// TODO Viv "rows and cols param"
	// etc.
}

var DataPref = "resp/"

func main() {
	// populate saved copies to allow test diffs
	for _, test := range Tests {
		b, err := HTTPget(test.url)

		if err != nil {
			log.Print(err)
		}

		h := sha1Hash(b)
		fmt.Printf("url: %s resp hash: %s\n", test.url, h)

		fn := DataPref + RespFilePrefix(test.desc) + h
		f, err := os.Create(fn)
		if err != nil {
			panic(err)
		}
		f.WriteString(string(b))
		f.Close()
	}
}

func HTTPget(s string) (b []byte, err error) {
	resp, err := http.Get(s)

	if err != nil {
		return b, err
	}

	defer resp.Body.Close()

	// or just io. in go 1.16+
	b, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return b, err
	}

	return b, err
}

func sha1Hash(b []byte) string {
	s := string(b)
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

// make resp file name prefix from test desc
func RespFilePrefix(testDesc string) string {
	return strings.Replace(testDesc, " ", "-", -1) + "-"
}

// parse sha1 from resp file name
func RespFileSha1(fn string) string {
	return fn[strings.LastIndex(fn, "-") + 1:]
}

// file file in DataPref directory that matches testDesc
func MatchingRespFile(testDesc string) (fn string, err error) {
	filesInDataPref, err := ioutil.ReadDir(DataPref)
	if err != nil {
		return fn, err
	}
	targetFnPrefix := RespFilePrefix(testDesc)
	for _, file := range filesInDataPref {
		fn := file.Name()
		if strings.HasPrefix(fn, targetFnPrefix) {
			return fn, err
		}
	}
	return fn, err
}
