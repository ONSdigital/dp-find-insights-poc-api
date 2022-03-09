package cantabular

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

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

type Pairs []struct {
	Code  graphql.String
	Label graphql.String
}

type VariableCodes struct {
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
		} `graphql:"table(variables: [$geotype,$var],filters: [{variable: $geotype, codes: $geos}])"`
	} `graphql:"dataset(name: $ds)"`
}

// Metadata is a slow, tactical solution
type Metadata struct {
	Code       string
	Name       string
	Categories []struct {
		Code string
		Name string
	} `json:"categories"`
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

	// 2011 cantabular geocodes have syn appended to the front
	var geosQL []graphql.String
	for _, v := range geos {
		geosQL = append(geosQL, graphql.String("syn"+v))
	}

	var query MetricFilter

	vars := map[string]interface{}{
		"ds":      graphql.String(ds),
		"geos":    geosQL,
		"geotype": graphql.String(GeoTypeMap()[geoType]),
		"var":     graphql.String(ShortVarMap()[code]),
	}

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

// QueryMetaData does some multiple data queries to get data structure
// XXX a poor work around for a lack of metadata.

func QueryMetaData(ds string, nomis bool) string {
	if ds == "" {
		ds = GetDataSet("")
	}

	revMap := make(map[string]string)
	for k, v := range ShortVarMap() {
		revMap[v] = k

	}

	var query VariableCodes
	vars := map[string]interface{}{
		"ds": graphql.String(ds),
	}
	SendQueryVars(&query, vars)

	var metadata []Metadata

	for _, v := range query.Dataset.Variables.Edges {
		if revMap[string(v.Node.Name)] == "" {
			continue
		}

		var name string
		if nomis {
			name = revMap[string(v.Node.Name)]
		} else {
			name = string(v.Node.Name)

		}
		md := Metadata{
			Code: name,
			//Code: string(v.Node.Name),
			Name: string(v.Node.Label),
		}
		var query2 ClassCodes
		vars := map[string]interface{}{
			"ds":   graphql.String(ds),
			"vars": graphql.String(v.Node.Name),
		}
		SendQueryVars(&query2, vars)
		for _, v2 := range query2.Dataset.Table.Dimensions {
			for _, v3 := range v2.Categories { // XXX not ordered!
				md.Categories = append(md.Categories, struct {
					Code string
					Name string
				}{Code: string(v3.Code), Name: string(v3.Label)})

			}
		}

		metadata = append(metadata, md)

	}

	bs, err := json.Marshal(metadata)
	if err != nil {
		log.Print(err)
	}

	return (string(bs))
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

func ParseMetric(geo, cats Pairs, values IntValues) string {
	first := []string{"cantabular"}
	second := []string{"geography_code"}

	var b strings.Builder
	w := csv.NewWriter(&b)

	for _, k := range cats {
		first = append(first, string(k.Label))
		second = append(second, string(k.Code))
	}

	w.Write(first)
	w.Write(second)

	k := 0
	for _, g := range geo {
		geo := strings.Split(string(g.Code), "syn")[1]
		line := []string{geo}

		for j := 0; j < len(cats); j++ {
			line = append(line, fmt.Sprintf("%d", values[k]))
			k++
		}
		w.Write(line)
	}

	w.Flush()

	return b.String()
}

func GeoTypeMap() map[string]string {

	return map[string]string{
		"Country": "Country",
		"Region":  "Region",
		"LAD":     "LA",
		"MSOA":    "MSOA",
	}
}

func ShortVarMap() map[string]string {
	// "matching" via command output and eg.
	// SELECT  nd.name,nc.* FROM nomis_desc nd, nomis_category nc where nd.short_nomis_code='KS103EW' and nd.id=nc.nomis_desc_id and nc.measurement_unit='Count' and nc.long_nomis_code not like '%0001';

	// these are syn2011 "keys" which we pretend, temporarily, are NOMIS short codes
	return map[string]string{
		"KS102EW": "AGE_T009A",
		"KS103EW": "MARSTAT_T006A",
		"KS202EW": "NATID_ALL_T009A",
		"KS206EW": "WELSHPUK112_T007A",
		"KS207WA": "WELSHPUK112_R003A",
		"KS208WA": "WELSHPUK112_R003A",
		"QS101EW": "RESIDTYPE",
		"QS104EW": "SEX",
		"QS201EW": "ETHPUK11_T009A",
		"QS203EW": "COB_R010A",
		"QS208EW": "RELIGIONEW",
		"QS301EW": "CARER",
		"QS302EW": "HEALTH_T004A", // HEALTH
		"QS303EW": "DISABILITY_T003B",
		"QS402EW": "TYPACCOM_T009A",
		"QS403EW": "TENHUK11_T010A",
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
	}

}

//  GetDataSet maps variable code to dataset
func GetDataSet(varCode string) string {

	mappy := map[string]string{
		"QS403EW": "People-Households",
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
