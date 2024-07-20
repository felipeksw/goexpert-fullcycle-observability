package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/felipeksw/goexpert-fullcycle-observability/service-b-weather/internal/dto"
	"go.opentelemetry.io/otel"
)

type WeatherFinder struct {
	httpClient *http.Client
}

func NewWeatherFinder(httpClient *http.Client) *WeatherFinder {
	return &WeatherFinder{httpClient: httpClient}
}

func (w *WeatherFinder) Execute(ctx context.Context, zipcode string) (interface{}, error) {
	tracer := otel.Tracer("service-b-weather")
	_, span := tracer.Start(ctx, "weather-search")
	defer span.End()

	req, err := http.NewRequest(http.MethodGet, "https://api.weatherapi.com/v1/current.json?key="+os.Getenv("KEY_WEATHER_API")+"&q="+url.QueryEscape(zipcode), nil)
	if err != nil {
		return nil, err
	}
	slog.Debug("[calling api]", "url", req.URL)

	//req.Header.Set("Content-type", "application/json")

	res, err := w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	_ = res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("API key is invalid")
	}

	if res.StatusCode == http.StatusBadRequest {
		return nil, errors.New("can not find zipcode")
	}

	var output dto.WeatherOutput
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
