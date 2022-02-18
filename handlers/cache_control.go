package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/cache"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
)

const (
	mimeCSV  = "text/csv"
	mimeJSON = "application/json"
)

type generateFunc func() ([]byte, error)

// respond returns cached data if it is available, or generates and caches new data.
func (svr *Server) respond(w http.ResponseWriter, r *http.Request, contentType string, generate generateFunc) {

	// add CORS header
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var err error
	var body []byte

	key := cache.CacheKey(r)

	// allocate a serialiser for this cache key
	ser := svr.cm.AllocateEntry(key)
	defer ser.Free()

	func() {
		ctx := r.Context()

		// lock cache key before doing any cache operations
		ser.Lock()
		defer ser.Unlock()

		if !noCache(r) {
			body, err = ser.Get(ctx)
			if err == nil {
				return
			}
		}

		body, err = generate()
		if err != nil {
			return
		}

		// if there is a problem saving response in cache, log it, but still send to client
		err = ser.Set(ctx, body)
		if err != nil {
			log.Printf("could not cache: %q (%d bytes): %v\n", key, len(body), err)
			err = nil
		}
	}()

	if err == nil {
		w.Header().Add("Content-Type", contentType)
		w.Write(body)
		return
	}

	code := http.StatusInternalServerError
	switch {
	case errors.Is(err, geodata.ErrMissingParams) || errors.Is(err, geodata.ErrInvalidParams):
		code = http.StatusBadRequest
	case errors.Is(err, geodata.ErrNoContent):
		code = http.StatusNoContent
	case errors.Is(err, geodata.ErrTooManyMetrics):
		code = http.StatusForbidden
	}
	sendError(w, code, err.Error())
}

// noCache is true if a Cache-Control header contains "no-cache"
// (This is just enough to let us get around caching during development.
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control)
func noCache(req *http.Request) bool {
	for _, value := range req.Header.Values("Cache-Control") {
		for _, token := range strings.Split(value, ",") {
			if strings.EqualFold(token, "no-cache") {
				return true
			}
		}
	}
	return false
}
