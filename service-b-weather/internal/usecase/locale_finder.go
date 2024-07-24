package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/felipeksw/goexpert-fullcycle-observability/service-b-weather/internal/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type LocaleFinder struct {
	httpClient *http.Client
}

func NewLocaleFinder(httpClient *http.Client) *LocaleFinder {
	return &LocaleFinder{httpClient: httpClient}
}

func (l *LocaleFinder) Execute(ctx context.Context, zipcode string) (interface{}, error) {
	tracer := otel.Tracer("service-b-weather")
	_, span := tracer.Start(ctx, "zipcode-search")
	defer span.End()

	input := &dto.LocaleInput{
		Cep: zipcode,
	}
	inputJson, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+os.Getenv("SERVICE_A_HOST")+":"+os.Getenv("SERVICE_A_PORT")+"/zipcode/", bytes.NewBuffer(inputJson))
	if err != nil {
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	req.Header.Set("Content-type", "application/json")
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
