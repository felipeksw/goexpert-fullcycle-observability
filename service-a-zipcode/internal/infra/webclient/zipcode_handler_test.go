package webclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipeksw/goexpert-fullcycle-observability/service-a-zipcode/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockLocaleFinder struct {
	mock.Mock
}

func (m *MockLocaleFinder) Execute(ctx context.Context, cep string) (*dto.LocaleOutput, error) {
	args := m.Called(cep)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.LocaleOutput), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestLocaleHandler_Handle(t *testing.T) {
	tests := []struct {
		name               string
		cep                string
		mockLocaleResponse *dto.LocaleOutput
		mockLocaleError    error
		expectedStatusCode int
		expectedResponse   interface{}
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
			name:               "successful response",
			cep:                "12345678",
			mockLocaleResponse: &dto.LocaleOutput{Localidade: "Localidade"},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   dto.LocaleOutput{Localidade: "Localidade"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLocaleFinder := new(MockLocaleFinder)
			handler := NewLocaleHandler(mockLocaleFinder)

			input := dto.LocaleInput{Cep: tt.cep}
			inputJson, _ := json.Marshal(input)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
			w := httptest.NewRecorder()

			mockLocaleFinder.On("Execute", tt.cep).Return(tt.mockLocaleResponse, tt.mockLocaleError)

			handler.Handle(w, req)

			res := w.Result()
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)

			var actualResponse interface{}
			if res.StatusCode == http.StatusOK {
				var output dto.LocaleOutput
				err := json.NewDecoder(res.Body).Decode(&output)
				require.NoError(t, err)
				actualResponse = output
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
