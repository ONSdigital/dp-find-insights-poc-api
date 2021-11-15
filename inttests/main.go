package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var Tests = []struct {
	desc     string
	url      string
	wantsha1 string
}{
	{"no parms", `https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew`, "91231b61eb345125ee4c10cf4594c6c961b8c623"},
	{"cols parm", `https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?cols=geography_code,total,_1`, "3adf545d8107914418bf9c3cf8d6ebb8ea075761"},
	{"rows param", `https://5laefo1cxd.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs119ew?rows=geography_code:E01000001,E01000003...E01000006`, "5407abf563b559e57b5a487f83f48965b9fc4afc"},
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

		// only save if right hash
		// check-in changes
		if h == test.wantsha1 {
			f, err := os.Create(DataPref + h)
			if err != nil {
				panic(err)
			}
			f.WriteString(string(b))
			f.Close()
		}
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
