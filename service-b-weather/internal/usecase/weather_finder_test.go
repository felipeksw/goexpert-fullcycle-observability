package usecase

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	mockAPIKey              = "your_api_key"
	mockLocale              = "Sao%20Paulo"
	mockInvalidLocale       = "Invalid%20Locale"
	mockWeatherResponseBody = `{"location": {"name": "Sao Paulo"}, "current": {}}`
)

func TestWeatherFinder_Execute_Success(t *testing.T) {
	mockContext := new(MockContext)
	mockRoundTripper := new(MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(mockWeatherResponseBody))),
	}, nil)

	origAPIKey := os.Getenv("KEY_WEATHER_API")
	_ = os.Setenv("KEY_WEATHER_API", mockAPIKey)
	defer func(key, value string) {
		_ = os.Setenv(key, value)
	}("KEY_WEATHER_API", origAPIKey)

	finder := NewWeatherFinder(mockClient)
	output, err := finder.Execute(mockContext, mockLocale)

	assert.Nil(t, err)
	assert.NotNil(t, output)
}

func TestWeatherFinder_Execute_HttpClientError(t *testing.T) {
	mockContext := new(MockContext)
	mockRoundTripper := new(MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(nil, errors.New("mocked http client error"))

	origAPIKey := os.Getenv("KEY_WEATHER_API")
	_ = os.Setenv("KEY_WEATHER_API", mockAPIKey)
	defer func(key, value string) {
		_ = os.Setenv(key, value)
	}("KEY_WEATHER_API", origAPIKey)

	finder := NewWeatherFinder(mockClient)
	_, err := finder.Execute(mockContext, mockLocale)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "mocked http client error")
}

func TestWeatherFinder_Execute_ApiKeyNotProvided(t *testing.T) {
	mockContext := new(MockContext)
	mockRoundTripper := new(MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(bytes.NewReader([]byte(""))),
	}, nil)

	origAPIKey := os.Getenv("KEY_WEATHER_API")
	_ = os.Setenv("KEY_WEATHER_API", mockAPIKey)
	defer func(key, value string) {
		_ = os.Setenv(key, value)
	}("KEY_WEATHER_API", origAPIKey)

	finder := NewWeatherFinder(mockClient)
	_, err := finder.Execute(mockContext, mockLocale)

	assert.NotNil(t, err)
	assert.Equal(t, "API key is not provided", err.Error())
}

func TestWeatherFinder_Execute_InvalidApiKey(t *testing.T) {
	mockContext := new(MockContext)
	mockRoundTripper := new(MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(bytes.NewReader([]byte(""))),
	}, nil)

	origAPIKey := os.Getenv("KEY_WEATHER_API")
	_ = os.Setenv("KEY_WEATHER_API", mockAPIKey)
	defer func(key, value string) {
		_ = os.Setenv(key, value)
	}("KEY_WEATHER_API", origAPIKey)

	finder := NewWeatherFinder(mockClient)
	_, err := finder.Execute(mockContext, mockLocale)

	assert.NotNil(t, err)
	assert.Equal(t, "API key is invalid", err.Error())
}

func TestWeatherFinder_Execute_NotFound(t *testing.T) {
	mockContext := new(MockContext)
	mockRoundTripper := new(MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"error": {"code": 1006, "message": "No matching location found."}}`))),
	}, nil)

	origAPIKey := os.Getenv("KEY_WEATHER_API")
	_ = os.Setenv("KEY_WEATHER_API", mockAPIKey)
	defer func(key, value string) {
		_ = os.Setenv(key, value)
	}("KEY_WEATHER_API", origAPIKey)

	finder := NewWeatherFinder(mockClient)
	_, err := finder.Execute(mockContext, mockInvalidLocale)

	assert.NotNil(t, err)
	assert.Equal(t, "can not find zipcode", err.Error())
}
