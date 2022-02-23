package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ONSdigital/dp-find-insights-poc-api/metadata"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	geodata "github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	maxmetrics := flag.Int("maxmetrics", 0, "max number of rows to accept from db query (default 0 means no limit)")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [command-options] query|ckmeans|ckmeansratio|metadata [subcommand-options]\n", filepath.Base(os.Args[0]))
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

	gdb, err := gorm.Open(postgres.Open(database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	md, err := metadata.New(gdb)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	switch flag.Arg(0) {
	case "query":
		query(ctx, app, flag.Args()[1:])
	case "ckmeans":
		ckmeans(ctx, app, flag.Args()[1:])
	case "metadata":
		mdquery(ctx, md, flag.Args()[1:])
	default:
		flag.Usage()
		os.Exit(2)
	}
}

func query(ctx context.Context, app *geodata.Geodata, argv []string) {
	var rows, cols, geotypes multiFlag

	flagset := flag.NewFlagSet("original", flag.ExitOnError)

	year := flagset.Int("year", 2011, "census year")
	bbox := flagset.String("bbox", "", "bounding box lon1,lat1,lon2,lat2 (any two opposite corners)")
	location := flagset.String("location", "", "central point for radius queries")
	radius := flagset.Int("radius", 0, "radius in meters")
	polygon := flagset.String("polygon", "", "polygon x1,y1,...,x1,y1 (closed linestring)")
	censustable := flagset.String("censustable", "", "censustable QS802EW 'nomis table' / grouping of census data categories")
	flagset.Var(&geotypes, "geotype", "geography types (LSOA, LAD, etc)")
	flagset.Var(&rows, "rows", "row or row range")
	flagset.Var(&cols, "cols", "column name(s) to return")
	flagset.Parse(argv)

	body, err := app.Query(ctx, *year, *bbox, *location, *radius, *polygon, geotypes, rows, cols, *censustable)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s", body)
}

func ckmeans(ctx context.Context, app *geodata.Geodata, argv []string) {
	var cat, geotype multiFlag

	flagset := flag.NewFlagSet("ckmeans", flag.ExitOnError)

	year := flagset.Int("year", 2011, "census year")
	flagset.Var(&cat, "cat", "category code(s) to provide ckmeans for")
	flagset.Var(&geotype, "geotype", "geography types (LSOA, LAD, etc)")
	k := flagset.Int("k", 5, "number of clusters/bins")
	divide_by := flagset.String("divide_by", "", "category code to divide all other categories by (optional)")
	flagset.Parse(argv)

	breaks, err := app.CKmeans(ctx, *year, cat, geotype, *k, *divide_by)
	if err != nil {
		log.Fatalln(err)
	}
	buf, err := json.MarshalIndent(breaks, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(string(append(buf, "\n"...)))
}

func mdquery(ctx context.Context, md *metadata.Metadata, argv []string) {
	flagset := flag.NewFlagSet("metadata", flag.ExitOnError)

	year := flagset.Int("year", 2011, "census year")
	filtertotals := flagset.Bool("filtertotals", false, "include totals")
	flagset.Parse(argv)

	result, err := md.Get(*year, *filtertotals)
	if err != nil {
		log.Fatalln(err)
	}
	if _, err = os.Stdout.Write(result); err != nil {
		log.Fatalln(err)
	}
}
