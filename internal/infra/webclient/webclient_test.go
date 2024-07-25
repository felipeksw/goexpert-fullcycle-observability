package webclient_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/infra/mockup"
	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/infra/webclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewWebclient(t *testing.T) {

	type ReqresResponseData struct {
		Id         int    `json:"id"`
		Email      string `json:"email"`
		First_name string `json:"first_name"`
		Last_name  string `json:"last_name"`
		Avatar     string `json:"avatar"`
	}
	type ReqresResponseSupport struct {
		Url  string `json:"url"`
		Text string `json:"text"`
	}
	type ReqresResponse struct {
		Data    ReqresResponseData    `json:"data"`
		Support ReqresResponseSupport `json:"support"`
	}

	mocReqresResponseSuccessBody := `{"data":{"id":3,"email":"emma.wong@reqres.in","first_name":"Emma","last_name":"Wong","avatar":"https://reqres.in/img/faces/3-image.jpg"},"support":{"url":"https://reqres.in/#support-heading","text":"To keep ReqRes free, contributions towards server costs are appreciated!"}}`

	mockRoundTripper := new(mockup.MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(mocReqresResponseSuccessBody))),
	}, nil)

	var urlQuery = map[string]string{}
	urlQuery["key01"] = "value01"
	urlQuery["key02"] = "value02"
	urlQuery["output"] = "json"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	wc, err := webclient.NewWebclient(ctx, mockClient, http.MethodGet, "https://dummy.restapiexample.com/api/v1/employee/1", urlQuery)
	assert.Nil(t, err)

	var w ReqresResponse

	err = wc.Do(func(p []byte) error {
		err = json.Unmarshal(p, &w)
		assert.Nil(t, err)
		return err
	})
	assert.Nil(t, err)
	assert.Equal(t, 3, w.Data.Id)
}

func TestNewWebclientJsonError(t *testing.T) {

	type ReqresResponseData struct {
		Id         string `json:"id"`
		Email      string `json:"email"`
		First_name string `json:"first_name"`
		Last_name  string `json:"last_name"`
		Avatar     string `json:"avatar"`
	}
	type ReqresResponseSupport struct {
		Url  string `json:"url"`
		Text string `json:"text"`
	}
	type ReqresResponse struct {
		Data    ReqresResponseData    `json:"data"`
		Support ReqresResponseSupport `json:"support"`
	}

	mocReqresResponseSuccessBody := `{"data":{"id":3,"email":"emma.wong@reqres.in","first_name":"Emma","last_name":"Wong","avatar":"https://reqres.in/img/faces/3-image.jpg"},"support":{"url":"https://reqres.in/#support-heading","text":"To keep ReqRes free, contributions towards server costs are appreciated!"}}`

	mockRoundTripper := new(mockup.MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(mocReqresResponseSuccessBody))),
	}, nil)

	var urlQuery = map[string]string{}
	urlQuery["key01"] = "value01"
	urlQuery["key02"] = "value02"
	urlQuery["output"] = "json"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	wc, err := webclient.NewWebclient(ctx, mockClient, http.MethodGet, "https://dummy.restapiexample.com/api/v1/employee/1", urlQuery)
	assert.Nil(t, err)

	var w ReqresResponse

	err = wc.Do(func(p []byte) error {
		return json.Unmarshal(p, &w)
	})
	assert.Contains(t, err.Error(), "cannot unmarshal")
}

func TestNewWebclientHttpStatusError(t *testing.T) {

	type ReqresResponseData struct {
		Id         int    `json:"id"`
		Email      string `json:"email"`
		First_name string `json:"first_name"`
		Last_name  string `json:"last_name"`
		Avatar     string `json:"avatar"`
	}
	type ReqresResponseSupport struct {
		Url  string `json:"url"`
		Text string `json:"text"`
	}
	type ReqresResponse struct {
		Data    ReqresResponseData    `json:"data"`
		Support ReqresResponseSupport `json:"support"`
	}

	mocReqresResponseNotFountBody := `{"message":"Error Occured! Page Not found, contact rstapi2example@gmail.com"}`

	mockRoundTripper := new(mockup.MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(bytes.NewReader([]byte(mocReqresResponseNotFountBody))),
	}, nil)

	var urlQuery = map[string]string{}
	urlQuery["key01"] = "value01"
	urlQuery["key02"] = "value02"
	urlQuery["output"] = "json"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	wc, err := webclient.NewWebclient(ctx, mockClient, http.MethodGet, "https://dummy.restapiexample.com/api/v1/employee/1", urlQuery)
	assert.Nil(t, err)

	var w ReqresResponse

	err = wc.Do(func(p []byte) error {
		err = json.Unmarshal(p, &w)
		assert.Nil(t, err)
		return err
	})
	assert.Contains(t, err.Error(), "Not Found")
	assert.Equal(t, http.MethodGet, wc.Request().Method)
}
