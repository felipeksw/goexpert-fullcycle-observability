package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipeksw/goexpert-fullcycle-observability/service-b-weather/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockFinder struct {
	mock.Mock
}

func (m *MockFinder) Execute(ctx context.Context, str string) (interface{}, error) {
	args := m.Called(str)
	if args.Get(0) != nil {
		return args.Get(0).(interface{}), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestWeatherHandler_Handle(t *testing.T) {
	tests := []struct {
		name                string
		cep                 string
		mockLocaleResponse  *dto.LocaleOutput
		mockLocaleError     error
		mockWeatherResponse *dto.WeatherOutput
		mockWeatherError    error
		expectedStatusCode  int
		expectedResponse    interface{}
	}{
		{
			name:               "invalid cep length - number < 8",
			cep:                "123",
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   dto.ErrorOutput{StatusCode: http.StatusUnprocessableEntity, Message: "invalid zipcode"},
		},
		{
			name:               "invalid cep length - number > 8",
			cep:                "123456789",
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   dto.ErrorOutput{StatusCode: http.StatusUnprocessableEntity, Message: "invalid zipcode"},
		},
		{
			name:               "locale finder error",
			cep:                "12345678",
			mockLocaleError:    errors.New("locale finder error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   dto.ErrorOutput{StatusCode: http.StatusInternalServerError, Message: "locale finder error"},
		},
		{
			name:               "locale not found",
			cep:                "12345678",
			mockLocaleResponse: &dto.LocaleOutput{Localidade: ""},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   dto.ErrorOutput{StatusCode: http.StatusNotFound, Message: "can not find zipcode"},
		},
		{
			name:               "weather finder error - unauthorized",
			cep:                "12345678",
			mockLocaleResponse: &dto.LocaleOutput{Localidade: "Localidade"},
			mockWeatherError:   errors.New("API key is invalid"),
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   dto.ErrorOutput{StatusCode: http.StatusUnauthorized, Message: "API key is invalid"},
		},
		{
			name:               "weather finder error - unauthorized",
			cep:                "12345678",
			mockLocaleResponse: &dto.LocaleOutput{Localidade: "Localidade"},
			mockWeatherError:   errors.New("API key is not provided"),
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   dto.ErrorOutput{StatusCode: http.StatusUnauthorized, Message: "API key is not provided"},
		},
		{
			name:               "weather finder error - not found",
			cep:                "99999999",
			mockLocaleResponse: &dto.LocaleOutput{Localidade: "Localidade"},
			mockWeatherError:   errors.New("can not find zipcode"),
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   dto.ErrorOutput{StatusCode: http.StatusNotFound, Message: "can not find zipcode"},
		},
		{
			name:               "weather finder error - internal server error",
			cep:                "12345678",
			mockLocaleResponse: &dto.LocaleOutput{Localidade: "Localidade"},
			mockWeatherError:   errors.New("weather finder error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   dto.ErrorOutput{StatusCode: http.StatusInternalServerError, Message: "weather finder error"},
		},
		{
			name:                "successful response",
			cep:                 "12345678",
			mockLocaleResponse:  &dto.LocaleOutput{Localidade: "Localidade"},
			mockWeatherResponse: &dto.WeatherOutput{Current: dto.CurrentWeather{TempC: 25.0, TempF: 77.0}},
			expectedStatusCode:  http.StatusOK,
			expectedResponse:    dto.ResultOutput{City: "Localidade", TempC: 25.0, TempF: 77.0, TempK: 298.15},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLocaleFinder := new(MockFinder)
			mockWeatherFinder := new(MockFinder)
			handler := NewWeatherHandler(mockWeatherFinder, mockLocaleFinder)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.cep), nil)
			w := httptest.NewRecorder()

			req.SetPathValue("cep", tt.cep)

			mockLocaleFinder.On("Execute", tt.cep).Return(tt.mockLocaleResponse, tt.mockLocaleError)
			if tt.mockLocaleResponse != nil && tt.mockLocaleResponse.Localidade != "" {
				mockWeatherFinder.On("Execute", tt.mockLocaleResponse.Localidade).Return(tt.mockWeatherResponse, tt.mockWeatherError)
			}

			handler.Handle(w, req)

			res := w.Result()
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)

			var actualResponse interface{}
			if res.StatusCode == http.StatusOK {
				var resultOutput dto.ResultOutput
				err := json.NewDecoder(res.Body).Decode(&resultOutput)
				require.NoError(t, err)
				actualResponse = resultOutput
			} else {
				var errorOutput dto.ErrorOutput
				err := json.NewDecoder(res.Body).Decode(&errorOutput)
				require.NoError(t, err)
				actualResponse = errorOutput
			}

			assert.Equal(t, tt.expectedResponse, actualResponse)
		})
	}
}
