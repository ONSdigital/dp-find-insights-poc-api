//go:build integration
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
	endpoint string
	query    string
}

const defaultBaseURL = `http://ec2-18-193-6-194.eu-central-1.compute.amazonaws.com:25252`
const baseURLLocal = `http://localhost:25252`
const censusEndpoint = `query/2011`
const metadataEndpoint = `metadata/2011`
const ckmeansEndpoint = `ckmeans/2011`
const geoEndpoint = "geo/2011"
const query2Endpoint = `query2/2011`

// For backward compatibility, baseURL defaults to the dev EC2 instance,
// but allow $TEST_TARGET_URL to override.
var baseURL = defaultBaseURL

func init() {
	target := os.Getenv("TEST_TARGET_URL")
	if target != "" {
		baseURL = target
	}
}

var Tests = []APITest{
	// skinny is deprecated
	// no rows found; expect csv headings only
	{
		"no rows found",
		censusEndpoint,
		`rows=noexist&cols=geography_code,geotype`,
	},
	// no cols found; expect csv headings only
	{
		"no cols found",
		censusEndpoint,
		`rows=E01002111&cols=geography_code,geotype,noexist`,
	},
	// all cols single row no geography census
	{
		"census tables all cols single row no geography",
		censusEndpoint,
		`rows=E01002111`,
	},
	// all cols single row census
	{
		"census tables all cols single row",
		censusEndpoint,
		`cols=geography_code&rows=E01002111`,
	},
	// single row + col census
	{
		"census tables rows param single cols param single",
		censusEndpoint,
		`rows=E01000001&cols=geography_code,QS701EW0004`,
	},
	// multi rows census
	{
		"census tables rows param multi single cols param single",
		censusEndpoint,
		`rows=E01000001&rows=E01000222&rows=E01000333&cols=geography_code,QS415EW0006`,
	},
	{
		"census tables rows param multi array cols param single",
		censusEndpoint,
		`rows=E01000100,E01000110,E01000200&cols=geography_code,QS208EW0009`,
	},
	{
		"census tables rows param multi range cols param single",
		censusEndpoint,
		`rows=E01001111...E01001211&cols=geography_code,QS701EW0010`,
	},
	{
		"census tables rows param multi mixed cols param single",
		censusEndpoint,
		`rows=E01000001&rows=E01000100,E01000110,E01000200&rows=E01001111...E01001211&cols=geography_code,KS202EW0008`,
	},
	// multi cols census
	{
		"census tables rows param single cols param multi single",
		censusEndpoint,
		`rows=E01000001&cols=QS701EW0001&cols=geography_code,QS112EW0001&cols=QS112EW0002`,
	},
	{
		"census tables rows param single cols param multi array",
		censusEndpoint,
		`rows=E01000001&cols=geography_code,QS701EW0001,QS415EW0002,QS415EW0003`,
	},
	{
		"census tables rows param single cols param multi range",
		censusEndpoint,
		`rows=E01000001&cols=geography_code,QS119EW0001...QS119EW0006`,
	},
	{
		"census tables rows param single cols param multi mixed",
		censusEndpoint,
		`rows=E01000001&cols=QS112EW0002&cols=geography_code,QS701EW0001,QS415EW0002,QS415EW0003&cols=QS112EW0001...QS118EW0011`,
	},
	// multi rows + cols census
	{
		"census tables rows param multi single cols param multi single",
		censusEndpoint,
		`rows=E01000001&rows=E01000222&rows=E01000333&cols=geography_code,QS701EW0001&cols=QS112EW0001&cols=QS112EW0002`,
	},
	{
		"rows param multi array cols param multi array",
		censusEndpoint,
		`rows=E01000100,E01000110,E01000200&cols=geography_code,QS701EW0001,QS415EW0002,QS415EW0003`,
	},
	{
		"census tables rows param multi range cols param multi range",
		censusEndpoint,
		`rows=E01001111...E01001211&cols=geography_code,QS112EW0001...QS118EW0011`,
	},
	{
		"census tables rows param multi mixed cols param multi mixed",
		censusEndpoint,
		`rows=E01000001&rows=E01000100,E01000110,E01000200&rows=E01001111...E01001211&cols=QS112EW0002&cols=geography_code,QS701EW0001,QS415EW0002,QS415EW0003&cols=QS112EW0001...QS118EW0011`,
	},
	// bbox census
	{
		"census tables bbox param cols param multi mixed",
		censusEndpoint,
		`bbox=0.1338,51.4635,0.1017,51.4647&cols=geography_code,QS112EW0002&cols=QS701EW0001,QS415EW0002,QS415EW0003&cols=QS112EW0001...QS118EW0011&geotype=LSOA`,
	},
	// polygon census
	{
		"census tables polygon param cols param multi mixed",
		censusEndpoint,
		`polygon=0.0844,51.4897,0.1214,51.4910,0.1338,51.4635,0.1017,51.4647,0.0844,51.4897&cols=geography_code,QS112EW0002&cols=QS701EW0001,QS415EW0002,QS415EW0003&cols=QS112EW0001...QS118EW0011&geotype=LSOA`,
	},
	// radius census
	{
		"census tables radius param cols param multi mixed",
		censusEndpoint,
		`location=0.1338,51.4635&radius=1000&cols=geography_code,QS112EW0002&cols=QS701EW0001,QS415EW0002,QS415EW0003&cols=QS112EW0001...QS118EW0011&geotype=LSOA`,
	},
	// rows, bbox, polygon & radius m-m-m-megamix census
	{
		"census tables rows bbox poylgon radius param cols param multi mixed",
		censusEndpoint,
		`rows=E01000100,E01000110,E01000200&bbox=0.1338,51.4635,0.1017,51.4647&polygon=0.0844,51.4897,0.1214,51.4910,0.1338,51.4635,0.1017,51.4647,0.0844,51.4897&location=0.1338,51.4635&radius=1000&cols=geography_code,QS112EW0002&cols=QS701EW0001,QS415EW0002,QS415EW0003&cols=QS112EW0001...QS118EW0011&geotype=LSOA`,
	},
	// censustable single geography
	{
		"census tables censustable with single geography",
		censusEndpoint,
		`rows=E01000001&censustable=QS101EW`,
	},
	// censustable single column single geography
	{
		"census tables censustable with single geography and additional single column",
		censusEndpoint,
		`cols=QS415EW0002&rows=E01000001&censustable=QS101EW`,
	},
	// censustable single column multi geography
	{
		"census tables censustable with single geography and multiple additional columns",
		censusEndpoint,
		`cols=QS701EW0001,QS415EW0002,QS415EW0003&rows=E01000001&censustable=QS101EW`,
	},
	// all LAD rows, single category
	{
		"all rows single category LAD",
		censusEndpoint,
		`rows=ALL&cols=geography_code,geotype,QS701EW0006&geotype=LAD`,
	},
	// metadata
	{
		"metadata",
		metadataEndpoint,
		"",
	},
	// ckmeans
	{
		"ckmeans",
		ckmeansEndpoint,
		"cat=QS208EW0002&geotype=LSOA&k=5",
	},
	{
		"ckmeans empty response",
		ckmeansEndpoint,
		"cat=noexist&geotype=LSOA&k=5",
	},
	// geo
	{
		"geo",
		geoEndpoint,
		"geocode=E09000004",
	},
	{
		"geo empty response",
		geoEndpoint,
		"geocode=noexist",
	},
}

var DataPref = "resp/"

func main() {
	// populate saved copies to allow test diffs - NB this will always run against the non-local server
	for _, test := range Tests {
		var queryString string
		if test.query != "" {
			queryString = "?" + test.query
		}

		url := fmt.Sprintf(`%s/%s%s`, baseURL, test.endpoint, queryString)
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

	// or just io. in go 1.16+
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, h, err
	}

	h = resp.Header

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("API responded with status code %v (%s)", resp.StatusCode, b))
	}

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
		if len(fn) < 40 {
			log.Printf("rogue file: %s", DataPref+fn)
			os.Exit(1)
		}
		fnPrefix := fn[:len(fn)-40]
		if fnPrefix == targetFnPrefix {
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
