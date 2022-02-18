package model

import "gorm.io/gorm"

// this is the schema for Postgres database)

type SchemaVer struct {
	gorm.Model // updated_at etc
	BuildTime  string
	GitCommit  string
	Version    string
}

// don't pluralise table name
func (SchemaVer) TableName() string {
	return "schema_ver"
}

type DataVer struct {
	gorm.Model       // updated_at etc
	ID         int32 `gorm:"primaryKey;autoIncrement:false"`
	CensusYear int32
	VerString  string
	Source     string
	Notes      string
	Public     bool
	GoMetrics  []GeoMetric `gorm:"foreignKey:DataVerID;references:ID"`
}

// don't pluralise table name
func (DataVer) TableName() string {
	return "data_ver"
}

// THIS TABLE NEEDS RESTRUCTURING & FIELDS RENAMING

type YearMapping struct {
	ID           int32 `gorm:"primaryKey"`
	Lsoa2011code string
	Lad2020code  string
}

// don't pluralise table name
func (YearMapping) TableName() string {
	return "lsoa2011_lad2020_lookup"
}

type GeoType struct {
	ID   int32 `gorm:"primaryKey;autoIncrement:false"`
	Name string
	Geos []Geo `gorm:"foreignKey:TypeID;references:ID"`
}

// don't pluralise table name
func (GeoType) TableName() string {
	return "geo_type"
}

type Geo struct {
	ID     int32 `gorm:"primaryKey"`
	TypeID int32
	Code   string `gorm:"index:unique"`
	Name   string
	Lat    float64
	Long   float64
	Valid  bool `gorm:"DEFAULT:true"`
	// wkb_geometry - added via ALTER
	// wkb_long_lat_geom - added via ALTER
	GoMetrics []GeoMetric `gorm:"foreignKey:GeoID;references:ID"`
}

// don't pluralise table name
func (Geo) TableName() string {
	return "geo"
}

type GeoMetric struct {
	ID         int32 `gorm:"primaryKey"`
	GeoID      int32 `gorm:"index"`
	CategoryID int32 `gorm:"index"`
	Metric     float64
	DataVerID  int32
}

// don't pluralise table name
func (GeoMetric) TableName() string {
	return "geo_metric"
}

type NomisCategory struct {
	// why do we need uniqueIndex? composite key!
	ID              int32 `gorm:"uniqueIndex;primaryKey"`
	NomisDescID     int32 `gorm:"primaryKey"`
	CategoryName    string
	MeasurementUnit string
	StatUnit        string
	LongNomisCode   string `gorm:"uniqueIndex"`
	Year            int32
	GoMetrics       []GeoMetric `gorm:"foreignKey:CategoryID;references:ID"`
}

// don't pluralise table name
func (NomisCategory) TableName() string {
	return "nomis_category"
}

type NomisDesc struct {
	ID              int32 `gorm:"uniqueIndex;primaryKey"`
	NomisTopicID    int32 `gorm:"primaryKey"`
	Name            string
	PopStat         string
	ShortNomisCode  string `gorm:"uniqueIndex"`
	Year            int32
	NomisCategories []NomisCategory `gorm:"foreignKey:NomisDescID;references:ID"`
}

// don't pluralise table name
func (NomisDesc) TableName() string {
	return "nomis_desc"
}

type NomisTopic struct {
	ID           int32 `gorm:"primaryKey"`
	TopNomisCode string
	Name         string
	NomisDescs   []NomisDesc `gorm:"foreignKey:NomisTopicID;references:ID"`
}

// don't pluralise table name
func (NomisTopic) TableName() string {
	return "nomis_topic"
}

// data prepopulated in Postgres database

// GetGeoTypeValues returns a slice of geo types
func GetGeoTypeValues() []string {
	// XXX LSOA to be removed
	return []string{"EW", "Country", "Region", "LAD", "MSOA", "LSOA"}
}

// GetGeoTypeValues returns a map of geo types for validation
func GetGeoTypeMap() map[string]bool {
	m := make(map[string]bool)

	for _, v := range GetGeoTypeValues() {
		m[v] = true
	}

	return m
}
