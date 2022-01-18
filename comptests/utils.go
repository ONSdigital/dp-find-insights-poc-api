package comptests

import (
	"fmt"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
)

func DoSQL(t *testing.T, db *database.Database, sql string) {
	_, err := db.DB().Exec(sql)
	if err != nil {
		t.Fatal(err)
	}
}

func ClearDB(db *database.Database) error {
	// order of delete-froms matters!
	tables := []string{
		"geo_metric",
		"geo",
		"nomis_category",
		"nomis_desc",
		"nomis_topic",
		"geo_type",
		"data_ver",
	}
	for _, table := range tables {
		_, err := db.DB().Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return err
		}
	}
	return nil
}
