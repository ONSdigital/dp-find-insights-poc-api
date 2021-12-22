package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/gosimple/slug"
)

/*

[
  {
    "code": "QS1",
    "name": "Population Basics",
    "slug": "population-basics",
    "tables": [
      {
        "code": "QS101EW",
        "name": "Residence type",
        "slug": "residence-type",
        "categories": [
          {
            "code": "QS101EW001",
            "name": "Lives in a household",
            "slug": "lives-in-a-household"
          },
          {
            "code": "QS101EW002",
            "name": "Lives in a communal establishment",
            "slug": "lives-in-a-communal-establishment"
          },
          {
            "code": "QS101EW003",
            "name": "Communal establishments with persons sleeping rough identified",
            "slug": "communal-establishments-with-persons-sleeping-rough-identified"
          }
        ]
      }
    ]
  }
]

*/

func TestStruct(t *testing.T) {

	mdr := MetadataResponse{
		Code: spointer("QS1"),
		Name: spointer("Population Basics"),
		Slug: spointer(slug.Make("Population Basics")),
		Tables: &Tables{
			//Code: spointer("QS101EW"), // table & cats wrong way round
			// tables isn't a slice!
			Name: spointer("Residence type"),
			Slug: spointer(slug.Make("Residence type")),
			// Code: spointer("QQQ"),
			Categories: &Categories{
				Metadata{Code: spointer("QS101EW001"), Name: spointer("Lives in a household"), Slug: spointer(slug.Make("Lives in a household"))}}},
		//Metadata{Code: spointer("QS101EW001"), Name: spointer("Lives in a household"), Slug: spointer("slug")}}},
	}

	b, err := json.Marshal(&mdr)
	if err != nil {
		log.Print(err)
	}

	var out bytes.Buffer
	if err := json.Indent(&out, b, "  ", "  "); err != nil {
		log.Print(err)
	} else {
		fmt.Println(out.String())
	}

	if string(b) != `{"code":"QS1","name":"Population Basics","slug":"population-basics","tables":{"categories":[{"code":"QS1","name":"Lives in a household","slug":"lives-in-a-household"}],"name":"Residence type","slug":"residence-type"}}` {
		t.Error(string(b))
	}

}

func spointer(s string) *string {
	return &s
}
