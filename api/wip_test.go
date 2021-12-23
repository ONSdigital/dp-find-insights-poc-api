package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/gosimple/slug"
	"github.com/ryboe/q"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*

OLD
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

New
[
    {
      "code": "QS1",
      "name": "Population Basics",
      "slug": "population-basics",
      "tables": [
        {
          "categories": [
            {
              "code": "QS101EW0001",
              "name": "All categories: Residence type",
              "slug": "all-categories-residence-type"
            },
            {
              "code": "QS101EW0002",
              "name": "Lives in a household",
              "slug": "lives-in-a-household"
            },
            {
              "code": "QS101EW0003",
              "name": "Lives in a communal establishment",
              "slug": "lives-in-a-communal-establishment"
            },
            {
              "code": "QS101EW0004",
              "name": "Communal establishments with persons sleeping rough identified",
              "slug": "communal-establishments-with-persons-sleeping-rough-identified"
            }
          ],
          "code": "QS101EW",
          "name": "Residence type",
          "slug": "residence-type"
        }
      ]
    }
  ]
*/

// UPDATE nomis_desc set nomis_topic_id=1 WHERE short_nomis_code like 'QS1%';
func TestReadSomeData(t *testing.T) {

	db, err := gorm.Open(postgres.Open(database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var topics []model.NomisTopic
	db.Preload("NomisDescs").First(&topics)
	topic := topics[0]

	var nd model.NomisDesc

	for _, nd = range topic.NomisDescs {
		if nd.ShortNomisCode == "QS101EW" {
			db.Preload("NomisCategories").Find(&nd)
			break
		}
	}

	var cats Categories
	for _, trip := range nd.NomisCategories {
		cats = append(cats, Triplet{Code: spointer(trip.LongNomisCode), Name: spointer(trip.CategoryName), Slug: spointer(slug.Make(trip.CategoryName))})
	}

	mdr := MetadataResponse{Metadata{
		Code: spointer(topic.TopNomisCode),
		Name: spointer(topic.Name),
		Slug: spointer(slug.Make(topic.Name)),
		Tables: &Tables{{
			Name:       spointer(nd.Name),
			Slug:       spointer(slug.Make(nd.Name)),
			Code:       spointer(nd.ShortNomisCode),
			Categories: &cats,
		},
		},
	},
	}

	b := jmarsh(mdr)

	if string(b) != `[{"code":"QS1","name":"Population Basics","slug":"population-basics","tables":[{"categories":[{"code":"QS101EW0001","name":"All categories: Residence type","slug":"all-categories-residence-type"},{"code":"QS101EW0002","name":"Lives in a household","slug":"lives-in-a-household"},{"code":"QS101EW0003","name":"Lives in a communal establishment","slug":"lives-in-a-communal-establishment"},{"code":"QS101EW0004","name":"Communal establishments with persons sleeping rough identified","slug":"communal-establishments-with-persons-sleeping-rough-identified"}],"code":"QS101EW","name":"Residence type","slug":"residence-type"}]}]` {
		t.Error(string(b))
	}

}

func TestReadData(t *testing.T) {

	db, err := gorm.Open(postgres.Open(database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var topics []model.NomisTopic
	db.Preload("NomisDescs").Find(&topics)
	//	fmt.Printf("%#v\n", topics[0])
	q.Q(topics[0])

}

func TestStruct(t *testing.T) {

	mdr := MetadataResponse{Metadata{
		Code: spointer("QS1"),
		Name: spointer("Population Basics"),
		Slug: spointer(slug.Make("Population Basics")),
		Tables: &Tables{{
			Name: spointer("Residence type"),
			Slug: spointer(slug.Make("Residence type")),
			Code: spointer("QS101EW"),
			Categories: &Categories{
				Triplet{Code: spointer("QS101EW001"), Name: spointer("Lives in a household"), Slug: spointer(slug.Make("Lives in a household"))},
				Triplet{Code: spointer("QS101EW002"), Name: spointer("Lives in a communal establishment"), Slug: spointer(slug.Make("Lives in a communal establishment"))},
			}},
		},
	},
	}

	b := jmarsh(mdr)

	if string(b) != `[{"code":"QS1","name":"Population Basics","slug":"population-basics","tables":[{"categories":[{"code":"QS101EW001","name":"Lives in a household","slug":"lives-in-a-household"},{"code":"QS101EW002","name":"Lives in a communal establishment","slug":"lives-in-a-communal-establishment"}],"code":"QS101EW","name":"Residence type","slug":"residence-type"}]}]` {
		t.Error(string(b))
	}

}

func spointer(s string) *string {
	return &s
}

func jmarsh(mdr MetadataResponse) []byte {
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

	return b
}
