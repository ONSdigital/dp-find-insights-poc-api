// FIX THIS SHIT UP
package cantabular

import (
	"fmt"
	"log"
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

	exp := "\"cantabular\",\"No qualifications\",\"Level 1 qualifications\",\"Level 2 qualifications\",\"Level 3 qualifications\",\"Level 4 qualifications and above\",\"Apprenticeships and other qualifications\",\"Not applicable\"\n\"geography_code\",\"10\",\"11\",\"12\",\"14\",\"15\",\"13,16\",\"-9\"\nE92000001, 9539010, 5680447, 6571675, 5306778, 12117875, 3952125, 9710552\nW92000004, 658065, 327844, 390605, 304275, 633755, 201280, 544600\n"
	if got != exp {
		fmt.Printf("%#v\n", got)
		t.Fail()
	}

}

func TestSpecificMetrics(t *testing.T) { // QUERY 2

	// need mapping for region and dataset
	ds := "Usual-Residents"
	code := "QS301EW"
	geotype := "Region"

	geoq, catsq, values := QueryMetric(ds, geotype, code)
	got := ParseMetric(geoq, catsq, values)

	exp := "\"cantabular\",\"Provides no unpaid care\",\"Provides unpaid care\",\"Not applicable\"\n\"geography_code\",\"1\",\"2-4\",\"-9\"\nE12000001, 2307839, 286530, 0\nE12000002, 6261878, 785955, 0\nE12000003, 4717798, 558367, 0\nE12000004, 4034230, 491420, 0\nE12000005, 4969823, 618354, 0\nE12000006, 5237546, 603132, 0\nE12000007, 7409443, 687162, 0\nE12000008, 7770397, 850101, 0\nE12000009, 4714266, 574221, 0\nW92000004, 2691486, 368938, 0\n"

	if got != exp {
		t.Errorf(got)
	}
}

/*
func TestRespFilterMetrics(t *testing.T) {
	for code := range shortVarMap() {

		var query MetricFilter
		SendQueryVars(&query, BuildCantParam("E92000001", code))
		got := ParseMetricFilter(&query)

		if !strings.HasPrefix(got, "\"cantabular") {
			t.Fatalf(got)
		}

	}
}


func TestSpecificFilterMetrics(t *testing.T) {

	{
		code := "QS501EW"
		geo := []string{"E92000001"}

			if getGeoTypeName(geo) != "Country" {
				t.Fail()

			}

		var query MetricFilter2
		SendQueryVars(&query, BuildCantParam(geo, code))

		exp := "\"cantabular\",\"No qualifications\",\"Level 1 qualifications\",\"Level 2 qualifications\",\"Level 3 qualifications\",\"Level 4 qualifications and above\",\"Apprenticeships and other qualifications\",\"Not applicable\"\n\"geography_code\",\"10\",\"11\",\"12\",\"14\",\"15\",\"13,16\",\"-9\"\nE92000001,9539010,5680447,6571675,5306778,12117875,3952125,9710552\n"

		geoq := query.Dataset.Table.Dimensions[0].Categories
		catsq := query.Dataset.Table.Dimensions[1].Categories
		values := query.Dataset.Table.Values
		got := ParseMetric2(geoq, catsq, values)

		if got != exp {
			fmt.Printf("%#v\n", got)
			t.Fail()
		}
	}

	{
		code := "QS302EW"
		geo := "E06000004"

		name := getGeoTypeName(geo)
		if name != "LA" {
			t.Errorf("got " + name)

		}

		var query MetricFilter2
		SendQueryVars(&query, BuildCantParam(geo, code))

		exp := "\"cantabular\",\"Very good or good health\",\"Fair health\",\"Bad or very bad health\",\"Not applicable\"\n\"geography_code\",\"1-2\",\"3\",\"4-5\",\"-9\"\nE06000004,141570,26818,13058,0\n"

		geoq := query.Dataset.Table.Dimensions[0].Categories
		catsq := query.Dataset.Table.Dimensions[1].Categories
		values := query.Dataset.Table.Values
		got := ParseMetric2(geoq, catsq, values)
		if got != exp {
			t.Errorf(got)
		}

	}

	{
		code := "QS406EW"
		geo := "E02003802"

		name := getGeoTypeName(geo)
		if name != "MSOA" {
			t.Errorf("got " + name)
		}

		ds := getDataSet("QS406EW")
		if ds != "People-Households" {
			t.Errorf("got " + ds)
		}

		var query MetricFilter2
		vars := BuildCantParam(geo, code)
		q.Q(vars)
		SendQueryVars(&query, vars)

		exp := "\"cantabular\",\"1 person in household\",\"2 people in household\",\"3 people in household\",\"4 people in household\",\"5 people in household\",\"6 or more people in household\",\"Not applicable\"\n\"geography_code\",\"1\",\"2\",\"3\",\"4\",\"5\",\"6-9\",\"-9\"\nE02003802,1049,2398,1845,2004,1015,580,184\n"
		geoq := query.Dataset.Table.Dimensions[0].Categories
		catsq := query.Dataset.Table.Dimensions[1].Categories
		values := query.Dataset.Table.Values
		got := ParseMetric2(geoq, catsq, values)

		if got != exp {
			fmt.Printf("%#v\n", got)
			t.Fail()
		}

	}

}
*/
