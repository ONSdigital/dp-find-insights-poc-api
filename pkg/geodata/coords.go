package geodata

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/sentinel"
	geom "github.com/twpayne/go-geom"
)

const maxRadius = 1000000 // largest "sane" Circle radius 1000km

// Defined UK bounding box for basic sanity checking.
var ukbbox = geom.NewBounds(geom.XY).SetCoords(
	geom.Coord{-7.57, 58.64}, // NW corner
	geom.Coord{1.76, 49.92},  // SE corner
)

// parseCoords parses a comma-separated list of floats and returns a []float64.
func parseCoords(s string) ([]float64, error) {
	coords := []float64{}
	for _, tok := range strings.Split(s, ",") {
		coord, err := strconv.ParseFloat(tok, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: error parsing coordinate %q: %s", sentinel.ErrInvalidParams, s, err)
		}
		coords = append(coords, coord)
	}
	return coords, nil
}

// checkValidCoords validates a flat coordinate slice.
// Returns nil if all points are valid geographic coordinates.
func checkValidCoords(coords []float64) error {
	if len(coords)%2 != 0 {
		return fmt.Errorf("%w: must be even number of coordinates", sentinel.ErrInvalidParams)
	}
	for i := 0; i < len(coords); i += 2 {
		if !isValidLon(coords[i]) {
			return fmt.Errorf("%w: longitude %g out of range", sentinel.ErrInvalidParams, coords[i])
		}
		if !isValidLat(coords[i+1]) {
			return fmt.Errorf("%w: latitude %g out of range", sentinel.ErrInvalidParams, coords[i+1])
		}
	}
	return nil
}

// isValidLon is true if lon is between -180 and 180
func isValidLon(lon float64) bool {
	return -180 <= lon && lon < 180
}

// isValidLat is true if lat is between -90 and 90
func isValidLat(lat float64) bool {
	return -90 <= lat && lat <= 90
}

// asLinestring takes a flat coordinate list and returns a string suitable for use
// in a SQL LINESTRING() call.
// So {1,2,3,4} returns "1 2,3 4".
// Must be an even number of elements in coords.
func asLineString(coords []float64) (string, error) {
	if len(coords)%2 != 0 {
		return "", fmt.Errorf("%w: uneven number of coords for LINESTRING", sentinel.ErrInvalidParams)
	}
	var a []string
	for i := 0; i < len(coords); i += 2 {
		s := fmt.Sprintf("%.13g %.13g", coords[i], coords[i+1])
		a = append(a, s)
	}
	return strings.Join(a, ","), nil
}
