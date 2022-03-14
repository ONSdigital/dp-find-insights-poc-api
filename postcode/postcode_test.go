package postcode

import (
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/cockroachdb/copyist"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	copyist.Register("postgres")
}

func TestLookupPC(t *testing.T) {

	defer copyist.Open(t).Close()

	dsn := database.GetDSN() + "?sslmode=disable"

	// use recorded DB response
	// "go test -v ../postcode -run TestLookupPC -record" to create new
	gdb, err := gorm.Open(postgres.New(postgres.Config{DriverName: "copyist_postgres", DSN: dsn}), &gorm.Config{})
	if err != nil {
		t.Errorf(err.Error())
	}

	p := New(gdb)

	{
		code, name, err := p.GetMSOA("EX39 5AA")

		if err != nil {
			t.Fail()
		}

		if code != "E02004223" {
			t.Errorf(code)
		}

		if name != "Bideford South & East" {
			t.Errorf(name)
		}
	}
	{
		code, name, err := p.GetMSOA("WS42BJ")

		if err != nil {
			t.Fail()
		}

		if code != "E02002133" {
			t.Errorf(code)
		}

		if name != "Walsall North East" {
			t.Errorf(name)
		}
	}
	{
		code, name, err := p.GetMSOA("  n16   7Hf ")

		if err != nil {
			t.Fail()
		}

		if code != "E02000350" {
			t.Errorf(code)
		}

		if name != "Stoke Newington East & Cazenove" {
			t.Errorf(name)
		}
	}

}

func TestNormalisePostcode(t *testing.T) {

	{
		_, err := normalisePostcode("rm     52dd")

		if err == nil {
			t.Fail()
		}
	}

	{
		n, _ := normalisePostcode("Rm5      2dD")

		if n != "RM5 2DD" {
			t.Errorf("got %s", n)

		}
	}

	{
		n, _ := normalisePostcode("rm5 2dD")

		if n != "RM5 2DD" {
			t.Errorf("got %s", n)

		}
	}

	{
		n, _ := normalisePostcode("RM52DD")

		if n != "RM5 2DD" {
			t.Errorf("got %s", n)

		}
	}

}
