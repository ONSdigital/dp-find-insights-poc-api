package model

import (
	"fmt"
	"log"
	"regexp"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TODO fuller SQL logs

// prepare DB
func SetupDB(dsn string) {

	_, pw, host, port, db := ParseDSN(dsn)

	{
		db, err := gorm.Open(postgres.Open(CreatDSN("postgres", pw, host, port, "postgres")), &gorm.Config{})
		if err != nil {
			log.Print(err)
		}

		// should replace creatdbuser.sh
		execSQL(db, []string{
			"CREATE DATABASE censustest",
			"CREATE USER insights WITH PASSWORD 'insights'",
			"ALTER USER insights WITH CREATEDB"})
	}

	{
		db, err := gorm.Open(postgres.Open(CreatDSN("postgres", pw, host, port, db)), &gorm.Config{})
		if err != nil {
			log.Print(err)
		}

		// should replace creatdb.sh
		execSQL(db, []string{"CREATE EXTENSION postgis"})
	}

	{
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err != nil {
			log.Print(err)
		}

		Migrate(db)

		// XXX checkme
		db.Save(&DataVer{ID: 1, CensusYear: 2011, VerString: "2.2", Public: true, Source: "Nomis Bulk API", Notes: "Release date 12/02/2013 Revised 17/01/2014"})

		DataPopulate(db)

	}

}

// special case for use in comptests - only setup db if it has not already been set up. This is safe to call multiple times against the same db.
func SetupDBOnceOnly(dsn string) {
	// assume if censustest db exists, the work is done
	_, pw, host, port, _ := ParseDSN(dsn)

	{
		db, err := gorm.Open(postgres.Open(CreatDSN("postgres", pw, host, port, "postgres")), &gorm.Config{})
		if err != nil {
			log.Print(err)
		}

		// check if db exists already, return if so
		var isDBCreated bool
		db.Raw("SELECT EXISTS (SELECT datname FROM pg_catalog.pg_database WHERE datname='censustest')").Scan(&isDBCreated)
		if isDBCreated {
			log.Println("Database already setup, skipping...")
			return
		}
	}

	// else run SetupDB
	SetupDB(dsn)
}

// setup schema
func Migrate(db *gorm.DB) {

	// XXX create/alter tables - doesn't delete cols or tables!
	// neither does it always change types correctly
	// More useful in dev than prod

	if err := db.AutoMigrate(
		&NomisTopic{},
		&SchemaVer{},
		&DataVer{},
		&GeoType{},
		&Geo{},
		&NomisDesc{},
		&NomisCategory{},
		&GeoMetric{},
		&YearMapping{},
	); err != nil {
		log.Print(err)
	}

	execSQL(db, []string{
		"ALTER TABLE geo ADD COLUMN wkb_geometry geometry(Geometry,4326)",
		"CREATE INDEX geo_wkb_geometry_geom_idx ON public.geo USING gist (wkb_geometry)",
		"ALTER TABLE geo ADD COLUMN wkb_long_lat_geom geometry(Geometry,4326)",
		"CREATE INDEX geo_long_lat_geom_idx ON public.geo USING gist ( wkb_long_lat_geom)"})

}

func DataPopulate(db *gorm.DB) {

	// populate topic -- top level metadata
	db.Save(&NomisTopic{ID: 1, TopNomisCode: "QS1", Name: "Population Basics"})
	db.Save(&NomisTopic{ID: 2, TopNomisCode: "QS2", Name: "Origins & Beliefs"})
	db.Save(&NomisTopic{ID: 3, TopNomisCode: "QS3", Name: "Health"})
	db.Save(&NomisTopic{ID: 4, TopNomisCode: "QS4", Name: "Housing"})
	db.Save(&NomisTopic{ID: 5, TopNomisCode: "QS5", Name: "Education"})
	db.Save(&NomisTopic{ID: 6, TopNomisCode: "QS6", Name: "Employment"})
	db.Save(&NomisTopic{ID: 7, TopNomisCode: "QS7", Name: "Travel to Work"})
	db.Save(&NomisTopic{ID: 8, TopNomisCode: "QS8", Name: "Residency"})

	// FK relationship for topic

	execSQL(db, []string{
		"UPDATE nomis_desc SET nomis_topic_id=1 WHERE short_nomis_code LIKE 'QS1%'",
		"UPDATE nomis_desc SET nomis_topic_id=2 WHERE short_nomis_code LIKE 'QS2%'",
		"UPDATE nomis_desc SET nomis_topic_id=3 WHERE short_nomis_code LIKE 'QS3%'",
		"UPDATE nomis_desc SET nomis_topic_id=4 WHERE short_nomis_code LIKE 'QS4%'",
		"UPDATE nomis_desc SET nomis_topic_id=5 WHERE short_nomis_code LIKE 'QS5%'",
		"UPDATE nomis_desc SET nomis_topic_id=6 WHERE short_nomis_code LIKE 'QS6%'",
		"UPDATE nomis_desc SET nomis_topic_id=7 WHERE short_nomis_code LIKE 'QS7%'",
		"UPDATE nomis_desc SET nomis_topic_id=8 WHERE short_nomis_code LIKE 'QS8%'",
	})

}

func execSQL(db *gorm.DB, ss []string) {
	for _, s := range ss {
		if err := db.Exec(s).Error; err != nil {
			log.Print(err)
		}
	}
}

func ParseDSN(dsn string) (user, pw, host, port, db string) {
	re := regexp.MustCompile(`postgres://(.*):(.*)@(.*):(.*)/(.*)`)
	match := re.FindStringSubmatch(dsn)

	if len(match) != 6 {
		log.Fatal("match fail")
	}

	return match[1], match[2], match[3], match[4], match[5]
}

func CreatDSN(user, pw, host, port, db string) (dsn string) {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pw, host, port, db)
}
