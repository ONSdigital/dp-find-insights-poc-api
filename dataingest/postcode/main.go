package main

//https://geoportal.statistics.gov.uk/datasets/a8d42df48f374a52907fe7d4f804a662/about

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func readCsvFile(filePath string) (records [][]string) {
	func() {
		f, err := os.Open(filePath)
		if err != nil {
			log.Fatal("Unable to read input file "+filePath,
				err)
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		records, err = csvReader.ReadAll()

		if err != nil {
			log.Fatal("Unable to parse file as CSV for "+
				filePath, err)
		}
	}()

	return records
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dsn := database.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Print(err)
	}

	fmt.Println(dsn)

	t0 := time.Now()
	parsePostcodeCSV(db, "PCD_OA_LSOA_MSOA_LAD_MAY20_UK_LU.csv")
	fmt.Printf("%#v\n", time.Since(t0).Seconds())
}

func parsePostcodeCSV(db *gorm.DB, file string) {

	records := readCsvFile(file)

	field := make(map[string]int)
	for i, k := range records[0] {
		field[k] = i
	}

	j := 0
	for i, line := range records {
		if i == 0 {
			continue
		}

		/*
			pcd: 7-character version of the postcode (e.g. 'BT1 1AA', 'BT486PL')
			pcd2: 8-character version of the postcode (e.g. 'BT1  1AA', 'BT48 6PL')
			pcds: one space between the district and sector-unit part of the postcode (e.g. 'BT1 1AA', 'BT48 6PL') - possibly the most common formatting of postcodes.
		*/

		pcds := line[field["pcds"]]
		msoa11cd := line[field["msoa11cd"]]
		// Scotland isn't of interest nor Northern Ireland nor Channel Islands nor Isle of Man
		if !strings.HasPrefix(msoa11cd, "S") && !strings.HasPrefix(msoa11cd, "N") && !strings.HasPrefix(msoa11cd, "L") && !strings.HasPrefix(msoa11cd, "M") && msoa11cd != "" {

			var geos model.Geo
			db.Where("type_id = 5 and code=?", msoa11cd).Find(&geos) // limit by MSOA
			if geos.ID == 0 {
				log.Fatalf("not found: %s", msoa11cd)
			}

			var pc model.PostCode
			pc.GeoID = geos.ID
			pc.Pcds = pcds
			db.Save(&pc)

			j++
			if j%100000 == 0 {
				fmt.Printf("~%0f%% ...\n", (float64(j)/2300000)*100)
			}
		}

	}

	log.Println(j)
}
