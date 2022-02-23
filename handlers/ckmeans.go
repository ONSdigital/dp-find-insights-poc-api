package handlers

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
)

func (svr *Server) GetCkmeansYear(w http.ResponseWriter, r *http.Request, year int, params api.GetCkmeansYearParams) {
	if !svr.assertAuthorized(w, r) || !svr.assertDatabaseEnabled(w, r) {
		return
	}

	generate := func() ([]byte, error) {
		var cat, geotype []string
		var divideBy string
		var k int
		if params.Cat != nil {
			cat = *params.Cat
		}
		if params.Geotype != nil {
			geotype = *params.Geotype
		}
		if params.K != nil {
			k = *params.K
		}
		if params.DivideBy != nil {
			divideBy = *params.DivideBy
		}
		if cat == nil || geotype == nil || k == 0 {
			return nil, fmt.Errorf("%w: cat, geotype and k required", geodata.ErrMissingParams)
		}

		ctx := r.Context()
		breaks, err := svr.querygeodata.CKmeans(ctx, year, cat, geotype, k, divideBy)
		if err != nil {
			return nil, err
		}
		return toJSON(breaks)
	}

	svr.respond(w, r, mimeJSON, generate)
}
