package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/table"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/where"
	"github.com/ONSdigital/dp-find-insights-poc-api/sentinel"
)

func (svr *Server) GetQuery(w http.ResponseWriter, r *http.Request, year int, params api.GetQueryParams) {
	if !svr.assertAuthorized(w, r) || !svr.assertDatabaseEnabled(w, r) {
		return
	}

	generate := func() ([]byte, error) {
		var rows []string
		var bbox string
		var geotype []string
		var location string
		var radius int
		var polygon string
		if params.Rows != nil {
			rows = *params.Rows
		}
		if params.Bbox != nil {
			bbox = *params.Bbox
		}
		if params.Geotype != nil {
			geotype = *params.Geotype
		}
		if params.Location != nil {
			location = *params.Location
		}
		if params.Radius != nil {
			radius = *params.Radius
		}
		if params.Polygon != nil {
			polygon = *params.Polygon
		}

		geocodes, err := svr.querygeodata.Query2(r.Context(), year, bbox, location, radius, polygon, geotype, rows)
		if err != nil {
			return nil, err
		}

		var cols []string
		if params.Cols != nil {
			cols = *params.Cols
		}
		var censustable string
		if params.Censustable != nil {
			censustable = *params.Censustable
		}

		// parse cols query strings into a ValueSet
		catset, err := where.ParseMultiArgs(cols)
		if err != nil {
			return nil, err
		}

		// extract special column names from ValueSet
		include, catset, err := geodata.ExtractSpecialCols(catset)
		if err != nil {
			return nil, err
		}

		// special case for dev: explicit cols="geocode" and no census table means just print geocodes column
		// (would just allow cols=geography_code, but that already means all columns)
		if len(catset.Singles) == 0 && len(catset.Ranges) == 0 && len(include) == 1 && include[0] == table.ColGeocodes && censustable == "" {
			return geocodeCSV(geocodes)
		}

		if year == 2011 {
			return svr.querygeodata.PGMetrics(r.Context(), year, geocodes, catset, include, censustable)
		}
		if len(geotype) != 1 {
			return nil, fmt.Errorf("%w: cantabular queries require a single geotype", sentinel.ErrInvalidParams)
		}

		return svr.querygeodata.CantabularMetrics(r.Context(), geocodes, catset, geotype[0])
	}

	svr.respond(w, r, mimeCSV, generate)
}

func geocodeCSV(geocodes []string) ([]byte, error) {
	var body bytes.Buffer
	cw := csv.NewWriter(&body)
	cw.Write([]string{"geocode"})
	for _, geocode := range geocodes {
		cw.Write([]string{geocode})
	}
	cw.Flush()
	return body.Bytes(), cw.Error()
}
