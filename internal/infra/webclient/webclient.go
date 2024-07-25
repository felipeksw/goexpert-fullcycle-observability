package webclient

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

type webClient struct {
	request *http.Request
	client  *http.Client
}

func (w *webClient) Request() *http.Request {
	return w.request
}

func NewWebclient(ctx context.Context, client *http.Client, method string, url string, query map[string]string) (*webClient, error) {

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		slog.Error("[http.NewRequest failed]", "error", err.Error())
		return nil, err
	}

	if ctx != nil {
		req = req.WithContext(ctx)
		slog.Debug("[Context Added]")
	}

	if query != nil {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	return &webClient{
		request: req,
		client:  client,
	}, nil
}

func (w *webClient) Do(ret func([]byte) error) error {

	slog.Debug("[http client Do host]", "host", w.request.URL.Host)
	slog.Debug("[http client Do full url]", "url", w.request.URL)

	resp, err := w.client.Do(w.request)
	if err != nil {
		slog.Debug("[http Client Do failed]", "error", err.Error())
		return errors.New("error to execute http request: " + w.request.URL.Host)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	defer func() {
		body = nil
	}()
	if err != nil {
		slog.Error("[io.ReadAll failed]", "error", err.Error())
		return err
	}

	slog.Debug("[http client Do status]", "status", resp.Status)
	slog.Debug("[http client Do statuscode]", "code", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return errors.New(w.request.URL.Host + ": " + http.StatusText(resp.StatusCode))
	}

	slog.Debug("[http client Do body]", "body", body)

	return ret(body)
}
