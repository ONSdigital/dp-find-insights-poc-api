package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/jackc/pgx/v4"
	"github.com/jszwec/csvutil"
	"github.com/spf13/cast"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// aws --region eu-central-1 s3 sync s3://find-insights-input-data-files/nomis/ .

const dataPref = "dataingest/addtodb/data/"

type dataIngest struct {
	gdb     *gorm.DB
	conn    *pgx.Conn
	dataVer string
	files   files
}

type files struct {
	meta []string
	data []string
	desc []string
}

// New takes optimal optimal dsn arg for testing override
func New(v string, dsns ...string) dataIngest {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var dsn string
	if len(dsns) == 0 {
		dsn = database.GetDSN()
	} else {
		dsn = dsns[0]
	}

	log.Printf("dsn=%s\n", dsn)

	// be nice to share same handle but I can't see how to do this!

	con, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}

	dbg, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal(err)
	}

	return dataIngest{gdb: dbg, conn: con, dataVer: v}
}

func (di *dataIngest) addGeoTypes() {

	if tx := di.gdb.Save(&model.GeoType{ID: 1, Name: "EW"}); tx.Error != nil {
		log.Fatal(tx.Error)
	}

	di.gdb.Save(&model.GeoType{ID: 2, Name: "Country"})
	di.gdb.Save(&model.GeoType{ID: 3, Name: "Region"})
	di.gdb.Save(&model.GeoType{ID: 4, Name: "LAD"})
	di.gdb.Save(&model.GeoType{ID: 5, Name: "MSOA"})
	di.gdb.Save(&model.GeoType{ID: 6, Name: "LSOA"})
}

func (di *dataIngest) getFiles(pref string) {

	_, err := os.Stat(pref)
	if os.IsNotExist(err) {
		log.Fatal(err)
	}

	err = filepath.Walk(pref, func(path string, info os.FileInfo, err error) error {

		if strings.Contains(info.Name(), "META") {
			di.files.meta = append(di.files.meta, path)
		}

		if strings.Contains(info.Name(), "DESC0") {
			di.files.desc = append(di.files.desc, path)
		}

		if strings.Contains(info.Name(), "DATA0") {
			di.files.data = append(di.files.data, path)
		}

		return err

	})

	if err != nil {
		log.Fatal(err)
	}
}

func (di *dataIngest) addCategoryData() map[string]int32 {

	m := make(map[string]int32)
	for _, fn := range di.files.desc {
		f, err := os.Open(fn)
		if err != nil {
			log.Print(err)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Print(err)
		}
		type DiscTable struct {
			ColumnVariableCode            string
			ColumnVariableMeasurementUnit string
			ColumnVariableStatisticalUnit string
			ColumnVariableDescription     string
		}

		var discTables []DiscTable
		if err := csvutil.Unmarshal(b, &discTables); err != nil {
			fmt.Println("error:", err)
		}

		y := cast.ToInt32(di.dataVer)
		for _, dt := range discTables {

			longNomisCode := dt.ColumnVariableCode
			shortNomisCode := longNomisCode[0:7]

			var desc model.NomisDesc
			di.gdb.Where("short_nomis_code = ?", shortNomisCode).First(&desc)

			nc := model.NomisCategory{
				NomisDescID:     desc.ID,
				CategoryName:    dt.ColumnVariableDescription,
				MeasurementUnit: dt.ColumnVariableMeasurementUnit,
				StatUnit:        dt.ColumnVariableStatisticalUnit,
				LongNomisCode:   longNomisCode,
				Year:            y,
			}

			if tx := di.gdb.Save(&nc); tx.Error != nil {
				log.Fatal(err)
			}

			m[longNomisCode] = nc.ID

		}
	}

	return m
}

