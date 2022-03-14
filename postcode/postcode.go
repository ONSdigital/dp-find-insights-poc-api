package postcode

import (
	"errors"
	"regexp"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"gorm.io/gorm"
)

type Postcode struct {
	gdb *gorm.DB
}

func New(gdb *gorm.DB) *Postcode {
	return &Postcode{gdb: gdb}
}

func (p *Postcode) GetMSOA(s string) (code, name string, err error) {
	s = strings.TrimSpace(s)
	s, err = normalisePostcode(s)
	if err != nil {
		return
	}

	var pc model.PostCode
	if err := p.gdb.Preload("Geo").Where(&model.PostCode{Pcds: s}).First(&pc).Error; err != nil {
		return "", "", err
	}
	return pc.Geo.Code, pc.Geo.Name, err

}

func normalisePostcode(s string) (string, error) {
	s = strings.ToUpper(s)
	re := regexp.MustCompile(`(\S*)\s*(\d\D\D)`)
	match := re.FindStringSubmatch(s)

	if len(match) != 3 || match[0] != s {
		return "", errors.New("invalid format")
	}

	return re.ReplaceAllString(s, "$1 $2"), nil
}
