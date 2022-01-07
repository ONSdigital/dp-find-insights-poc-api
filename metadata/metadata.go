package metadata

import (
	"encoding/json"
	"log"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/gosimple/slug"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Metadata struct {
	db *gorm.DB
}

func New() (*Metadata, error) {
	db, err := gorm.Open(postgres.Open(database.GetDSN()), &gorm.Config{
		//	Logger: logger.Default.LogMode(logger.Info), // display SQL
	})

	return &Metadata{db: db}, err
}

func (md *Metadata) Get() (b []byte, err error) {
	var topics []model.NomisTopic

	md.db.Preload("NomisDescs", func(db *gorm.DB) *gorm.DB { return db.Order("short_nomis_code") }).Find(&topics)

	var mdr api.MetadataResponse

	var newTabs api.Tables

	for _, topic := range topics {

		var nd model.NomisDesc

		for _, nd = range topic.NomisDescs {
			md.db.Preload("NomisCategories").Find(&nd)

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
		log.Print(err)
	}

	return b, nil
}

func spointer(s string) *string {
	return &s
}
