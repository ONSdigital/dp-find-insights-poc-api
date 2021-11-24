// +build integration

package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type APITest = struct {
	desc     string
	url      string
}

var Tests = []APITest{
	{
		"no params",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew`,
	},
	{
		"cols param single",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=_1`,
	},
	{
		"cols param multi",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=geography_code,total,_1`,
	},
	{
		"rows param single",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?rows=geography_code:E01000001`,
	},
	{
		"rows param multi single",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?rows=geography_code:E01000011&rows=geography_code:E01000012&rows=geography_code:E01000013`,
	},
	{
		"rows param multi array",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?rows=geography_code:E01000001,E01000002,E01000003,E01000005`,
	},
	{
		"rows param range",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?rows=geography_code:E01000005...E01000010`,
	},
	{
		"rows param mixed",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?rows=geography_code:E01000001,E01000005...E01000010&rows=geography_code:E01000079`,
	},
	{
		"rows cols params single col single row",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=_1&rows=geography_code:E01000081`,
	},
	{
		"rows cols params multi col single row",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=geography_code,total,_1&rows=geography_code:E01000081`,
	},
	{
		"rows cols params multi col multi single row",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=geography_code,total,_1&rows=geography_code:E01000081&rows=geography_code:E01000073`,
	},
	{
		"rows cols params multi col multi array row",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=geography_code,total,_1&rows=geography_code:E01000124,E01000130,E01000264`,
	},
	{
		"rows cols params multi col multi range row",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=geography_code,total,_1&rows=geography_code:E01000296...E01000355`,
	},
	{
		"rows cols params multi col multi mixed row",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=geography_code,total,_1&rows=geography_code:E01000275,E01000281,E01001146...E01001194&rows=geography_code:E01000027`,
	},
	{
		"rows cols params single col multi single row",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=_1&rows=geography_code:E01000081&rows=geography_code:E01000073`,
	},
	{
		"rows cols params single col multi array row",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=_1&rows=geography_code:E01000124,E01000130,E01000264`,
	},
	{
		"rows cols params single col multi range row",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=_1&rows=geography_code:E01000296...E01000355`,
	},
	{
		"rows cols params single col multi mixed row",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=_1&rows=geography_code:E01000275,E01000281,E01001146...E01001194&rows=geography_code:E01000027`,
	},
}

var DataPref = "resp/"

func main() {
	// populate saved copies to allow test diffs
	for _, test := range Tests {
		b, _, err := HTTPget(test.url)

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

func HTTPget(s string) (b []byte, h map[string][]string, err error) {
	resp, err := http.Get(s)

	if err != nil {
		return b, h, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := errors.New(fmt.Sprintf("API responded with status code %v", resp.StatusCode))
		return b, h, err
	}

	// or just io. in go 1.16+
	b, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return b, h, err
	}

	h = resp.Header

	return b, h, err
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

// find file(s) in DataPref directory that matches testDesc
func MatchingRespFile(testDesc string) (fns []string, err error) {
	filesInDataPref, err := ioutil.ReadDir(DataPref)
	if err != nil {
		return fns, err
	}
	targetFnPrefix := RespFilePrefix(testDesc)
	for _, file := range filesInDataPref {
		fn := file.Name()
		if strings.HasPrefix(fn, targetFnPrefix) {
			fns = append(fns, fn)
		}
	}
	return fns, err
}

func IsStringInSlice(str string, s []string) (bool) {
	for _, ele := range(s) {
		if ele == str {
			return true
		}
	}
	return false
}
