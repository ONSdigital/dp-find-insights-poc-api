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

// Categories defines model for Categories.
type Categories []Metadata

// Error defines model for Error.
type Error struct {
	// error message
	Error string `json:"error"`
}

// GetDevCkmeansParams defines parameters for GetDevCkmeans.
type GetDevCkmeansParams struct {
	Cat     *string `json:"cat,omitempty"`
	Geotype *string `json:"geotype,omitempty"`
	K       *int    `json:"k,omitempty"`
}

// Metadata defines model for Metadata.
type Metadata struct {
	Code *string `json:"code,omitempty"`
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
}

// MetadataResponse defines model for MetadataResponse.
type MetadataResponse struct {
	Code   *string `json:"code,omitempty"`
	Name   *string `json:"name,omitempty"`
	Slug   *string `json:"slug,omitempty"`
	Tables *Tables `json:"tables,omitempty"`
}

// Tables defines model for Tables.
type Tables struct {
	Categories *Categories `json:"categories,omitempty"`
	Name       *string     `json:"name,omitempty"`
	Slug       *string     `json:"slug,omitempty"`
}

// GetDevHelloDatasetParams defines parameters for GetDevHelloDataset.
type GetDevHelloDatasetParams struct {
	Rows        *[]string `json:"rows,omitempty"`
	Cols        *[]string `json:"cols,omitempty"`
	Bbox        *string   `json:"bbox,omitempty"`
	Geotype     *[]string `json:"geotype,omitempty"`
	Location    *string   `json:"location,omitempty"`
	Radius      *int      `json:"radius,omitempty"`
	Polygon     *string   `json:"polygon,omitempty"`
	Censustable *string   `json:"censustable,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// calculate ckmeans over a given category and geography type
	// (GET /dev/ckmeans)
	GetDevCkmeans(w http.ResponseWriter, r *http.Request, params GetDevCkmeansParams)
	// query census
	// (GET /dev/hello/{dataset})
	GetDevHelloDataset(w http.ResponseWriter, r *http.Request, dataset string, params GetDevHelloDatasetParams)
	// Get Metadata
	// (GET /metadata)
	GetMetadata(w http.ResponseWriter, r *http.Request)
	// spec
	// (GET /swagger)
	GetSwagger(w http.ResponseWriter, r *http.Request)
	// spec
	// (GET /swaggerui)
	GetSwaggerui(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// GetDevCkmeans operation middleware
func (siw *ServerInterfaceWrapper) GetDevCkmeans(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetDevCkmeansParams

	// ------------- Optional query parameter "cat" -------------
	if paramValue := r.URL.Query().Get("cat"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "cat", r.URL.Query(), &params.Cat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter cat: %s", err), http.StatusBadRequest)
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

	// ------------- Optional query parameter "k" -------------
	if paramValue := r.URL.Query().Get("k"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "k", r.URL.Query(), &params.K)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter k: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetDevCkmeans(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

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

	// ------------- Optional query parameter "location" -------------
	if paramValue := r.URL.Query().Get("location"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "location", r.URL.Query(), &params.Location)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter location: %s", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "radius" -------------
	if paramValue := r.URL.Query().Get("radius"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "radius", r.URL.Query(), &params.Radius)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter radius: %s", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "polygon" -------------
	if paramValue := r.URL.Query().Get("polygon"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "polygon", r.URL.Query(), &params.Polygon)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter polygon: %s", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "censustable" -------------
	if paramValue := r.URL.Query().Get("censustable"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "censustable", r.URL.Query(), &params.Censustable)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter censustable: %s", err), http.StatusBadRequest)
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

// GetMetadata operation middleware
func (siw *ServerInterfaceWrapper) GetMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetMetadata(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetSwagger operation middleware
func (siw *ServerInterfaceWrapper) GetSwagger(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSwagger(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetSwaggerui operation middleware
func (siw *ServerInterfaceWrapper) GetSwaggerui(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSwaggerui(w, r)
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
		r.Get(options.BaseURL+"/dev/ckmeans", wrapper.GetDevCkmeans)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/dev/hello/{dataset}", wrapper.GetDevHelloDataset)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/metadata", wrapper.GetMetadata)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/swagger", wrapper.GetSwagger)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/swaggerui", wrapper.GetSwaggerui)
	})

	return r
}
