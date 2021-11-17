package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/demo"
)

func main() {
	var rows, cols multiFlag

	dataset := flag.String("dataset", "", "name of dataset to query")
	flag.Var(&rows, "rows", "row or row range")
	flag.Var(&cols, "cols", "column name(s) to return")
	flag.Parse()

	if *dataset == "" {
		usage()
	}
	fmt.Printf("rows:\n")
	for _, r := range rows {
		fmt.Printf("\t%s\n", r)
	}

	fmt.Printf("cols:\n")
	for _, c := range cols {
		fmt.Printf("\t%s\n", c)
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

	// Set up our demo app
	app, err := demo.New(db)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	body, err := app.Query(ctx, *dataset, rows, cols)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s", body)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s --dataset <dataset> --rows rowspec[,...] --cols col[,...]\n", os.Args[0])
	os.Exit(2)
}
