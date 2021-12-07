package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
)

type Server struct {
	private      bool             // true if private endpoint feature flag is enabled
	querygeodata *geodata.Geodata // if nil, database not available
}

func New(private bool, querygeodata *geodata.Geodata) *Server {
	return &Server{
		private:      private,
		querygeodata: querygeodata,
	}
}

func (svr *Server) GetHello(w http.ResponseWriter, r *http.Request) {
	response := &api.HelloResponse{
		Message: "Hello, world!",
	}
	sendResult(w, http.StatusOK, response)
}

func (svr *Server) GetDevHelloDataset(w http.ResponseWriter, r *http.Request, dataset string, params api.GetDevHelloDatasetParams) {
	if !svr.private {
		sendError(w, http.StatusNotFound, "endpoint not enabled")
		return
	}
	if svr.querygeodata == nil {
		sendError(w, http.StatusNotImplemented, "database not enabled")
		return
	}
	var rows []string
	var cols []string
	var bbox string
	var geotype []string
	var location string
	var radius int
	if dataset == "" {
		sendError(w, http.StatusBadRequest, "dataset missing")
		return
	}
	if params.Rows != nil {
		rows = *params.Rows
	}
	if params.Cols != nil {
		cols = *params.Cols
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

	ctx := r.Context()
	csv, err := svr.querygeodata.Query(ctx, dataset, bbox, location, radius, geotype, rows, cols)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, geodata.ErrTooManyMetrics) {
			status = http.StatusForbidden
		} else if errors.Is(err, geodata.ErrMissingParams) || errors.Is(err, geodata.ErrInvalidTable) {
			status = http.StatusBadRequest
		}
		sendError(w, status, err.Error())
		return
	}
	sendCSV(w, http.StatusOK, csv)
}

func sendResult(w http.ResponseWriter, code int, result interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("SendResult: %s", err)
	}
}

func sendError(w http.ResponseWriter, code int, message string) {
	e := api.Error{
		Error: message,
	}
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		log.Printf("SendResult: %s", err)
	}
}

func sendCSV(w http.ResponseWriter, code int, message string) {
	w.Header().Add("Content-Type", "text/csv")
	w.WriteHeader(code)
	w.Write([]byte(message))
}
