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
	flag.Var(&geotypes, "geotype", "geography types (LSOA, LAD, etc)")
	flag.Var(&rows, "rows", "row or row range")
	flag.Var(&cols, "cols", "column name(s) to return")
	maxmetrics := flag.Int("maxmetrics", 0, "max skinny rows to accept (default 0 means no limit)")
	flag.Parse()

	if *dataset == "" {
		usage()
	}
	fmt.Printf("bbox: %s\n", *bbox)
	fmt.Printf("location: %s\n", *location)
	fmt.Printf("radius: %d meters\n", *radius)
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

	// Open postgres connection
	//
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
	)
	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	// Set up our geodata app
	app, err := geodata.New(db, *maxmetrics)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	body, err := app.Query(ctx, *dataset, *bbox, *location, *radius, geotypes, rows, cols)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s", body)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s --dataset <dataset> [--rows rowspec[,...]|--bbox p1lon,p1lat,p2lon,pl2lat|--location lon,lat --radius meters] [--geotype LSOA|LAD,...] [--cols col[,...]] [--maxmetrics n]\n", os.Args[0])
	os.Exit(2)
}
