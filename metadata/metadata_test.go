//go:build comptest
// +build comptest

package metadata

import (
	"log"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/comptests"
	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = comptests.DefaultDSN

var db *gorm.DB

func init() {
	comptests.SetupDockerDB(dsn)
	model.SetupUpdateDB(dsn)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		//		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Print(err)
	}
}

func TestMetaDataTest(t *testing.T) {
	// inside transaction rolled back
	func() {
		tx := db.Begin()
		defer tx.Rollback()

		// this is prepopulated so make the result much smaller!
		db.Exec("DELETE FROM nomis_topic WHERE id>1")

		db.Exec("INSERT INTO NOMIS_DESC (id,name,pop_stat,short_nomis_code,year,nomis_topic_id) VALUES (15,'Families with dependent children','All families in households; All dependent children in households','QS118EW',2011,1)")

		db.Exec("INSERT INTO NOMIS_CATEGORY (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year) VALUES (211,15,'All categories: Dependent children in family','Count','Family','QS118EW0001',2011)")

		md, _ := New(db)

		filterTotals := false
		b, err := md.Get(2011, filterTotals)
		if err != nil {
			t.Error(err)
		}

		if string(b) != result() {
			println(string(b))
			t.Fail()
		}
	}()
}

func result() string {
	return `[{"code":"QS1","name":"Population Basics","slug":"population-basics","tables":[{"categories":[{"code":"QS118EW0001","name":"All categories: Dependent children in family","slug":"all-categories-dependent-children-in-family"},{"code":"QS118EW0002","name":"foo blah etc","slug":"foo-blah-etc"}],"code":"QS118EW","name":"Families with dependent children","slug":"families-with-dependent-children"}]}]`
}

func TestMetaDataFiltertotals(t *testing.T) {
	// inside transaction rolled back
	func() {
		tx := db.Begin()
		defer tx.Rollback()

		// this is prepopulated so make the result much smaller!
		db.Exec("DELETE FROM nomis_topic WHERE id>1")

		db.Exec("INSERT INTO NOMIS_DESC (id,name,pop_stat,short_nomis_code,year,nomis_topic_id) VALUES (15,'Families with dependent children','All families in households; All dependent children in households','QS118EW',2011,1)")

		db.Exec("INSERT INTO NOMIS_CATEGORY (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year) VALUES (211,15,'All categories: Dependent children in family','Count','Family','QS118EW0001',2011)")
		db.Exec("INSERT INTO NOMIS_CATEGORY (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year) VALUES (212,15,'foo blah etc','Count','Family','QS118EW0002',2011)")

		md, _ := New(db)

		filterTotals := true
		b, err := md.Get(2011, filterTotals)
		if err != nil {
			t.Error(err)
		}

		if string(b) != resultFilterTotals() {
			println(string(b))
			t.Fail()
		}
	}()
}

func resultFilterTotals() string {
	return `[{"code":"QS1","name":"Population Basics","slug":"population-basics","tables":[{"categories":[{"code":"QS118EW0002","name":"foo blah etc","slug":"foo-blah-etc"}],"code":"QS118EW","name":"Families with dependent children","slug":"families-with-dependent-children","total":{"code":"QS118EW0001","name":"All categories: Dependent children in family","slug":"all-categories-dependent-children-in-family"}}]}]`
}
