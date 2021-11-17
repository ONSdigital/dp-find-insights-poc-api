package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Print(err)
	}

	// XXX create/alter tables - doesn't delete cols or tables!
	// neither does it always change types correctly
	// More useful in dev than prod

	if err := db.AutoMigrate(
		&model.GeoType{},
		&model.Geo{},
		&model.NomisDesc{},
		&model.NomisCategory{},
		&model.GeoMetric{},
		&model.YearMapping{},
	); err != nil {
		log.Print(err)
	}
}
