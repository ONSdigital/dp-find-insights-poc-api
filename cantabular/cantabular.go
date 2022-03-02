package cantabular

// adhoc query tool to investigate cantabular 2011 instance

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/ryboe/q"
	"github.com/shurcooL/graphql"
)

const URL = "https://ftb-api-ext.ons.sensiblecode.io/graphql"

//const URL = "http://127.0.0.1:8080"

type AuthTripper struct {
	User string
	Pass string
}

func (at AuthTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(at.User, at.Pass)
	return http.DefaultTransport.RoundTrip(req)
}

// list datasets
type DataSets struct {
	Datasets []struct{ Name graphql.String }
}

type IntValues []graphql.Int

type Pairs []struct { // Rename
	Code  graphql.String
	Label graphql.String
}

// like Nomis short codes and geo_type names mixed up (?)
type VariableCodes struct { // rename
	Dataset struct {
		Variables struct {
			Edges []struct {
				Node struct {
					Name  graphql.String
					Label graphql.String
				}
			}
		}
	} `graphql:"dataset(name: $ds)"`
}

// return classifications like long nomis codes
type ClassCodes struct { // rename
	Dataset struct {
		Table struct {
			Dimensions []struct {
				Categories Pairs
			}
		} `graphql:"table(variables: [$vars])"`
	} `graphql:"dataset(name: $ds)"`
}

type Metric struct {
	Dataset struct {
		Table struct {
			Dimensions []struct {
				Categories Pairs
			}
			Values IntValues
		} `graphql:"table(variables: [$geotype,$var])"`
	} `graphql:"dataset(name: $ds)"`
}

type MetricFilter struct {
	Dataset struct {
		Table struct {
			Dimensions []struct {
				Categories Pairs
			}

			Values IntValues
		} `graphql:"table(variables: [$var1,$var2],filters: [{variable: $var1, codes: $geos}])"`
	} `graphql:"dataset(name: $ds)"`
}

/*
{ dataset(name: "Usual-Residents") {
    table(
      variables: ["COUNTRY", "HLQPUK11_T007A"]
      filters: [{variable: "COUNTRY", codes: ["synE92000001"]}]
    ) {
      dimensions {
        categories {
          label
          code
        }
      }
      values
      error
    }
  }
}
*/
// MetricFilter is a cli type query1
// could be entrypoint for REST endpoint
func QueryMetricFilter(ds, geo, geoType, code string) (geoq, catsQL Pairs, values IntValues) {
	geos := strings.Split(geo, ",")

	if ds == "" {
		ds = GetDataSet(code)
	}

	var geosQL []graphql.String
	for _, v := range geos {
		geosQL = append(geosQL, graphql.String("syn"+v)) // XXX
	}

	var query MetricFilter

	vars := map[string]interface{}{
		"ds":   graphql.String(ds),
		"geos": geosQL,
		"var1": graphql.String(GeoTypeMap()[geoType]),
		"var2": graphql.String(ShortVarMap()[code]),
	}

	q.Q(vars)
	SendQueryVars(&query, vars)

	geoq = query.Dataset.Table.Dimensions[0].Categories
	catsQL = query.Dataset.Table.Dimensions[1].Categories
	values = query.Dataset.Table.Values

	return geoq, catsQL, values
}

/*
{ dataset(name: "Usual-Residents") {
    table(
      variables: ["LA", "HLQPUK11_T007A"]
    ) {
      dimensions {
        categories {
          label
          code
        }
      }
      values
      error
    }
  }
}
*/
// QueryMetric is is a cli query2
// could be entrypoint for REST endpoint
func QueryMetric(ds, geoType, code string) (geoq, catsQL Pairs, values IntValues) {
	if ds == "" {
		ds = GetDataSet(code)
	}

	vars := map[string]interface{}{
		"ds":      graphql.String(ds),
		"geotype": graphql.String(GeoTypeMap()[geoType]),
		"var":     graphql.String(ShortVarMap()[code]),
	}

	var query Metric
	SendQueryVars(&query, vars)

	geoq = query.Dataset.Table.Dimensions[0].Categories
	catsQL = query.Dataset.Table.Dimensions[1].Categories
	values = query.Dataset.Table.Values

	return geoq, catsQL, values
}

func SendQueryVars(query interface{}, vars map[string]interface{}) interface{} {
	if os.Getenv("CANT_USER") == "" || os.Getenv("CANT_PW") == "" {
		log.Fatal("define CANT_USER & CANT_PW")
	}
	hclient := &http.Client{Transport: AuthTripper{User: os.Getenv("CANT_USER"), Pass: os.Getenv("CANT_PW")}}
	client := graphql.NewClient(URL, hclient)
	if err := client.Query(context.Background(), query, vars); err != nil {
		log.Fatal(err)
	}
	return query
}

