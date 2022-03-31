package handlers

// XXX move config into Server and don't call config.Get

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
	"github.com/ONSdigital/dp-find-insights-poc-api/cache"
	"github.com/ONSdigital/dp-find-insights-poc-api/config"
	"github.com/ONSdigital/dp-find-insights-poc-api/metadata"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
	"github.com/ONSdigital/dp-find-insights-poc-api/postcode"
	Swagger "github.com/ONSdigital/dp-find-insights-poc-api/swagger"
	"github.com/ONSdigital/log.go/v2/log"
)

type Server struct {
	private      bool             // true if private endpoints are enabled
	querygeodata *geodata.Geodata // if nil, database not available
	md           *metadata.Metadata
	cm           *cache.Manager
	pc           *postcode.Postcode
}

func New(private bool, querygeodata *geodata.Geodata, md *metadata.Metadata, cm *cache.Manager, pc *postcode.Postcode) *Server {
	return &Server{
		private:      private,
		querygeodata: querygeodata,
		md:           md,
		cm:           cm,
		pc:           pc,
	}
}

func (svr *Server) GetSwagger(w http.ResponseWriter, r *http.Request) {
	spec, _ := Swagger.GetOpenAPISpec()
	b, err := spec.MarshalJSON()
	if err != nil {
		sendError(r.Context(), w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write(b)
}

func (svr *Server) GetSwaggerui(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "html")
	w.WriteHeader(http.StatusOK)
	ctx := r.Context()
	c, err := config.Get()
	if err != nil {
		sendError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}
	b, err := Swagger.GetSwaggerUIPage("http://"+c.BindAddr+"/swagger", "")
	if err != nil {
		sendError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write(b)
}

func (svr *Server) GetMetadataYear(w http.ResponseWriter, r *http.Request, year int, params api.GetMetadataYearParams) {
	if !svr.assertAuthorized(w, r) || !svr.assertDatabaseEnabled(w, r) {
		return
	}

	generate := func() ([]byte, error) {
		var filtertotals bool
		if params.Filtertotals != nil {
			filtertotals = *params.Filtertotals
		} else {
			filtertotals = false
		}

		return svr.md.Get(r.Context(), year, filtertotals)
	}

	svr.respond(w, r, mimeCSV, generate)
}

func (svr *Server) GetMsoaPostcode(w http.ResponseWriter, r *http.Request, pc string) {

	generate := func() ([]byte, error) {
		code, name, err := svr.pc.GetMSOA(pc)
		if err != nil {
			return nil, err
		}
		// JSON?
		return []byte(code + ", " + name + "\r\n"), nil
	}

	svr.respond(w, r, mimeCSV, generate)
}

func (svr *Server) GetQueryYear(w http.ResponseWriter, r *http.Request, year int, params api.GetQueryYearParams) {
	if !svr.assertAuthorized(w, r) || !svr.assertDatabaseEnabled(w, r) {
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

func (svr *Server) GetClearCache(w http.ResponseWriter, r *http.Request) {
	if !svr.assertPrivate(w, r) || !svr.assertAuthorized(w, r) {
		return
	}

	ctx := r.Context()
	err := svr.cm.Clear(ctx)
	if err == nil {
		return
	}
	sendError(ctx, w, http.StatusInternalServerError, "problem clearing cache", log.Data{"error": err.Error()})
}

func (svr *Server) Preflight(w http.ResponseWriter, r *http.Request, path string, year int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
}

// assertPrivate sends an error to the client if private endpoints are not enabled.
// Returns true if private endpoints are enabled.
func (svr *Server) assertPrivate(w http.ResponseWriter, r *http.Request) bool {
	if svr.private {
		return true
	}
	sendError(r.Context(), w, http.StatusNotFound, "endpoint not enabled")
	return false
}

// assertAuthorized send an error to the client if they are not authorized.
// Returns true if authorized.
func (svr *Server) assertAuthorized(w http.ResponseWriter, req *http.Request) bool {
	// check Auth header
	c, _ := config.Get()
	if !c.EnableHeaderAuth {
		return true
	}

	ahdr := "Authorization"
	auth := req.Header.Get(ahdr)
	if auth == c.APIToken {
		return true
	}

	xhdr := "X-Forwarded-For"
	sendError(
		req.Context(),
		w,
		http.StatusUnauthorized,
		"unauthorized",
		log.Data{ahdr: auth, xhdr: req.Header.Get(xhdr)},
	)
	return false
}

// assertDatabaseEnabled sends and error to the client if database is not enabled.
// Returns true if the database is enabled.
func (svr *Server) assertDatabaseEnabled(w http.ResponseWriter, req *http.Request) bool {
	if svr.querygeodata != nil {
		return true
	}
	sendError(req.Context(), w, http.StatusNotImplemented, "database not enabled")
	return false
}

// sendError returns an http status code and message to the client, and logs a message.
// message is the string sent to the client in the body of the response.
// code is the http status code.
func sendError(ctx context.Context, w http.ResponseWriter, code int, message string, opts ...log.Data) {
	event := "request error"
	if code == http.StatusInternalServerError {
		event = "internal error"
	}

	data := log.Data{"response": message}
	for _, opt := range opts {
		for k, v := range opt {
			data[k] = v
		}
	}
	log.Info(ctx, event, log.Data{"response": message}, data)

	e := api.Error{
		Error: message,
	}
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		log.Error(ctx, "json encode", err)
	}
}

func toJSON(v interface{}) ([]byte, error) {
	buf, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return nil, err
	}
	return append(buf, "\n"...), nil
}
