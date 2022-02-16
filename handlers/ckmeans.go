package handlers

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
	"github.com/ONSdigital/dp-find-insights-poc-api/config"
)

func (svr *Server) GetCkmeansYear(w http.ResponseWriter, r *http.Request, year int, params api.GetCkmeansYearParams) {
	if !svr.private {
		sendError(w, http.StatusNotFound, "endpoint not enabled")
		return
	}

	// check Auth header
	c, _ := config.Get()
	if c.EnableHeaderAuth {
		auth := r.Header.Get("Authorization")
		if auth != c.APIToken {
			sendError(w, http.StatusUnauthorized, "unauthorized")
			fmt.Printf("failed auth header '%s' from '%s'", auth, r.Header.Get("X-Forwarded-For"))
			return
		}
	}

	if svr.querygeodata == nil {
		sendError(w, http.StatusNotImplemented, "database not enabled")
		return
	}

	generate := func() ([]byte, error) {
		var cat, geotype string
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
		if cat == "" || geotype == "" || k == 0 {
			return nil, fmt.Errorf("%w: cat, geotype and k required", ErrBadRequest)
		}

		ctx := r.Context()
		breaks, err := svr.querygeodata.CKmeans(ctx, year, cat, geotype, k)
		if err != nil {
			return nil, err
		}
		return toJSON(breaks)
	}

	svr.respond(w, r, mimeJSON, generate)
}

func (svr *Server) GetCkmeansratioYear(w http.ResponseWriter, r *http.Request, year int, params api.GetCkmeansratioYearParams) {
	if !svr.private {
		sendError(w, http.StatusNotFound, "endpoint not enabled")
		return
	}

	// check Auth header
	c, _ := config.Get()
	if c.EnableHeaderAuth {
		auth := r.Header.Get("Authorization")
		if auth != c.APIToken {
			sendError(w, http.StatusUnauthorized, "unauthorized")
			fmt.Printf("failed auth header '%s' from '%s'", auth, r.Header.Get("X-Forwarded-For"))
			return
		}
	}

	if svr.querygeodata == nil {
		sendError(w, http.StatusNotImplemented, "database not enabled")
		return
	}

	generate := func() ([]byte, error) {
		var cat1, cat2, geotype string
		var k int
		if params.Cat1 != nil {
			cat1 = *params.Cat1
		}
		if params.Cat2 != nil {
			cat2 = *params.Cat2
		}
		if params.Geotype != nil {
			geotype = *params.Geotype
		}
		if params.K != nil {
			k = *params.K
		}
		if cat1 == "" || cat2 == "" || geotype == "" || k == 0 {
			return nil, fmt.Errorf("%w: cat1, cat2, geotype and k required", ErrBadRequest)
		}

		ctx := r.Context()
		breaks, err := svr.querygeodata.CKmeansRatio(ctx, year, cat1, cat2, geotype, k)
		if err != nil {
			return nil, err
		}
		return toJSON(breaks)
	}

	svr.respond(w, r, mimeJSON, generate)
}
