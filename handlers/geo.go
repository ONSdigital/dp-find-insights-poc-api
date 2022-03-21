package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
)

func (svr *Server) GetGeo(w http.ResponseWriter, r *http.Request, year int, params api.GetGeoParams) {
	if !svr.assertAuthorized(w, r) || !svr.assertDatabaseEnabled(w, r) {
		return
	}

	// geocode is not marked required in swaggger.yaml because we expect different
	// combinations of query parameters in future
	if params.Geocode == nil {
		sendError(w, http.StatusBadRequest, "geocode query parameter required")
		return
	}

	generate := func() ([]byte, error) {
		resp, err := svr.querygeodata.Geo(r.Context(), year, *params.Geocode)
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
