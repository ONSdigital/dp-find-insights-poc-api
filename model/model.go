package model

// this is the schema for the "new" (row-based Postgres database)

type GeoType struct {
	ID          int32 `gorm:"primaryKey;autoIncrement:false"`
	GeoTypeName string
	Geos        []Geo `gorm:"foreignKey:GeoTypeID;references:ID"`
}

// don't pluralise table name
func (GeoType) TableName() string {
	return "geo_type"
}

type Geo struct {
	ID        int32 `gorm:"primaryKey;autoIncrement:false"`
	GeoTypeID int32
	GeoCode   string
	GeoName   string
	GoMetrics []GeoMetric `gorm:"foreignKey:GeoID;references:ID"`
}

// don't pluralise table name
func (Geo) TableName() string {
	return "geo"
}

type GeoMetric struct {
	ID         int32 `gorm:"primaryKey"`
	GeoID      int32
	CategoryID int32
	Metric     float64
	Year       int32
}

// don't pluralise table name
func (GeoMetric) TableName() string {
	return "geo_metric"
}

type NomisCategory struct {
	// why do we need uniqueIndex?
	ID              int32 `gorm:"uniqueIndex;primaryKey"`
	NomisDescID     int32 `gorm:"primaryKey"`
	CategoryName    string
	MeasurementUnit string
	StatUnit        string
	LongNomisCode   string
	Year            int32
	GoMetrics       []GeoMetric `gorm:"foreignKey:CategoryID;references:ID"`
}

// don't pluralise table name
func (NomisCategory) TableName() string {
	return "nomis_category"
}

type NomisDesc struct {
	ID              int32 `gorm:"primaryKey"`
	LongDesc        string
	ShortDesc       string
	ShortNomisCode  string
	Year            int32
	NomisCategories []NomisCategory `gorm:"foreignKey:NomisDescID;references:ID"`
}

// don't pluralise table name
func (NomisDesc) TableName() string {
	return "nomis_desc"
}
