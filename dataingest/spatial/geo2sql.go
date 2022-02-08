package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/ONSdigital/dp-find-insights-poc-api/dataingest/spatial/table"
	"github.com/twpayne/go-geom/encoding/geojson"
)

// couldn't find this as a go-geom type?
type FeatureCollection struct {
	Type     string            `json:"type"`
	Name     string            `json:"name"`
	CRS      geojson.CRS       `json:"crs"`
	Features []geojson.Feature `json:"features"`
}

func main() {
	filename := flag.String("f", "", "geojson file to import")
	tablename := flag.String("t", "", "name of table to create")
	mode := flag.String("m", "copy", "SQL copy or insert")
	keepcase := flag.Bool("k", false, "preserve case of property (column) names in features")
	flag.Parse()

	if *filename == "" || *tablename == "" {
		log.Fatalf("%s: filename and table required\n", os.Args[0])
	}
	if *mode != "copy" && *mode != "insert" {
		log.Fatalf("%s: mode must be copy or insert\n", os.Args[0])
	}

	collection, err := loadCollection(*filename)
	if err != nil {
		log.Fatal(err)
	}

	tab, err := table.New(*tablename, collection.Features, *keepcase, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	if err = tab.CreateTable(); err != nil {
		log.Fatal(err)
	}

	if *mode == "copy" {
		err = tab.CopyRows(collection.Features)
	} else {
		err = tab.InsertRows(collection.Features)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func loadCollection(name string) (*FeatureCollection, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var collection FeatureCollection
	decoder := json.NewDecoder(f)
	if err = decoder.Decode(&collection); err != nil {
		return nil, err
	}
	return &collection, nil
}
