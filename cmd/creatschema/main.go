package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
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

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
	)

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

	if err := db.Exec(`ALTER TABLE geo ADD COLUMN wkb_geometry geometry(Geometry,4326)`); err != nil {
		if !strings.Contains(err.Error.Error(), "SQLSTATE 42701") {
			log.Print(err.Error.Error())
		}
	}

	if haveDump {
		ndump, _ := pgDump()
		f, err := os.Create("sql/schema.sql") // XXX
		if err != nil {
			log.Print(err)
		}
		f.WriteString(ndump)

		if ndump != odump {
			bs, _ := exec.Command("git", "diff", "sql/schema.sql").Output()
			//bs, _ := exec.Command("diff", tf.Name(), "sql/schema.sql").Output()
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
