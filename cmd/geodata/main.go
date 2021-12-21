package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	geodata "github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
)

func main() {
	var rows, cols, geotypes multiFlag

	dataset := flag.String("dataset", "", "name of dataset to query")
	bbox := flag.String("bbox", "", "bounding box lon1,lat1,lon2,lat2 (any two opposite corners)")
	location := flag.String("location", "", "central point for radius queries")
	radius := flag.Int("radius", 0, "radius in meters")
	polygon := flag.String("polygon", "", "polygon x1,y1,...,x1,y1 (closed linestring)")
	flag.Var(&geotypes, "geotype", "geography types (LSOA, LAD, etc)")
	flag.Var(&rows, "rows", "row or row range")
	flag.Var(&cols, "cols", "column name(s) to return")
	maxmetrics := flag.Int("maxmetrics", 0, "max skinny rows to accept (default 0 means no limit)")
	censustable := flag.String("censustable", "", "censustable QS802EW 'nomis table' / grouping of census data categories")
	flag.Parse()

	if *dataset == "" {
		usage()
	}
	fmt.Printf("bbox: %s\n", *bbox)
	fmt.Printf("location: %s\n", *location)
	fmt.Printf("radius: %d meters\n", *radius)
	fmt.Printf("polygon: %s\n", *polygon)
	fmt.Printf("rows:\n")
	for _, r := range rows {
		fmt.Printf("\t%s\n", r)
	}

	fmt.Printf("cols:\n")
	for _, c := range cols {
		fmt.Printf("\t%s\n", c)
	}

	fmt.Printf("geotypes:\n")
	for _, t := range geotypes {
		fmt.Printf("\t%s\n", t)
	}

	fmt.Printf("censustable: %s\n", *censustable)
	// Open postgres connection
	//

	db, err := database.Open("pgx", database.GetDSN())
	if err != nil {
		log.Fatalln(err)
	}

	// Set up our geodata app
	app, err := geodata.New(db, *maxmetrics)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	body, err := app.Query(ctx, *dataset, *bbox, *location, *radius, *polygon, geotypes, rows, cols, *censustable)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s", body)
}

func usage() {
	fmt.Fprintf(
		os.Stderr,
		"usage: %s --dataset <dataset> [--rows rowspec[,...]|--bbox p1lon,p1lat,p2lon,pl2lat|--location lon,lat --radius meters|--polygon x1,y1,...,x1,y1] [--geotype LSOA|LAD,...] [--cols col[,...]] [--maxmetrics n] [--censustable QS802EW]\n",
		os.Args[0],
	)
	os.Exit(2)
}
