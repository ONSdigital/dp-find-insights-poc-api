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

// SetupUpdate is used both to create and update the database
// It's OK to call this more than once on the same DB
func SetupUpdateDB(dsn string) {

	_, pw, host, port, db := ParseDSN(dsn)

	{
		gdb, err := gorm.Open(postgres.Open(CreatDSN("postgres", pw, host, port, "postgres")), &gorm.Config{})
		if err != nil {
			log.Print(err)
		}

		execSQL(gdb, []string{
			"CREATE USER insights WITH PASSWORD 'insights'", // XXX pw hardcoded
			"CREATE DATABASE " + db + " WITH OWNER insights",
			"ALTER USER insights WITH CREATEDB"})
	}

	{
		gdb, err := gorm.Open(postgres.Open(CreatDSN("postgres", pw, host, port, db)), &gorm.Config{})
		if err != nil {
			log.Print(err)
		}

		execSQL(gdb, []string{"CREATE EXTENSION IF NOT EXISTS postgis"})
	}

	{
		gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err != nil {
			log.Print(err)
		}

		Migrate(gdb)

		// XXX checkme
		gdb.Save(&DataVer{ID: 1, CensusYear: 2011, VerString: "2.2", Public: true, Source: "Nomis Bulk API", Notes: "20220117 2i based on metadata/i2.txt (fewer QS + some KS rows)"})

		DataPopulate(gdb)

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
	SetupUpdateDB(dsn)
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

	// id=0 is undefined topic, can't see how to do this with gorm!
	// we need this when we import data before FK set up as default value
	// in "nomis-bulk-to-postgres/add_to_db.py" function "add_meta_tables"

	// suppress error on multiple insert
	execSQL(db, []string{
		"INSERT INTO NOMIS_TOPIC (id) VALUES (0) ON CONFLICT (id) DO NOTHING",
	})

	// populate topic -- top level metadata
	db.Save(&NomisTopic{ID: 1, TopNomisCode: "QS1", Name: "Population Basics"})
	db.Save(&NomisTopic{ID: 2, TopNomisCode: "QS2", Name: "Origins & Beliefs"})
	db.Save(&NomisTopic{ID: 3, TopNomisCode: "QS3", Name: "Health"})
	db.Save(&NomisTopic{ID: 4, TopNomisCode: "QS4", Name: "Housing"})
	db.Save(&NomisTopic{ID: 5, TopNomisCode: "QS5", Name: "Education"})
	db.Save(&NomisTopic{ID: 6, TopNomisCode: "QS6", Name: "Employment"})
	db.Save(&NomisTopic{ID: 7, TopNomisCode: "QS7", Name: "Travel to Work"})
	db.Save(&NomisTopic{ID: 8, TopNomisCode: "QS8", Name: "Residency"})

	db.Save(&NomisTopic{ID: 100, TopNomisCode: "KS1", Name: "Population Basics"})
	db.Save(&NomisTopic{ID: 200, TopNomisCode: "KS2", Name: "Origins & Beliefs"})
	db.Save(&NomisTopic{ID: 400, TopNomisCode: "KS4", Name: "Housing"})

	// XXX DC6 "Population Basics"
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
		"UPDATE nomis_desc SET nomis_topic_id=100 WHERE short_nomis_code LIKE 'KS1%'",
		"UPDATE nomis_desc SET nomis_topic_id=200 WHERE short_nomis_code LIKE 'KS2%'",
		"UPDATE nomis_desc SET nomis_topic_id=400 WHERE short_nomis_code LIKE 'KS4%'",
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
