package metadata

import (
	"context"
	"encoding/json"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gosimple/slug"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Metadata struct {
	gdb *gorm.DB
}

// New takes optimal optimal db args for testing override
func New(dbs ...*gorm.DB) (*Metadata, error) {
	var err error
	if len(dbs) == 0 {

		// TODO this func should accept a persistent *sql.DB from
		// handler/hander.go and make gdb from that eg.
		// gorm.Open(postgres.New(postgres.Config{Conn: db.DB()}))

		var gdb *gorm.DB
		gdb, err = gorm.Open(postgres.Open(database.GetDSN()), &gorm.Config{
			//	Logger: logger.Default.LogMode(logger.Info), // display SQL
		})

		return &Metadata{gdb: gdb}, err
	}

	return &Metadata{gdb: dbs[0]}, err
}

func (md *Metadata) Get() (b []byte, err error) {
	var topics []model.NomisTopic

	md.gdb.Preload("NomisDescs", func(gdb *gorm.DB) *gorm.DB { return gdb.Order("short_nomis_code") }).Find(&topics)

	var mdr api.MetadataResponse

	for _, topic := range topics {

		// this is a dummy row used to import
		if topic.ID == 0 {
			continue
		}

		var newTabs api.Tables
		var nd model.NomisDesc

		for _, nd = range topic.NomisDescs {
			md.gdb.Preload("NomisCategories", func(gdb *gorm.DB) *gorm.DB { return gdb.Order("long_nomis_code") }).Find(&nd)

			var cats api.Categories
			for _, trip := range nd.NomisCategories {
				cats = append(cats, api.Triplet{Code: spointer(trip.LongNomisCode), Name: spointer(trip.CategoryName), Slug: spointer(slug.Make(trip.CategoryName))})
			}

			newTabs = append(newTabs,
				api.Table{
					Name:       spointer(nd.Name),
					Slug:       spointer(slug.Make(nd.Name)),
					Code:       spointer(nd.ShortNomisCode),
					Categories: &cats,
				})
		}

		mdr = append(mdr, api.Metadata{
			Code:   spointer(topic.TopNomisCode),
			Name:   spointer(topic.Name),
			Slug:   spointer(slug.Make(topic.Name)),
			Tables: &newTabs,
		})

	}

	b, err = json.Marshal(&mdr)

	if err != nil {
		log.Error(context.Background(), err.Error(), err)
		return b, err
	}

	return b, nil
}

func spointer(s string) *string {
	return &s
}
