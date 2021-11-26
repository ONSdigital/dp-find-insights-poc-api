package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/demo"
)

type Server struct {
	private   bool       // true if private endpoint feature flag is enabled
	queryDemo *demo.Demo // if nil, database not available
}

func New(private bool, queryDemo *demo.Demo) *Server {
	return &Server{
		private:   private,
		queryDemo: queryDemo,
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
	if svr.queryDemo == nil {
		sendError(w, http.StatusNotImplemented, "database not enabled")
		return
	}
	var rows []string
	var cols []string
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

	ctx := r.Context()
	csv, err := svr.queryDemo.Query(ctx, dataset, rows, cols)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, demo.ErrTooManyMetrics) {
			status = http.StatusForbidden
		} else if errors.Is(err, demo.ErrMissingParams) || errors.Is(err, demo.ErrInvalidTable) {
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
