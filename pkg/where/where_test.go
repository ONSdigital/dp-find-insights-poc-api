package where_test

import (
	"reflect"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/where"
)

func TestNoRows(t *testing.T) {
	got, err := where.ParseRows(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatal("expected empty Rows")
	}
}

func TestMissingColumnName(t *testing.T) {
	_, err := where.ParseRows([]string{"missing-colon"})
	if err == nil {
		t.Fatal("expected error about missing bad rowspec")
	}
}

func TestMissingRowSpec(t *testing.T) {
	_, err := where.ParseRows([]string{"col:"})
	if err == nil {
		t.Fatal("expected error about empty spec")
	}
}

func TestSingleColumnValue(t *testing.T) {
	want := where.Rows{
		"col": &where.ValueSet{
			Singles: []string{
				"val",
			},
		},
	}

	got, err := where.ParseRows([]string{"col:val"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatal("expected single value in single col")
	}
}

func TestMissingValueInRange(t *testing.T) {
	_, err := where.ParseRows([]string{"col:low..."})
	if err == nil {
		t.Fatal("expect error about missing column in range")
	}
}

func TestColumnRange(t *testing.T) {
	want := where.Rows{
		"col": &where.ValueSet{
			Ranges: []*where.ValueRange{
				{
					Low:  "lo",
					High: "hi",
				},
			},
		},
	}

	got, err := where.ParseRows([]string{"col:lo...hi"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatal("expected a range in col")
	}
}

func TestTooManyValuesInRange(t *testing.T) {
	_, err := where.ParseRows([]string{"col:low...med...high"})
	if err == nil {
		t.Fatal("expected error about bad range")
	}
}

func TestSingleAndRange(t *testing.T) {
	want := where.Rows{
		"col": &where.ValueSet{
			Singles: []string{
				"one",
			},
			Ranges: []*where.ValueRange{
				{
					Low:  "lo",
					High: "hi",
				},
			},
		},
	}

	got, err := where.ParseRows([]string{"col:one", "col:lo...hi"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatal("single and range failed")
	}
}

func TestMultipleColumns(t *testing.T) {
	want := where.Rows{
		"colA": &where.ValueSet{
			Singles: []string{
				"one",
			},
		},
		"colB": &where.ValueSet{
			Ranges: []*where.ValueRange{
				{
					Low:  "lo",
					High: "hi",
				},
			},
		},
	}

	got, err := where.ParseRows([]string{"colA:one", "colB:lo...hi"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatal("multiple columns failed")
	}
}

func TestEmptyWhere(t *testing.T) {
	r := where.Rows{}

	clause := where.Clause(r)
	if clause != "" {
		t.Fatal("expected empty where clause")
	}
}

func TestSingleValueWhere(t *testing.T) {
	r := where.Rows{
		"col": &where.ValueSet{
			Singles: []string{
				"val",
			},
		},
	}

	clause := where.Clause(r)
	if clause != "WHERE col IN ( 'val' )" {
		t.Fatal("expected single value where")
	}
}

func TestMultipleValueWhere(t *testing.T) {
	r := where.Rows{
		"col": &where.ValueSet{
			Singles: []string{
				"val1",
				"val2",
			},
		},
	}

	clause := where.Clause(r)
	if clause != "WHERE col IN ( 'val1','val2' )" {
		t.Fatal("expected multiple value where")
	}
}

func TestRangeWhere(t *testing.T) {
	r := where.Rows{
		"col": &where.ValueSet{
			Ranges: []*where.ValueRange{
				{
					Low:  "lo",
					High: "hi",
				},
			},
		},
	}

	clause := where.Clause(r)
	if clause != "WHERE col BETWEEN 'lo' AND 'hi'" {
		t.Fatal("expected single column range")
	}
}

func TestValueAndRangeWhere(t *testing.T) {
	r := where.Rows{
		"col": &where.ValueSet{
			Singles: []string{
				"val",
			},
			Ranges: []*where.ValueRange{
				{
					Low:  "lo",
					High: "hi",
				},
			},
		},
	}

	clause := where.Clause(r)
	if clause != "WHERE col IN ( 'val' ) OR col BETWEEN 'lo' AND 'hi'" {
		t.Fatal("expected where clause with IN and BETWEEN")
	}
}

func TestComplexWhere(t *testing.T) {

	r := where.Rows{
		"colA": &where.ValueSet{
			Singles: []string{
				"valA",
			},
			Ranges: []*where.ValueRange{
				{
					Low:  "loA",
					High: "hiA",
				},
			},
		},
		"colB": &where.ValueSet{
			Singles: []string{
				"valB",
			},
			Ranges: []*where.ValueRange{
				{
					Low:  "loB",
					High: "hiB",
				},
			},
		},
	}

	got := where.Clause(r)

	// columns are treated in a random order since they are in a range,
	// so check for either variation
	wantA := "WHERE colA IN ( 'valA' ) OR colA BETWEEN 'loA' AND 'hiA' OR colB IN ( 'valB' ) OR colB BETWEEN 'loB' AND 'hiB'"
	wantB := "WHERE colB IN ( 'valB' ) OR colB BETWEEN 'loB' AND 'hiB' OR colA IN ( 'valA' ) OR colA BETWEEN 'loA' AND 'hiA'"
	if got != wantA && got != wantB {
		t.Fatalf("got %s, want %s or %s", got, wantA, wantB)
	}
}
