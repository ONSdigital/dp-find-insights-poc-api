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
type Categories []Triplet

// Error defines model for Error.
type Error struct {
	// error message
	Error string `json:"error"`
}

// Metadata defines model for Metadata.
type Metadata struct {
	Code   *string `json:"code,omitempty"`
	Name   *string `json:"name,omitempty"`
	Slug   *string `json:"slug,omitempty"`
	Tables *Tables `json:"tables,omitempty"`
}

// MetadataResponse defines model for MetadataResponse.
type MetadataResponse []Metadata

// Table defines model for Table.
type Table struct {
	Categories *Categories `json:"categories,omitempty"`
	Code       *string     `json:"code,omitempty"`
	Name       *string     `json:"name,omitempty"`
	Slug       *string     `json:"slug,omitempty"`
}

// Tables defines model for Tables.
type Tables []Table

// Triplet defines model for Triplet.
type Triplet struct {
	Code *string `json:"code,omitempty"`
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
}

// GetCkmeansYearParams defines parameters for GetCkmeansYear.
type GetCkmeansYearParams struct {
	Cat     *string `json:"cat,omitempty"`
	Geotype *string `json:"geotype,omitempty"`
	K       *int    `json:"k,omitempty"`
}

// GetCkmeansratioYearParams defines parameters for GetCkmeansratioYear.
type GetCkmeansratioYearParams struct {
	Cat1    *string `json:"cat1,omitempty"`
	Cat2    *string `json:"cat2,omitempty"`
	Geotype *string `json:"geotype,omitempty"`
	K       *int    `json:"k,omitempty"`
}

// GetQueryYearParams defines parameters for GetQueryYear.
type GetQueryYearParams struct {
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
	// (GET /ckmeans/{year})
	GetCkmeansYear(w http.ResponseWriter, r *http.Request, year int, params GetCkmeansYearParams)
	// calculate ckmeans for the ratio between two given categories (cat1 / cat2) for a given geography type
	// (GET /ckmeansratio/{year})
	GetCkmeansratioYear(w http.ResponseWriter, r *http.Request, year int, params GetCkmeansratioYearParams)
	// Get Metadata
	// (GET /metadata)
	GetMetadata(w http.ResponseWriter, r *http.Request)
	// query census
	// (GET /query/{year})
	GetQueryYear(w http.ResponseWriter, r *http.Request, year int, params GetQueryYearParams)
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

// GetCkmeansYear operation middleware
func (siw *ServerInterfaceWrapper) GetCkmeansYear(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "year" -------------
	var year int

	err = runtime.BindStyledParameter("simple", false, "year", chi.URLParam(r, "year"), &year)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter year: %s", err), http.StatusBadRequest)
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetCkmeansYearParams

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
		siw.Handler.GetCkmeansYear(w, r, year, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetCkmeansratioYear operation middleware
func (siw *ServerInterfaceWrapper) GetCkmeansratioYear(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "year" -------------
	var year int

	err = runtime.BindStyledParameter("simple", false, "year", chi.URLParam(r, "year"), &year)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter year: %s", err), http.StatusBadRequest)
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetCkmeansratioYearParams

	// ------------- Optional query parameter "cat1" -------------
	if paramValue := r.URL.Query().Get("cat1"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "cat1", r.URL.Query(), &params.Cat1)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter cat1: %s", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "cat2" -------------
	if paramValue := r.URL.Query().Get("cat2"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "cat2", r.URL.Query(), &params.Cat2)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter cat2: %s", err), http.StatusBadRequest)
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
		siw.Handler.GetCkmeansratioYear(w, r, year, params)
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

// GetQueryYear operation middleware
func (siw *ServerInterfaceWrapper) GetQueryYear(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "year" -------------
	var year int

	err = runtime.BindStyledParameter("simple", false, "year", chi.URLParam(r, "year"), &year)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter year: %s", err), http.StatusBadRequest)
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetQueryYearParams

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
		siw.Handler.GetQueryYear(w, r, year, params)
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
		r.Get(options.BaseURL+"/ckmeans/{year}", wrapper.GetCkmeansYear)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/ckmeansratio/{year}", wrapper.GetCkmeansratioYear)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/metadata", wrapper.GetMetadata)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/query/{year}", wrapper.GetQueryYear)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/swagger", wrapper.GetSwagger)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/swaggerui", wrapper.GetSwaggerui)
	})

	return r
}
