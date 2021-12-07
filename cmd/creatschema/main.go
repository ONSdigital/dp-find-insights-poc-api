package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running (can include -dirty)
	Version string
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if BuildTime == "" || GitCommit == "" {
		log.Fatal("run from Makefile target")
	}

	var line string
	if strings.Contains(Version, "dirty") {
		fmt.Print("Enter 'y' to confirm deploy of unchecked in changes? ")
		fmt.Scanln(&line)
		if line != "y" {
			log.Fatal("exiting")
		}
	}

	dsn := database.GetDSN()

	fmt.Printf("using dsn: '%s' continue (y/n)? ", dsn)

	fmt.Scanln(&line)
	if line != "y" {
		log.Fatal("exiting")
	}

	fmt.Println("migrating DB")

	odump, haveDump := pgDump()

	var tf *os.File
	if haveDump {
		var err error
		tf, err = os.CreateTemp("/tmp", "*.sql")
		if err != nil {
			panic(err)
		}
		tf.WriteString(odump)
	} else {
		log.Print("'pg_dump' not detected in PATH not doing schema dumps")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Print(err)
	}

	// XXX create/alter tables - doesn't delete cols or tables!
	// neither does it always change types correctly
	// More useful in dev than prod

	if err := db.AutoMigrate(
		&model.SchemaVer{},
		&model.DataVer{},
		&model.GeoType{},
		&model.Geo{},
		&model.NomisDesc{},
		&model.NomisCategory{},
		&model.GeoMetric{},
		&model.YearMapping{},
	); err != nil {
		log.Print(err)
	}

	// refactor
	if err := db.Exec(`ALTER TABLE geo ADD COLUMN wkb_geometry geometry(Geometry,4326)`).Error; err != nil {
		log.Print(err)
	}

	if err := db.Exec(`CREATE INDEX geo_wkb_geometry_geom_idx ON public.geo USING gist (wkb_geometry);`).Error; err != nil {
		log.Print(err)
	}

	// refactor XXX

	if err := db.Exec(`ALTER TABLE geo ADD COLUMN wkb_long_lat_geom geometry(Geometry,4326)`).Error; err != nil {
		log.Print(err)
	}

	if err := db.Exec(`CREATE INDEX geo_long_lat_geom_idx ON public.geo USING gist ( wkb_long_lat_geom);`).Error; err != nil {
		log.Print(err)
	}

	if haveDump {
		ndump, _ := pgDump()
		f, err := os.Create("sql/schema.sql")
		if err != nil {
			log.Print(err)
		}
		f.WriteString(ndump)

		if ndump != odump {
			bs, _ := exec.Command("git", "diff", "sql/schema.sql").Output()
			fmt.Println(string(bs))
			fmt.Println("check-in sql/schema.sql")
		} else {
			fmt.Println("no schema changes")
			os.Exit(0)
		}
	}

	db.Save(&model.SchemaVer{BuildTime: BuildTime, GitCommit: GitCommit, Version: Version})

	// populate data_ver
	db.Save(&model.DataVer{ID: 1, CensusYear: 2011, VerString: "2.2", Public: true, Source: "Nomis Bulk API", Notes: "Release date 12/02/2013 Revised 17/01/2014"})

}

func pgDump() (string, bool) {
	bs, err := exec.Command("pg_dump", "--schema-only").Output()
	if err == nil {
		return string(bs), true
	}

	return "", false
}
