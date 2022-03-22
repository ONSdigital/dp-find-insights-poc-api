package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
)

func (svr *Server) GetGeo(w http.ResponseWriter, r *http.Request, year int, params api.GetGeoParams) {
	if !svr.assertAuthorized(w, r) || !svr.assertDatabaseEnabled(w, r) {
		return
	}

	// need either geocode OR geoname, but not both!
	var geocode, geoname string
	if params.Geocode != nil {
		geocode = *params.Geocode
	}
	if params.Geoname != nil {
		geoname = *params.Geoname
	}
	if (geocode == "" && geoname == "") || (geocode != "" && geoname != "") {
		sendError(w, http.StatusBadRequest, "geocode OR geoname query parameter required")
		return
	}

	generate := func() ([]byte, error) {
		resp, err := svr.querygeodata.Geo(r.Context(), year, geocode, geoname)
		if err != nil {
			return nil, err
		}
		buf, err := toJSON(resp)
		if err != nil {
			return nil, err
		}
		return []byte(buf), err
	}

	svr.respond(w, r, mimeJSON, generate)
}
