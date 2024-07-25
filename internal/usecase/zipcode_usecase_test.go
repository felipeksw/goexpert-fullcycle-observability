package usecase_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/entity"
	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/infra/mockup"
	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
)

func TestNewAddressByZipcodeSuccess(t *testing.T) {

	mockCepSuccess := "01001000"
	mockZipcodeResponseSuccessBody := `{"cep":"01001-000","localidade":"São Paulo","erro":""}`

	mockRoundTripper := new(mockup.MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(mockZipcodeResponseSuccessBody))),
	}, nil)

	tracer := otel.Tracer("test")

	mockZipcodeDto, err := entity.NewZipcode(mockCepSuccess)
	assert.Nil(t, err)

	addressDto, err := usecase.NewAddressByZipcode(context.Background(), tracer, *mockZipcodeDto, mockClient)
	assert.Nil(t, err)

	add, err := json.Marshal(addressDto)
	assert.Nil(t, err)

	assert.Equal(t, []byte(mockZipcodeResponseSuccessBody), add)
}

func TestNewAddressByZipcodeCepNotFound(t *testing.T) {

	mockCepSuccess := "01001009"
	mockZipcodeResponseErrorBody := `{"erro":"true"}`

	mockRoundTripper := new(mockup.MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(mockZipcodeResponseErrorBody))),
	}, nil)

	tracer := otel.Tracer("test")

	mockZipcodeDto, err := entity.NewZipcode(mockCepSuccess)
	assert.Nil(t, err)

	addressDto, err := usecase.NewAddressByZipcode(context.Background(), tracer, *mockZipcodeDto, mockClient)
	assert.Equal(t, "zip code not found", err.Error())
	assert.Nil(t, addressDto)
}
