package usecase

import (
	"github.com/stretchr/testify/mock"
	"net/http"
	"time"
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
