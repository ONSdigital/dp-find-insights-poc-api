package metadata

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
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
		var gdb *gorm.DB
		gdb, err = gorm.Open(postgres.Open(database.GetDSN()), &gorm.Config{
			//	Logger: logger.Default.LogMode(logger.Info), // display SQL
		})

		return &Metadata{gdb: gdb}, err
	}

	return &Metadata{gdb: dbs[0]}, err
}

func (md *Metadata) Get(ctx context.Context, year int, filterTotals bool) ([]byte, error) {
	var topics []model.NomisTopic

	md.gdb.Preload(
		"NomisDescs",
		func(gdb *gorm.DB) *gorm.DB {
			return gdb.Order("short_nomis_code").Where("year = ?", year)
		},
	).Find(&topics)

	var mdr api.MetadataResponse

	for _, topic := range topics {

		// this is a dummy row used to import
		if topic.ID == 0 {
			continue
		}

		var newTabs api.Tables
		var nd model.NomisDesc

		for _, nd = range topic.NomisDescs {
			md.gdb.Preload(
				"NomisCategories",
				func(gdb *gorm.DB) *gorm.DB {
					return gdb.Order("long_nomis_code").Where("year = ?", year)
				},
			).Find(&nd)

			// partially populate table here to allow optional inclusion of Total if filterTotals == true
			table := api.Table{
				Name: spointer(nd.Name),
				Slug: spointer(slug.Make(nd.Name)),
				Code: spointer(nd.ShortNomisCode),
			}

			var cats api.Categories
			for _, trip := range nd.NomisCategories {
				cat := api.Triplet{Code: spointer(trip.LongNomisCode), Name: spointer(trip.CategoryName), Slug: spointer(slug.Make(trip.CategoryName))}
				if filterTotals && isTotalCat(trip.LongNomisCode) {
					table.Total = &cat
				} else {
					cats = append(cats, cat)
				}
			}
			table.Categories = &cats
			newTabs = append(newTabs, table)
		}

		mdr = append(mdr, api.Metadata{
			Code:   spointer(topic.TopNomisCode),
			Name:   spointer(topic.Name),
			Slug:   spointer(slug.Make(topic.Name)),
			Tables: &newTabs,
		})

	}

	b, err := json.Marshal(&mdr)
	if err != nil {
		log.Error(ctx, "json marshal", err)
		return b, err
	}

	return b, nil
}

func spointer(s string) *string {
	return &s
}

func (md *Metadata) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	db, err := md.gdb.DB() // extract gorm's underlying db handle
	if err != nil {
		state.Update(healthcheck.StatusCritical, err.Error(), 0)
		return nil
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		state.Update(healthcheck.StatusCritical, err.Error(), 0)
		return nil
	}
	state.Update(healthcheck.StatusOK, "gorm healthy", 0)
	conn.Close()
	return nil
}

func isTotalCat(catName string) bool {
	return strings.HasSuffix(catName, "0001")
}
