package entity_test

import (
	"errors"
	"log/slog"
	"math"
	"strings"
	"testing"

	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewLocaleWeather(t *testing.T) {

	type localeWeatherLote struct {
		locale string
		temp   float64
		err    error
	}
	errTemp := errors.New("temperature is outside the earth range")
	errLoca := errors.New("location can not be empty")

	table := []localeWeatherLote{
		{"  ", 58.001, errLoca},
		{"", -89.001, errLoca},
		{"", 22.3, errLoca},
		{" ", 17.9, errLoca},
		{"Osasco", 58.001, errTemp},
		{"Penha", -89.001, errTemp},
		{"Ribeirão Preto", 57.999, nil},
		{"Toronto", -88.999, nil},
		{"San Luis", 0, nil},
		{"São Paulo", 22, nil},
		{"Campinas", 24.5, nil},
	}
	for _, item := range table {
		localeWeatherDto, err := entity.NewLocaleWeather(item.locale, item.temp)

		if item.err != nil {
			slog.Info("[err != nil]", "locale", item.locale, "temp", item.temp, "erro", item.err.Error())

			assert.Error(t, err, item.err)
		} else {
			slog.Info("[err == nil]", "locale", item.locale, "temp", item.temp)

			assert.Nil(t, err)
			assert.Equal(t, strings.TrimSpace(item.locale), localeWeatherDto.Locale)
			assert.Equal(t, math.Round((item.temp)*10)/10, math.Round((localeWeatherDto.TempC)*10)/10)
			assert.Equal(t, math.Round((item.temp*1.8+32)*10)/10, math.Round((localeWeatherDto.TempF)*10)/10)
			assert.Equal(t, math.Round((item.temp+273)*10)/10, math.Round((localeWeatherDto.TempK)*10)/10)
		}
	}
}
