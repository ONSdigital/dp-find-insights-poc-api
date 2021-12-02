// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// Error defines model for Error.
type Error struct {
	// error message
	Error string `json:"error"`
}

// HelloResponse defines model for HelloResponse.
type HelloResponse struct {
	// Message returned by hello world endpoint
	Message string `json:"message"`
}

// GetDevHelloDatasetParams defines parameters for GetDevHelloDataset.
type GetDevHelloDatasetParams struct {
	Rows    *[]string `json:"rows,omitempty"`
	Cols    *[]string `json:"cols,omitempty"`
	Bbox    *string   `json:"bbox,omitempty"`
	Geotype *string   `json:"geotype,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// query census
	// (GET /dev/hello/{dataset})
	GetDevHelloDataset(w http.ResponseWriter, r *http.Request, dataset string, params GetDevHelloDatasetParams)
	// Hello world endpoint
	// (GET /hello)
	GetHello(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// GetDevHelloDataset operation middleware
func (siw *ServerInterfaceWrapper) GetDevHelloDataset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "dataset" -------------
	var dataset string

	err = runtime.BindStyledParameter("simple", false, "dataset", chi.URLParam(r, "dataset"), &dataset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter dataset: %s", err), http.StatusBadRequest)
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetDevHelloDatasetParams

	// ------------- Optional query parameter "rows" -------------
	if paramValue := r.URL.Query().Get("rows"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "rows", r.URL.Query(), &params.Rows)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter rows: %s", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "cols" -------------
	if paramValue := r.URL.Query().Get("cols"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "cols", r.URL.Query(), &params.Cols)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter cols: %s", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "bbox" -------------
	if paramValue := r.URL.Query().Get("bbox"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "bbox", r.URL.Query(), &params.Bbox)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter bbox: %s", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "geotype" -------------
	if paramValue := r.URL.Query().Get("geotype"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "geotype", r.URL.Query(), &params.Geotype)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter geotype: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetDevHelloDataset(w, r, dataset, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetHello operation middleware
func (siw *ServerInterfaceWrapper) GetHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHello(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL     string
	BaseRouter  chi.Router
	Middlewares []MiddlewareFunc
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/dev/hello/{dataset}", wrapper.GetDevHelloDataset)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/hello", wrapper.GetHello)
	})

	return r
}
