package entity

import (
	"errors"
	"log/slog"
	"math"
	"strings"

	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/dto"
)

type localWeatherEntity struct {
	locale string
	tempC  float64
	tempF  float64
	tempK  float64
}

func (w *localWeatherEntity) TempC() float64 {
	return w.tempC
}

func (w *localWeatherEntity) Locale() string {
	return w.locale
}

func NewLocaleWeather(locale string, tempC float64) (*dto.LocalWeatherDto, error) {

	locale = strings.TrimSpace(locale)

	var tc = &localWeatherEntity{
		locale: locale,
		tempC:  tempC,
		tempF:  0,
		tempK:  0,
	}
	slog.Debug("struct", "localWeatherEntity", tc)

	err := tc.IsValid()
	if err != nil {
		slog.Error("[invalid locale]", "error", err.Error())
		return nil, err
	}

	return &dto.LocalWeatherDto{
		Locale: tc.locale,
		TempC:  math.Round((tc.tempC)*10) / 10,
		TempF:  math.Round((tc.tempC*1.8+32)*10) / 10,
		TempK:  math.Round((tc.tempC+273)*10) / 10,
	}, nil
}

func (z *localWeatherEntity) IsValid() error {

	if len(z.Locale()) < 1 {
		return errors.New("location can not be empty")
	}

	if z.TempC() > 58 || z.tempC < -89 {
		return errors.New("temperature is outside the earth range")
	}
	return nil
}
