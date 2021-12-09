package geodata

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
		return nil, errors.New("polygon must have even number of coordinates")
	}

	var points []Point
	for i := 0; i < len(a); i += 2 {
		lon, err := strconv.ParseFloat(a[i], 64)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", a[i], err)
		}
		lat, err := strconv.ParseFloat(a[i+1], 64)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", a[i+1], err)
		}

		p := Point{
			Lon: lon,
			Lat: lat,
		}
		points = append(points, p)
	}

	if len(points) < 4 {
		return nil, errors.New("must be at least 4 points in polygon")
	}
	if points[len(points)-1] != points[0] {
		return nil, errors.New("first and last point must match")
	}

	return points, nil
}
