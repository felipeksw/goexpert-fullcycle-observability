package entity

import (
	"errors"
	"log/slog"
	"regexp"

	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/dto"
)

type zipcodeEntity struct {
	zipcode string
}

func NewZipcode(zipcode string) (*dto.ZipcodeDto, error) {

	var zc = &zipcodeEntity{
		zipcode: zipcode,
	}

	err := zc.IsValid()
	if err != nil {
		slog.Error("[invalid zipcode]", "error", err.Error())
		return nil, err
	}

	return &dto.ZipcodeDto{
		Zipcode: zc.zipcode,
	}, nil
}

func (z *zipcodeEntity) IsValid() error {

	var re = regexp.MustCompile(`^[0-9]{8}$`)

	if !re.MatchString(z.zipcode) {
		return errors.New("invalid zipcode")
	}
	return nil
}
