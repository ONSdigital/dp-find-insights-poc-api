package model

import (
	"testing"
)

func TestGetGeoTypeMap(t *testing.T) {
	g := GetGeoTypeMap()

	if !g["LAD"] {
		t.Fail()
	}

	if g["XXX"] {
		t.Fail()
	}
}
