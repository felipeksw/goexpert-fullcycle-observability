package usecase

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/felipeksw/goexpert-fullcycle-observability/service-a-zipcode/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	mockCep                = "01001000"
	mockInvalidCep         = "InvalidLocale"
	mockLocaleResponseBody = `{"cep": "01001000", "logradouro": "Avenida Paulista", "bairro": "Jardins", "localidade": "São Paulo", "uf": "SP"}`
)

type MockRoundTripper struct {
	mock.Mock
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) != nil {
		return args.Get(0).(*http.Response), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockContext struct{}

func (mc *MockContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (mc *MockContext) Done() <-chan struct{} {
	return nil
}

func (mc *MockContext) Err() error {
	return nil
}

func (mc *MockContext) Value(key interface{}) interface{} {
	return nil
}

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
