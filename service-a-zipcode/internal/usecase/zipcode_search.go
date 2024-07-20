package usecase

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/felipeksw/goexpert-fullcycle-observability/service-a-zipcode/internal/dto"
	"go.opentelemetry.io/otel"
)

type LocaleFinder struct {
	httpClient *http.Client
}

func NewLocaleFinder(httpClient *http.Client) *LocaleFinder {
	return &LocaleFinder{httpClient: httpClient}
}

func (l *LocaleFinder) Execute(ctx context.Context, zipcode string) (*dto.LocaleOutput, error) {
	tracer := otel.Tracer("service-a-zipcode")
	_, span := tracer.Start(ctx, "zipcode-search")
	defer span.End()

	req, err := http.NewRequest(http.MethodGet, "https://viacep.com.br/ws/"+url.QueryEscape(zipcode)+"/json/", nil)
	if err != nil {
		return nil, err
	}
	slog.Debug("[calling api]", "url", req.URL)

	res, err := l.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	_ = res.Body.Close()

	var output dto.LocaleOutput
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
