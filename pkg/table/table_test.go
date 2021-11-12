package table

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewErrors(t *testing.T) {
	primary := "GEO"
	var tests = []struct {
		primary  string
		colnames []string
	}{
		{"", []string{"COLA"}},
		{primary, []string{}},
		{primary, []string{""}},
		{primary, []string{"COLA", "COLA"}},
	}

	for _, test := range tests {
		tbl, err := New(test.primary, test.colnames)
		if err == nil {
			t.Errorf("expected error")
			continue
		}
		if tbl != nil {
			t.Errorf("expected tbl to be nil")
		}
	}
}

func TestNew(t *testing.T) {
	primary := "GEO"

	var tests = []struct {
		colnames []string
		want     string
	}{
		{[]string{"COL"}, fmt.Sprintf("%s,COL\n", primary)},
		{[]string{"COLA", "COLB"}, fmt.Sprintf("%s,COLA,COLB\n", primary)},
	}

	for _, test := range tests {
		tbl, err := New(primary, test.colnames)
		if err != nil {
			t.Error(err)
			continue
		}
		buf := strings.Builder{}
		err = tbl.Generate(&buf)
		if err != nil {
			t.Error(err)
			continue
		}
		if buf.String() != test.want {
			t.Errorf("got %q, want %q", buf.String(), test.want)
			continue
		}
	}
}

func TestColIndex(t *testing.T) {
	var tests = []struct {
		colname string
		want    int
	}{
		{"", -1},
		{"noexist", -1},
		{"COLA", 1},
		{"COLB", 2},
	}

	tbl, err := New("GEO", []string{"COLA", "COLB"})
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		i := tbl.colIndex(test.colname)
		if i != test.want {
			t.Errorf("%s: %d, want %d", test.colname, i, test.want)
		}
	}
}

func TestFindRow(t *testing.T) {
	tbl, err := New("GEO", []string{"COL"})
	if err != nil {
		t.Fatal(err)
	}

	code := "GEO"
	row := tbl.findRow(code)
	if row[0] != code {
		t.Errorf("expected new row[0] to be %s", code)
	}

	if len(tbl.rows) != 1 {
		t.Errorf("expected table to have one row")
	}

	if tbl.rows[0][0] != code {
		t.Errorf("expected table[0][0] to be %s", code)
	}

	row = tbl.findRow(code)
	if row[0] != code {
		t.Errorf("expected existing row[0] to be %s", code)
	}
}

func TestSetCellError(t *testing.T) {
	tbl, err := New("GEO", []string{"COL"})
	if err != nil {
		t.Fatal(err)
	}

	err = tbl.SetCell("GEO", "WRONG", "a value")
	if err == nil {
		t.Fatalf("expected SetCell to return error")
	}
}

func TestSetCell(t *testing.T) {
	primary := "GEO"

	rows := [][]string{
		{"HERE", "COLA", "10"},
		{"HERE", "COLB", "20"},
		{"THERE", "COLA", "30"},
		{"THERE", "COLB", "40"},
	}

	tbl, err := New(primary, []string{"COLA", "COLB"})
	if err != nil {
		t.Fatal(err)
	}

	for _, row := range rows {
		err := tbl.SetCell(row[0], row[1], row[2])
		if err != nil {
			t.Fatal(err)
		}
	}

	want := fmt.Sprintf("%s,COLA,COLB\nHERE,10,20\nTHERE,30,40\n", primary)

	buf := strings.Builder{}
	err = tbl.Generate(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != want {
		t.Fatalf("got %q, want %q", buf.String(), want)
	}
}
