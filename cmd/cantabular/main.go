package main

// adhoc query tool to investigate cantabular 2011 instance

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ONSdigital/dp-find-insights-poc-api/cantabular"
	"github.com/shurcooL/graphql"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// check ds use
	query1 := flag.Bool("query1", false, "query uses -code, -geo and -geotype, eg. '-query -code QS501EW -geo E92000001 -geotype Country'")
	query2 := flag.Bool("query2", false, "query uses -geotype and -code, eg. '-query2 -code QS501EW -geotype Region'")
	code := flag.String("code", "", "query: code , eg. QS501EW")
	geo := flag.String("geo", "", "query: geo, eg. E92000001")
	geotype := flag.String("geotype", "", "query: geotype, eg. Country,Region,LAD,MSOA")

	datasets := flag.Bool("datasets", false, "list datasets, eg. Usual-Residents")
	ds := flag.String("ds", "Usual-Residents", "set dataset to query")
	class := flag.String("class", "", "classifications under variables eg. pass AGE_T022A (or MSOA) to get categories under it (like old longcodes)")
	variables := flag.Bool("variables", false, "list variables, results eg. 'AGE_T022A : Age of individual (21 categories)' (like old short codes)")
	flag.Parse()

	// MetricFilter type query
	if *query1 {

		checkParams(*code, *geotype)
		if *code == "" || *geo == "" || *geotype == "" {
			fmt.Println("must define -code, -geo and -geotype")
			os.Exit(1)
		}

		geoq, catsq, values := cantabular.QueryMetricFilter(*ds, *geo, *geotype, *code)
		got := cantabular.ParseMetric(geoq, catsq, values)

		fmt.Println(got)
		os.Exit(0)

	}

	// Pure Metric query
	if *query2 {
		checkParams(*code, *geotype)

		if *code == "" || *geotype == "" {
			fmt.Println("must define -code and -geotype")
			os.Exit(1)
		}

		geoq, catsq, values := cantabular.QueryMetric(*ds, *geotype, *code)
		got := cantabular.ParseMetric(geoq, catsq, values)

		fmt.Println(got)
		os.Exit(0)
	}

	if *datasets {
		var query cantabular.DataSets
		cantabular.SendQueryVars(&query, nil)
		cantabular.ParseResp(&query)
		os.Exit(0)
	}

	if *variables {
		var query cantabular.VariableCodes
		vars := map[string]interface{}{
			"ds": graphql.String(*ds),
		}
		cantabular.SendQueryVars(&query, vars)
		cantabular.ParseResp(&query)
		fmt.Println("\nUSED: '" + *ds + "'")
		os.Exit(0)
	}

	if len(*class) > 0 {
		var query cantabular.ClassCodes
		vars := map[string]interface{}{
			"ds":   graphql.String(*ds),
			"vars": graphql.String(*class),
		}
		cantabular.SendQueryVars(&query, vars)
		cantabular.ParseResp(&query)
		fmt.Println("\nUSED: '" + *ds + "'")

		os.Exit(0)
	}

	flag.PrintDefaults()

}

func checkParams(code, geoType string) {
	shorts := cantabular.ShortVarMap()
	if shorts[code] == "" {
		fmt.Println("err: use -code from following list")
		for k := range shorts {
			fmt.Print(k + " ")

		}
		os.Exit(1)
	}

	gts := cantabular.GeoTypeMap()
	if gts[geoType] == "" {
		fmt.Println("err: use -geotype from following list")
		for k := range gts {
			fmt.Print(k + " ")

		}
		os.Exit(1)
	}

}
