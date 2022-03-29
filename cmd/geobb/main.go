package main

import (
	"fmt"
	"log"

	"github.com/ONSdigital/dp-find-insights-poc-api/geobb"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	gdb, err := gorm.Open(postgres.Open(database.GetDSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Print(err)
	}

	g := geobb.GeoBB{Gdb: gdb}

	fmt.Printf("%s\n", g.AsJSON(geobb.Params{Welsh: false, Pretty: false, Geos: []string{"LAD", "MSOA"}}))

}
