package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	geodata "github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
)

func main() {
	maxmetrics := flag.Int("maxmetrics", 0, "max number of rows to accept from db query (default 0 means no limit)")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [command-options] original|ckmeans [subcommand-options]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	db, err := database.Open("pgx", database.GetDSN())
	if err != nil {
		log.Fatalln(err)
	}

	app, err := geodata.New(db, *maxmetrics)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	switch flag.Arg(0) {
	case "original":
		original(ctx, app, flag.Args()[1:])
	case "ckmeans":
		ckmeans(ctx, app, flag.Args()[1:])
	default:
		flag.Usage()
		os.Exit(2)
	}
}

func original(ctx context.Context, app *geodata.Geodata, argv []string) {
	var rows, cols, geotypes multiFlag

	flagset := flag.NewFlagSet("original", flag.ExitOnError)

	bbox := flagset.String("bbox", "", "bounding box lon1,lat1,lon2,lat2 (any two opposite corners)")
	location := flagset.String("location", "", "central point for radius queries")
	radius := flagset.Int("radius", 0, "radius in meters")
	polygon := flagset.String("polygon", "", "polygon x1,y1,...,x1,y1 (closed linestring)")
	censustable := flagset.String("censustable", "", "censustable QS802EW 'nomis table' / grouping of census data categories")
	flagset.Var(&geotypes, "geotype", "geography types (LSOA, LAD, etc)")
	flagset.Var(&rows, "rows", "row or row range")
	flagset.Var(&cols, "cols", "column name(s) to return")
	flagset.Parse(argv)

	body, err := app.Query(ctx, "census", *bbox, *location, *radius, *polygon, geotypes, rows, cols, *censustable)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s", body)
}

func ckmeans(ctx context.Context, app *geodata.Geodata, argv []string) {
	flagset := flag.NewFlagSet("ckmeans", flag.ExitOnError)

	cat := flagset.String("cat", "", "category code")
	geotype := flagset.String("geotype", "", "geography type (LSOA,...)")
	k := flagset.Int("k", 5, "number of clusters/bins")
	flagset.Parse(argv)

	breaks, err := app.CKmeans(ctx, *cat, *geotype, *k)
	if err != nil {
		log.Fatalln(err)
	}
	for _, breakpoint := range breaks {
		fmt.Printf("%0.13g\n", breakpoint)
	}
}
