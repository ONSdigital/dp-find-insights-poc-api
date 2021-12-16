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
	desc         string
	baseURL      string
	baseURLLocal string
	query        string
}

const baseURL = `https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny`
const baseURLLocal = `http://localhost:25252/dev/hello/skinny`

var Tests = []APITest{
	// all rows single col
	{
		"all rows single col",
		baseURL,
		baseURLLocal,
		`cols=QS802EW0009`,
	},
	// all cols single row
	{
		"all cols single row",
		baseURL,
		baseURLLocal,
		`rows=E01002111`,
	},
	// single row + col
	{
		"rows param single cols param single",
		baseURL,
		baseURLLocal,
		`rows=E01000001&cols=QS119EW0001`,
	},
	// multi rows
	{
		"rows param multi single cols param single",
		baseURL,
		baseURLLocal,
		`rows=E01000001&rows=E01000222&rows=E01000333&cols=QS117EW0001`,
	},
	{
		"rows param multi array cols param single",
		baseURL,
		baseURLLocal,
		`rows=E01000100,E01000110,E01000200&cols=QS119EW0003`,
	},
	{
		"rows param multi range cols param single",
		baseURL,
		baseURLLocal,
		`rows=E01001111...E01001211&cols=QS118EW0011`,
	},
	{
		"rows param multi mixed cols param single",
		baseURL,
		baseURLLocal,
		`rows=E01000001&rows=E01000100,E01000110,E01000200&rows=E01001111...E01001211&cols=QS118EW0011`,
	},
	// multi cols
	{
		"rows param single cols param multi single",
		baseURL,
		baseURLLocal,
		`rows=E01000001&cols=QS119EW0001&cols=QS118EW0001&cols=QS117EW0001`,
	},
	{
		"rows param single cols param multi array",
		baseURL,
		baseURLLocal,
		`rows=E01000001&cols=QS119EW0001,QS119EW0002,QS119EW0003`,
	},
	{
		"rows param single cols param multi range",
		baseURL,
		baseURLLocal,
		`rows=E01000001&cols=QS118EW0001...QS118EW0011`,
	},
	{
		"rows param single cols param multi mixed",
		baseURL,
		baseURLLocal,
		`rows=E01000001&cols=QS117EW0001&cols=QS119EW0001,QS119EW0002,QS119EW0003&cols=QS118EW0001...QS118EW0011`,
	},
	// multi rows + cols
	{
		"rows param multi single cols param multi single",
		baseURL,
		baseURLLocal,
		`rows=E01000001&rows=E01000222&rows=E01000333&cols=QS119EW0001&cols=QS118EW0001&cols=QS117EW0001`,
	},
	{
		"rows param multi array cols param multi array",
		baseURL,
		baseURLLocal,
		`rows=E01000100,E01000110,E01000200&cols=QS119EW0001,QS119EW0002,QS119EW0003`,
	},
	{
		"rows param multi range cols param multi range",
		baseURL,
		baseURLLocal,
		`rows=E01001111...E01001211&cols=QS118EW0001...QS118EW0011`,
	},
	{
		"rows param multi mixed cols param multi mixed",
		baseURL,
		baseURLLocal,
		`rows=E01000001&rows=E01000100,E01000110,E01000200&rows=E01001111...E01001211&cols=QS117EW0001&cols=QS119EW0001,QS119EW0002,QS119EW0003&cols=QS118EW0001...QS118EW0011`,
	},
	// bbox
	{
		"bbox param cols param multi mixed",
		baseURL,
		baseURLLocal,
		`bbox=0.1338,51.4635,0.1017,51.4647&cols=QS117EW0001&cols=QS119EW0001,QS119EW0002,QS119EW0003&cols=QS118EW0001...QS118EW0011&geotype=LSOA`,
		//`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?bbox=51.4635,0.1338,51.4647,0.1017&cols=QS117EW0001&cols=QS119EW0001,QS119EW0002,QS119EW0003&cols=QS118EW0001...QS118EW0011&geotype=LSOA`,
	},
}

var DataPref = "resp/"

func main() {
	// populate saved copies to allow test diffs - NB this will always run against the non-local server
	for _, test := range Tests {
		url := fmt.Sprintf(`%s?%s`, test.baseURL, test.query)
		b, _, err := HTTPget(url)

		if err != nil {
			log.Print(err)
		}

		h := sha1Hash(b)
		fmt.Printf("url: %s resp hash: %s\n", url, h)

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
	return fn[strings.LastIndex(fn, "-")+1:]
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

func IsStringInSlice(str string, s []string) bool {
	for _, ele := range s {
		if ele == str {
			return true
		}
	}
	return false
}
