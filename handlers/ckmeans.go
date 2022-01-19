package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
	"github.com/ONSdigital/dp-find-insights-poc-api/config"
)

func (svr *Server) GetDevCkmeans(w http.ResponseWriter, r *http.Request, params api.GetDevCkmeansParams) {
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

	// add CORS header
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
		sendError(w, http.StatusBadRequest, "cat, geotype and k required")
		return
	}

	ctx := r.Context()
	breaks, err := svr.querygeodata.CKmeans(ctx, cat, geotype, k)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	_ = encoder.Encode(breaks)
}

func (svr *Server) GetDevCkmeansratio(w http.ResponseWriter, r *http.Request, params api.GetDevCkmeansratioParams) {
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

	// add CORS header
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
		sendError(w, http.StatusBadRequest, "cat1, cat2, geotype and k required")
		return
	}

	ctx := r.Context()
	breaks, err := svr.querygeodata.CKmeansRatio(ctx, cat1, cat2, geotype, k)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	_ = encoder.Encode(breaks)
}
