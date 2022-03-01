package handlers

import (
	"net/http"
)

func (svr *Server) GetGeoYearRegion(w http.ResponseWriter, r *http.Request, year int, region string) {
	if !svr.assertAuthorized(w, r) || !svr.assertDatabaseEnabled(w, r) {
		return
	}

	generate := func() ([]byte, error) {
		geo, err := svr.querygeodata.Geo(r.Context(), year, region)
		if err != nil {
			return nil, err
		}
		buf, err := toJSON(geo)
		if err != nil {
			return nil, err
		}
		return []byte(buf), err

	}

	svr.respond(w, r, mimeJSON, generate)
}
