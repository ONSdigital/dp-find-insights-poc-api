package geodata

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/sentinel"
)

type Point struct {
	Lon float64
	Lat float64
}

func (c *Point) String() string {
	return fmt.Sprintf("%.13g %.13g", c.Lon, c.Lat)
}

func ParsePolygon(s string) ([]Point, error) {
	a := strings.Split(s, ",")
	if len(a)%2 != 0 {
		return nil, fmt.Errorf("%w: polygon must have even number of coordinates", sentinel.ErrInvalidParams)
	}

	var points []Point
	for i := 0; i < len(a); i += 2 {
		lon, err := strconv.ParseFloat(a[i], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %s", sentinel.ErrInvalidParams, a[i], err)
		}
		lat, err := strconv.ParseFloat(a[i+1], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %s", sentinel.ErrInvalidParams, a[i+1], err)
		}

		p := Point{
			Lon: lon,
			Lat: lat,
		}
		points = append(points, p)
	}

	if len(points) < 4 {
		return nil, fmt.Errorf("%w: must be at least 4 points in polygon", sentinel.ErrInvalidParams)
	}
	if points[len(points)-1] != points[0] {
		return nil, fmt.Errorf("%w: first and last point must match", sentinel.ErrInvalidParams)
	}

	return points, nil
}
