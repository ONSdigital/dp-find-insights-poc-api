// +build comparison

package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/jackc/pgx/v4"
)

/*
Not run from unit tests due to use of build tag

Ad hoc tool to compare a database generated by the
https://github.com/ONSdigital/nomis-bulk-to-postgres/ progress (as defined by
the PG_ env) with a second database created by the go process and available on
hard coded DSN below.
*/

func TestData(t *testing.T) {

	dsn := database.GetDSN()
	fmt.Printf("%#v\n", dsn)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	conn2, err := pgx.Connect(context.Background(), "postgres://insights:insights@localhost:54322/censusgo2i")
	if err != nil {
		log.Fatal(err)
	}
	defer conn2.Close(context.Background())

	for i := 0; i < 1000; i++ {
		nomis := getNomis(conn)
		geocode := getGeoCode(conn)

		println(nomis)
		println(geocode)

		if strings.Contains(geocode, "best") {
			geocode = strings.ReplaceAll(geocode, "best_fit_", "")
		}

		var geometric float64
		if err := conn.QueryRow(context.Background(), `
SELECT gm.metric FROM geo_metric gm, nomis_category c, geo g
WHERE gm.geo_id = g.id and gm.category_id=c.id
AND c.long_nomis_code=$1 and g.code=$2
	`, nomis, geocode).Scan(&geometric); err != nil {
			log.Print(err)
		}

		fmt.Printf("geo=%f\n", geometric)

		var geometric2 float64
		if err := conn2.QueryRow(context.Background(), `
		SELECT gm.metric FROM geo_metric gm, nomis_category c, geo g
		WHERE gm.geo_id = g.id and gm.category_id=c.id
		AND c.long_nomis_code=$1 and g.code=$2
			`, nomis, geocode).Scan(&geometric2); err != nil {
			log.Print(err)
		}

		fmt.Printf("geo2=%f\n", geometric2)

		if geometric != geometric2 {
			t.Fail()
		}

	}
}

func getNomis(conn *pgx.Conn) (code string) {
	if err := conn.QueryRow(context.Background(), "select long_nomis_code from nomis_category order by random() limit 1").Scan(&code); err != nil {
		log.Fatal(err)
	}

	return code
}

func getGeoCode(conn *pgx.Conn) (code string) {
	if err := conn.QueryRow(context.Background(), "select code from geo order by random() limit 1").Scan(&code); err != nil {
		log.Fatal(err)
	}

	return code
}