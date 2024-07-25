package entity_test

import (
	"errors"
	"testing"

	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewZipcode(t *testing.T) {

	type zipcodeLote struct {
		zipcode string
		err     error
		status  bool
	}

	err8Digits := errors.New("zip code must be 8 numeric digits")

	table := []zipcodeLote{
		{"", err8Digits, false},
		{"1300000z", err8Digits, false},
		{"130000010", err8Digits, false},
		{"13000001012345678", err8Digits, false},
		{"#30000010", err8Digits, false},
		{"13000001$", err8Digits, false},
		{"^#3000001", err8Digits, false},
		{"'3000001", err8Digits, false},
		{"\"3000001", err8Digits, false},
		{"130000-010", err8Digits, false},
		{"ABCDEFGH", err8Digits, false},
		{"13000001", nil, true},
		{"00000000", nil, true},
		{"99999999", nil, true},
	}
	for _, item := range table {
		zipcodeDto, err := entity.NewZipcode(item.zipcode)
		if item.status {
			assert.Nil(t, err)
			assert.Equal(t, item.zipcode, zipcodeDto.Zipcode)
		} else {
			assert.Error(t, err, item.err)
		}
	}
}
