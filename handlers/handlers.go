package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
	"github.com/ONSdigital/dp-find-insights-poc-api/cache"
	"github.com/ONSdigital/dp-find-insights-poc-api/config"
	"github.com/ONSdigital/dp-find-insights-poc-api/metadata"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
	Swagger "github.com/ONSdigital/dp-find-insights-poc-api/swagger"
)

type Server struct {
	private      bool             // true if private endpoint feature flag is enabled
	querygeodata *geodata.Geodata // if nil, database not available
	md           *metadata.Metadata
	cm           *cache.Manager
}

func New(private bool, querygeodata *geodata.Geodata, md *metadata.Metadata, cm *cache.Manager) *Server {
	return &Server{
		private:      private,
		querygeodata: querygeodata,
		md:           md,
		cm:           cm,
	}
}

func (svr *Server) GetSwagger(w http.ResponseWriter, r *http.Request) {
	spec, _ := Swagger.GetOpenAPISpec()
	b, err := spec.MarshalJSON()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write(b)
}

func (svr *Server) GetSwaggerui(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "html")
	w.WriteHeader(http.StatusOK)
	c, err := config.Get()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	b, err := Swagger.GetSwaggerUIPage("http://"+c.BindAddr+"/swagger", "")
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write(b)
}

func (svr *Server) GetMetadataYear(w http.ResponseWriter, r *http.Request, year int, params api.GetMetadataYearParams) {
	generate := func() ([]byte, error) {
		var filtertotals bool
		if params.Filtertotals != nil {
			filtertotals = *params.Filtertotals
		} else {
			filtertotals = false
		}

		return svr.md.Get(year, filtertotals)
	}

	svr.respond(w, r, mimeCSV, generate)
}

func (svr *Server) GetQueryYear(w http.ResponseWriter, r *http.Request, year int, params api.GetQueryYearParams) {
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
		var rows []string
		var cols []string
		var bbox string
		var geotype []string
		var location string
		var radius int
		var polygon string
		var censustable string
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
		if params.Polygon != nil {
			polygon = *params.Polygon
		}
		if params.Censustable != nil {
			censustable = *params.Censustable
		}

		ctx := r.Context()
		csv, err := svr.querygeodata.Query(ctx, year, bbox, location, radius, polygon, geotype, rows, cols, censustable)
		return []byte(csv), err
	}

	svr.respond(w, r, mimeCSV, generate)
}

func sendError(w http.ResponseWriter, code int, message string) {
	log.Printf("request error: %s", message)
	e := api.Error{
		Error: message,
	}
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		log.Printf("SendResult: %s", err)
	}
}

func toJSON(v interface{}) ([]byte, error) {
	buf, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return nil, err
	}
	return append(buf, "\n"...), nil
}
