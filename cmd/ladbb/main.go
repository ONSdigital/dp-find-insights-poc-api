package main

import (
	"fmt"
	"log"

	"github.com/ONSdigital/dp-find-insights-poc-api/ladbb"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	gdb, err := gorm.Open(postgres.Open(database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Print(err)
	}

	g := ladbb.LadGeom{Gdb: gdb}

	fmt.Printf("%s\n", g.AsJSON())

}