// adds to "geo" tables (but without geo.name!)
// adds to "geo_metric"
func (di *dataIngest) addGeoGeoMetricData(longToCatid map[string]int32) {
	geoCodeToID := make(map[string]int64)

	con := di.conn

	con.Exec(context.Background(), "SET synchronous_commit TO off")
	re := regexp.MustCompile(`DATA0(\d)\.CSV`)

	ctx := context.Background()

	// pgx is supposed to auto prepare but it's easy enough...
	_, err := con.Prepare(ctx, "geo-insert", "INSERT INTO geo (code,name,type_id) VALUES ($1,$2,$3) RETURNING id")
	if err != nil {
		log.Print(err)
	}

	num := len(di.files.data)

	t0 := time.Now()

	for i, fn := range di.files.data {

		match := re.FindStringSubmatch(fn)

		if len(match) == 0 {
			continue
		}

		fmt.Printf("file %d of %d, %.2f step min(s), name=%s\n", i, num, time.Since(t0).Minutes(), fn)

		geoType := cast.ToInt32(match[1])

		// XXX we do need MSOA see Trello #282
		if geoType == 5 {
			log.Println("skipping MSOA...")
			continue
		}

		recs := readCsvFile(fn)

		headers := recs[0]

		header := make(map[int]string)

		for k, v := range headers {
			header[k] = v
		}

		num := len(recs)

		fmt.Printf("processing: %#v recs\n", num)

		var rows [][]interface{}

		// lines in file "E01000001,1465,50..."
		for i, row := range recs {

			// skip headers
			if i == 0 {
				continue
			}

			var geoID int64

			// columns over row "geographyCode,KS102EW0001,..."
			for j, col := range row {

				key := header[j]

				if key == "GeographyCode" {

					if geoCodeToID[col] == 0 {

						err := con.QueryRow(ctx, "geo-insert", col, "NA", geoType).Scan(&geoID)
						if err != nil {
							log.Fatal(err)
						}

						geoCodeToID[col] = geoID
					} else {
						geoID = geoCodeToID[col]

					}
					continue
				}

				if geoID > 0 {
					rows = append(rows, []interface{}{1, geoID, longToCatid[key], cast.ToFloat64(col)})
				}

			}

		} // end rows

		fmt.Println("Bulk copy...")

		var count int64
		count, err = con.CopyFrom(ctx,
			pgx.Identifier{"geo_metric"},
			[]string{"data_ver_id", "geo_id", "category_id", "metric"},
			pgx.CopyFromRows(rows),
		)

		fmt.Printf("count: %#v\n", count)

		if err != nil {
			log.Print(err)
		}

	} // end files
}

// TODO v4 rename Classification
func (di *dataIngest) addClassificationData() {

	for _, f := range di.files.meta {

		recs := readCsvFile(f)

		m := make(map[string]string)
		for i, v := range recs[0] {
			m[v] = recs[1][i]
		}

		// skip some duff data in Nomis Bulk 2011
		if m["DatasetTitle"] != "Cyfradd" && m["DatasetTitle"] != "Pellter teithio i'r gwaith " && m["DatasetTitle"] != "" && di.dataVer == "2011" {

			di.gdb.Save(&model.NomisDesc{
				Name:           m["DatasetTitle"],
				PopStat:        m["StatisticalPopulations"],
				ShortNomisCode: m["DatasetId"],
				Year:           2011,
			})
		}

	}
}

func readCsvFile(filePath string) (records [][]string) {

	func() {
		f, err := os.Open(filePath)
		if err != nil {
			log.Fatal("Unable to read input file "+filePath, err)
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		records, err = csvReader.ReadAll()

		if err != nil {
			log.Fatal("Unable to parse file as CSV for "+filePath, err)
		}
	}()

	return records
}

// probably means this can be removed from model/provision
// This should be refactored into a config file

func (di *dataIngest) popTopics() {
	di.gdb.Save(&model.NomisTopic{ID: 1, TopNomisCode: "QS1", Name: "Population Basics"})
	di.gdb.Save(&model.NomisTopic{ID: 2, TopNomisCode: "QS2", Name: "Origins & Beliefs"})
	di.gdb.Save(&model.NomisTopic{ID: 3, TopNomisCode: "QS3", Name: "Health"})
	di.gdb.Save(&model.NomisTopic{ID: 4, TopNomisCode: "QS4", Name: "Housing"})
	di.gdb.Save(&model.NomisTopic{ID: 5, TopNomisCode: "QS5", Name: "Education"})
	di.gdb.Save(&model.NomisTopic{ID: 6, TopNomisCode: "QS6", Name: "Employment"})
	di.gdb.Save(&model.NomisTopic{ID: 7, TopNomisCode: "QS7", Name: "Travel to Work"})
	di.gdb.Save(&model.NomisTopic{ID: 8, TopNomisCode: "QS8", Name: "Residency"})
	di.gdb.Save(&model.NomisTopic{ID: 400, TopNomisCode: "KS4", Name: "Housing"})

	di.gdb.Save(&model.NomisTopic{ID: 100, TopNomisCode: "KS1", Name: "Population Basics"})
	di.gdb.Save(&model.NomisTopic{ID: 200, TopNomisCode: "KS2", Name: "Origins & Beliefs"})
	di.gdb.Save(&model.NomisTopic{ID: 400, TopNomisCode: "KS4", Name: "Housing"})
}

func (di *dataIngest) putVersion() {
	di.gdb.Save(&model.DataVer{ID: 1, CensusYear: 2011, VerString: "2.2", Public: true, Source: "Nomis Bulk API", Notes: "20220204 2i using go addtodb"})
}

func main() {
	t0 := time.Now()

	di := New("2011") // TODO get from command line
	di.getFiles(dataPref)
	di.popTopics()
	di.addGeoTypes()
	di.addClassificationData()
	longToCatid := di.addCategoryData()
	di.addGeoGeoMetricData(longToCatid)
	di.putVersion()

	fmt.Printf("%#v\n", time.Since(t0).Seconds())
}
