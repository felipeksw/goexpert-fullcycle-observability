package usecase_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/dto"
	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/infra/mockup"
	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewWeatherByAddressSuccess(t *testing.T) {

	mockWeatherResponseSuccessBody := `{"Location":{"name":"San Paulo","region":"Sao Paulo"},"Current":{"temp_c":14.2}}`
	mockAddressSuccess := dto.AddressDto{
		Cep:        "01001-000",
		Localidade: "São Paulo",
		Error:      "",
	}

	mockRoundTripper := new(mockup.MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(mockWeatherResponseSuccessBody))),
	}, nil)

	weatherDto, err := usecase.NewWeatherByAddress(context.Background(), mockAddressSuccess, mockClient)
	assert.Nil(t, err)

	wea, err := json.Marshal(weatherDto)
	assert.Nil(t, err)

	assert.Equal(t, []byte(mockWeatherResponseSuccessBody), wea)
}

func TestNewWeatherByAddressAddressNotFount(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	mockWeatherResponseErrorBody := "api.weatherapi.com: "

	slog.Info("[test code 400]", "status", http.StatusText(400))

	mockAddressError := dto.AddressDto{
		Cep:        "01001-000",
		Localidade: "SãoPauloo",
		Error:      "",
	}

	mockRoundTripper := new(mockup.MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusBadRequest,
		Status:     http.StatusText(http.StatusBadRequest),
		Body:       io.NopCloser(bytes.NewReader([]byte(mockWeatherResponseErrorBody))),
	}, nil)

	weatherDto, err := usecase.NewWeatherByAddress(context.Background(), mockAddressError, mockClient)
	slog.Info("[test struct]", "weatherDto", weatherDto)

	assert.Nil(t, weatherDto)

	slog.Info("[test NewWeatherByAddress]", "error", err.Error())

	assert.Contains(t, err.Error(), "Bad Request")
}
