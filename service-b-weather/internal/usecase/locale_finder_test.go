package usecase

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/felipeksw/goexpert-fullcycle-observability/service-b-weather/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	mockCep                = "01001000"
	mockInvalidCep         = "InvalidLocale"
	mockLocaleResponseBody = `{"cep": "01001000", "logradouro": "Avenida Paulista", "bairro": "Jardins", "localidade": "São Paulo", "uf": "SP"}`
)

func TestLocaleFinder_Execute_Success(t *testing.T) {
	mockContext := new(MockContext)
	mockRoundTripper := new(MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	expectedBody := []byte(mockLocaleResponseBody)
	expectedOutput := &dto.LocaleOutput{Localidade: "São Paulo"}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(expectedBody)),
	}, nil)

	finder := NewLocaleFinder(mockClient)
	output, err := finder.Execute(mockContext, mockCep)

	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestLocaleFinder_Execute_HttpClientError(t *testing.T) {
	mockContext := new(MockContext)
	mockRoundTripper := new(MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(nil, errors.New("mocked http client error"))

	finder := NewLocaleFinder(mockClient)
	_, err := finder.Execute(mockContext, mockCep)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "mocked http client error")
}

func TestLocaleFinder_Execute_InvalidResponse(t *testing.T) {
	mockContext := new(MockContext)
	mockRoundTripper := new(MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	expectedOutput := &dto.LocaleOutput{}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"erro" : true}`))),
	}, nil)

	finder := NewLocaleFinder(mockClient)
	output, err := finder.Execute(mockContext, mockInvalidCep)

	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, output)
}
