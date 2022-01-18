// +build comptest

package geodata

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/comptests"
	"github.com/ONSdigital/dp-find-insights-poc-api/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// passing -args -kill=true to the test will kill docker postgres
var kill = flag.Bool("kill", false, "docker kill postgres")

func setup() {
	const dsn = comptests.DefaultDSN
	comptests.SetupDockerDB(dsn)
	model.SetupDBOnceOnly(dsn)
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

}

func teardown() {
	if *kill {
		comptests.KillDockerDB()
	}
}

// test runner function for comptests
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
