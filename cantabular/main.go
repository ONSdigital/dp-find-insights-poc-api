package main

// adhoc query tool to investigate cantabular 2011 instance

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

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

// like Nomis short codes and geo_type names mixed up (?)
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

// return metrics
type ClassCodes struct {
	Dataset struct {
		Table struct {
			Dimensions []struct {
				Categories []struct {
					Code  graphql.String
					Label graphql.String
				}
			}
		} `graphql:"table(variables: [$vars])"` // XXX
	} `graphql:"dataset(name: $ds)"`
}

func sendQueryVars(query interface{}, vars map[string]interface{}) interface{} {
	if os.Getenv("CANT_USER") == "" || os.Getenv("CANT_PW") == "" {
		log.Fatal("define CANT_USER & CANT_PW")
	}
	hclient := &http.Client{Transport: AuthTripper{User: os.Getenv("CANT_USER"), Pass: os.Getenv("CANT_PW")}}
	client := graphql.NewClient(URL, hclient)
	if err := client.Query(context.Background(), query, vars); err != nil {
		log.Print(err)
	}
	return query
}

func parseResp(query interface{}) {
	// wish there were a better way!
	qt := reflect.TypeOf(query).String()
	switch qt {
	case "*main.DataSets":
		ds := query.(*DataSets)
		for _, v := range ds.Datasets {
			fmt.Println(v.Name)
		}
	case "*main.VariableCodes":
		scodes := query.(*VariableCodes)
		for _, v := range scodes.Dataset.Variables.Edges {
			fmt.Print(v.Node.Name + " : ")
			fmt.Println(v.Node.Label)
		}
	case "*main.ClassCodes":
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

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	datasets := flag.Bool("datasets", false, "list datasets, eg. Usual-Residents")
	ds := flag.String("ds", "Usual-Residents", "set dataset to query")
	class := flag.String("class", "", "classifications under variables eg. pass AGE_T022A (or MSOA) to get categories under it (like old longcodes)")
	variables := flag.Bool("variables", false, "list variables, results eg. 'AGE_T022A : Age of individual (21 categories)' (like old short codes)")
	flag.Parse()

	if *datasets {
		var query DataSets
		sendQueryVars(&query, nil)
		parseResp(&query)
		os.Exit(0)
	}

	if *variables {
		var query VariableCodes
		vars := map[string]interface{}{
			"ds": graphql.String(*ds),
		}
		sendQueryVars(&query, vars)
		parseResp(&query)
		fmt.Println("\nUSED: '" + *ds + "'")
		os.Exit(0)
	}

	if len(*class) > 0 {
		var query ClassCodes
		vars := map[string]interface{}{
			"ds":   graphql.String(*ds),
			"vars": graphql.String(*class),
		}
		sendQueryVars(&query, vars)
		parseResp(&query)
		fmt.Println("\nUSED: '" + *ds + "'")

		os.Exit(0)
	}

	flag.PrintDefaults()

}
