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
		"rows param single cols param single",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?rows=E01000001&cols=QS119EW0001`,
	},
	// multi rows
	{
		"rows param multi single cols param single",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?rows=E01000001&rows=E01000222&rows=E01000333&cols=QS117EW0001`,
	},
	{
		"rows param multi array cols param single",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?rows=E01000100,E01000110,E01000200&cols=QS119EW0003`,
	},
	{
		"rows param multi range cols param single",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?rows=E01001111...E01001211&cols=QS118EW0011`,
	},
	{
		"rows param multi mixed cols param single",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?rows=E01000001&rows=E01000100,E01000110,E01000200&rows=E01001111...E01001211&cols=QS118EW0011`,
	},
	// multi cols
	{
		"rows param single cols param multi single",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?rows=E01000001&cols=QS119EW0001&cols=QS118EW0001&cols=QS117EW0001`,
	},
	{
		"rows param single cols param multi array",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?rows=E01000001&cols=QS119EW0001,QS119EW0002,QS119EW0003`,
	},
	{
		"rows param single cols param multi range",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?rows=E01000001&cols=QS118EW0001...QS118EW0011`,
	},
	{
		"rows param single cols param multi mixed",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?rows=E01000001&cols=QS117EW0001&cols=QS119EW0001,QS119EW0002,QS119EW0003&cols=QS118EW0001...QS118EW0011`,
	},
	// multi rows + cols
	{
		"rows param multi single cols param multi single",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?rows=E01000001&rows=E01000222&rows=E01000333&cols=QS119EW0001&cols=QS118EW0001&cols=QS117EW0001`,
	},
	{
		"rows param multi array cols param multi array",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?rows=E01000100,E01000110,E01000200&cols=QS119EW0001,QS119EW0002,QS119EW0003`,
	},
	{
		"rows param multi range cols param multi range",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?rows=E01001111...E01001211&cols=QS118EW0001...QS118EW0011`,
	},
	{
		"rows param multi mixed cols param multi mixed",
		`https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/skinny?rows=E01000001&rows=E01000100,E01000110,E01000200&rows=E01001111...E01001211&cols=QS117EW0001&cols=QS119EW0001,QS119EW0002,QS119EW0003&cols=QS118EW0001...QS118EW0011`,
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
