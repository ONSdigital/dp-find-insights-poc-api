// The database package is super thin layer on sql that provides a Checker.
//
package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/aws"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type Database struct {
	driver string
	db     *sql.DB
}

// These unique types and the Provide* functions below are
// glue to make wire work.
//
type Password string
type DSN string

func ProvidePassword(aws *aws.Clients) (Password, error) {
	pwd := os.Getenv("PGPASSWORD")
	if pwd != "" {
		return Password(pwd), nil
	}
	pwd, err := aws.GetSecret(os.Getenv("FI_PG_SECRET_ID"))
	return Password(pwd), err
}

func ProvideDSN(pwd Password) DSN {
	return DSN(GetDSN(string(pwd)))
}

func ProvideDatabase(dsn DSN) (*Database, error) {
	db, err := Open("pgx", string(dsn))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Open(driverName, dsn string) (*Database, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	return &Database{
		driver: driverName,
		db:     db,
	}, nil
}

func (db *Database) DB() *sql.DB {
	return db.db
}

func (db *Database) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	conn, err := db.db.Conn(ctx)
	if err != nil {
		state.Update(healthcheck.StatusCritical, err.Error(), 0)
		return nil
	}
	state.Update(healthcheck.StatusOK, db.driver+" healthy", 0)
	conn.Close()
	return nil
}

func (db *Database) Close() error {
	return db.db.Close()
}

func GetDSN(pw ...string) string {
	if len(pw) == 0 {
		pw = append(pw, os.Getenv("PGPASSWORD"))
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("PGUSER"),
		pw[0],
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
	)
}
