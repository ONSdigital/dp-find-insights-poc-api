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

	model.Migrate(db)

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
		}
	}

	db.Save(&model.SchemaVer{BuildTime: BuildTime, GitCommit: GitCommit, Version: Version})

	// REFACTOR into something like SetupData
	// populate data_ver
	db.Save(&model.DataVer{ID: 1, CensusYear: 2011, VerString: "2.2", Public: true, Source: "Nomis Bulk API", Notes: "Release date 12/02/2013 Revised 17/01/2014"})
	db.Save(&model.NomisTopic{TopNomisCode: "QS1", Name: "Population Basics"})
	db.Save(&model.NomisTopic{TopNomisCode: "QS2", Name: "Origins & Beliefs"})
	db.Save(&model.NomisTopic{TopNomisCode: "QS3", Name: "Health"})
	db.Save(&model.NomisTopic{TopNomisCode: "QS4", Name: "Housing"})
	db.Save(&model.NomisTopic{TopNomisCode: "QS5", Name: "Education"})
	db.Save(&model.NomisTopic{TopNomisCode: "QS6", Name: "Employment"})
	db.Save(&model.NomisTopic{TopNomisCode: "QS7", Name: "Travel to Work"})
	db.Save(&model.NomisTopic{TopNomisCode: "QS8", Name: "Residency"})

}

func pgDump() (string, bool) {
	bs, err := exec.Command("pg_dump", "--schema-only").Output()
	if err == nil {
		return string(bs), true
	}

	return "", false
}
