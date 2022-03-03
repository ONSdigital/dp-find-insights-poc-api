package cantabular

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestQueryMetricFilters(t *testing.T) { // QUERY 1
	ds := "Usual-Residents"
	geo := "E92000001,W92000004"
	geotype := "Country"
	code := "QS501EW"

	geoq, catsq, values := QueryMetricFilter(ds, geo, geotype, code)
	got := ParseMetric(geoq, catsq, values)

	// XXX note one field has "13,16"
	exp := "cantabular,No qualifications,Level 1 qualifications,Level 2 qualifications,Level 3 qualifications,Level 4 qualifications and above,Apprenticeships and other qualifications,Not applicable\ngeography_code,10,11,12,14,15,\"13,16\",-9\nE92000001,9539010,5680447,6571675,5306778,12117875,3952125,9710552\nW92000004,658065,327844,390605,304275,633755,201280,544600\n"
	if got != exp {
		fmt.Printf("%#v\n", got)
		t.Fail()
	}

}

func TestQueryMetricFilterOtherDS(t *testing.T) { // QUERY 1
	geo := "E92000001,W92000004"
	geotype := "Country"
	code := "QS406EW"

	geoq, catsq, values := QueryMetricFilter("", geo, geotype, code)
	got := ParseMetric(geoq, catsq, values)

	exp := "cantabular,1 person in household,2 people in household,3 people in household,4 people in household,5 people in household,6 or more people in household,Not applicable\ngeography_code,1,2,3,4,5,6-9,-9\nE92000001,5465645,13535176,11630475,11182493,6183087,4657106,1054765\nW92000004,316241,786253,672344,647015,357566,263329,59263\n"
	if got != exp {
		fmt.Printf("%#v\n", got)
		t.Fail()
	}

}

func TestSpecificMetrics(t *testing.T) { // QUERY 2
	ds := "Usual-Residents"
	code := "QS301EW"
	geotype := "Region"

	geoq, catsq, values := QueryMetric(ds, geotype, code)
	got := ParseMetric(geoq, catsq, values)

	exp := "cantabular,Provides no unpaid care,Provides unpaid care,Not applicable\ngeography_code,1,2-4,-9\nE12000001,2307839,286530,0\nE12000002,6261878,785955,0\nE12000003,4717798,558367,0\nE12000004,4034230,491420,0\nE12000005,4969823,618354,0\nE12000006,5237546,603132,0\nE12000007,7409443,687162,0\nE12000008,7770397,850101,0\nE12000009,4714266,574221,0\nW92000004,2691486,368938,0\n"

	if got != exp {
		fmt.Printf("%#v\n", got)
		t.Fail()
	}
}

func TestSpecificMetricsOtherDS(t *testing.T) { // QUERY 2
	code := "QS416EW"
	geotype := "Region"

	geoq, catsq, values := QueryMetric("", geotype, code)
	got := ParseMetric(geoq, catsq, values)

	exp := "cantabular,No cars or vans in household,1 car or van in household,2 or more cars or vans in household,Not applicable\ngeography_code,0,1,2-4,-9\nE12000001,620632,1005179,949070,49018\nE12000002,1456726,2688158,2860409,131208\nE12000003,1061663,2057408,2116106,108899\nE12000004,704466,1640205,2146475,98765\nE12000005,995417,2055221,2510158,102261\nE12000006,705831,2124298,2991478,111561\nE12000007,2734412,3428989,1986363,120091\nE12000008,1042660,3029708,4498032,207868\nE12000009,646812,1925768,2672339,125092\nW92000004,495272,1114274,1433202,59263\n"

	if got != exp {
		fmt.Printf("%#v\n", got)
		t.Fail()
	}
}

// query all codes as a crude benchmark
func TestRespFilterMetrics(t *testing.T) {
	for code := range ShortVarMap() {

		geoq, catsq, values := QueryMetricFilter("", "E92000001", "Country", code)
		got := ParseMetric(geoq, catsq, values)
		if !strings.HasPrefix(got, "cantabular") {
			t.Fatalf(got)
		}

	}
}