// ParseResp is used for the command line investigate API commands
// probably doesn't make sense for REST API
func ParseResp(query interface{}) {
	// wish there were a better way!
	qt := reflect.TypeOf(query).String()
	switch qt {
	// these help with looking at data in the API
	case "*cantabular.DataSets":
		ds := query.(*DataSets)
		for _, v := range ds.Datasets {
			fmt.Println(v.Name)
		}
	case "*cantabular.VariableCodes":
		scodes := query.(*VariableCodes)
		for _, v := range scodes.Dataset.Variables.Edges {
			fmt.Print(v.Node.Name + " : ")
			fmt.Println(v.Node.Label)
		}
	case "*cantabular.ClassCodes":
		scodes := query.(*ClassCodes)
		for _, v := range scodes.Dataset.Table.Dimensions {
			for _, v2 := range v.Categories {
				fmt.Print(v2.Code + " = ")
				fmt.Println(v2.Label)
			}
		}

	default:
		log.Fatal(qt + " unrecognised")
	}

}

func ParseMetric(geo, cats Pairs, values IntValues) (resp string) {
	first := []string{"\"cantabular\""}
	second := []string{"\"geography_code\""}
	var lines []string

	for _, k := range cats {
		first = append(first, fmt.Sprintf("%q", string(k.Label)))
		second = append(second, fmt.Sprintf("%q", string(k.Code)))
	}

	k := 0
	for _, g := range geo {
		// factor
		geo := strings.Split(string(g.Code), "syn")[1]
		line := []string{geo}

		for j := 0; j < len(cats); j++ {
			line = append(line, fmt.Sprintf("%d", values[k]))
			k++
		}
		lines = append(lines, strings.Join(line, ", "))
	}

	resp = fmt.Sprintf("%s\n", strings.Join(first, ","))
	resp += fmt.Sprintf("%s\n", strings.Join(second, ","))
	resp += fmt.Sprintf("%s\n", strings.Join(lines, "\n"))

	q.Q(resp)

	return resp
}

func GeoTypeMap() map[string]string {

	return map[string]string{
		"Country": "Country",
		"Region":  "Region", // XXX checkme
		"LAD":     "LA",
		"MSOA":    "MSOA",
	}
}

func ShortVarMap() map[string]string {

	// maybe this should be in the database?
	// although list is short & likely to change..

	return map[string]string{
		"KS102EW": "AGE_T009A",
		"KS202EW": "NATID_ALL_T009A",
		"KS206EW": "WELSHPUK112_T007A",
		"KS207WA": "WELSHPUK112_R003A",
		"KS208WA": "WELSHPUK112_R003A",
		"QS104EW": "SEX",
		"QS113EW": "MARSTAT_T006A",
		"QS201EW": "ETHPUK11_T009A",
		"QS203EW": "COB_R010A",
		"QS208EW": "RELPUK11_R005A",
		"QS301EW": "CARER_R003A",
		"QS302EW": "HEALTH_T004A", // HEALTH
		"QS303EW": "DISABILITY_T003B",
		"QS402EW": "TYPACCOM_T009A",
		"QS406EW": "SIZHUK11_T007A",
		"QS415EW": "CENHEATHUK11_T003A",
		"QS416EW": "CARSNO_T004A",
		"QS501EW": "HLQPUK11_T007A", // EDUCATION
		"QS601EW": "ECOPUK11_R006A",
		"QS604EW": "HOURS",
		"QS605EW": "INDGPUK11_T009A",
		"QS606EW": "OCCPUK113_T010A",
		"QS701EW": "TRANSPORT_R005A",
		"QS702EW": "AGGDTWPEW11_R010A",
		//"DC6102EW": "STUDENT_AGE_T002A",
		//"QS402EW":  "TENHUK11_T007B",
		//"QS411EW": "BEDROOMS_T006A",
		//"QS501EW":  "HLQPUK11_T007A",
	}

}

func GetDataSet(varCode string) string {

	mappy := map[string]string{
		"QS406EW": "People-Households",
		"KS206EW": "People-Households",
		"QS402EW": "People-Households",
		"QS416EW": "People-Households",
		"QS415EW": "People-Households",
	}

	if mappy[varCode] != "" {
		return mappy[varCode]
	}

	return "Usual-Residents"
}
