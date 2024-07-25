package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/dto"
	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/infra/webclient"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func NewWeatherByAddress(ctx context.Context, tracer trace.Tracer, a dto.AddressDto, client *http.Client) (*dto.WeatherDto, error) {

	ctx, span := tracer.Start(ctx, "NewWeatherByServiceB")
	defer span.End()

	var urlQuery = map[string]string{}
	urlQuery["key"] = os.Getenv("WEATHER_API_KEY")
	urlQuery["q"] = a.Localidade
	urlQuery["aqi"] = "no"

	wcReq, err := webclient.NewWebclient(ctx, client, http.MethodGet, "https://api.weatherapi.com/v1/current.json", urlQuery)
	if err != nil {
		slog.Error("[weatherapi webserver client]", "error", err.Error())
		return nil, err
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(wcReq.Request().Header))

	slog.Debug("[wcReq.Request().Header]", "Header", wcReq.Request().Header)

	var w dto.WeatherDto

	err = wcReq.Do(func(p []byte) error {
		err = json.Unmarshal(p, &w)
		if err != nil {
			slog.Error("[weather body unmarshal]", "error", err.Error())
		}
		return err
	})
	if err != nil {
		slog.Error("[weather do]", "error", err.Error())
		return nil, err

	}

	slog.Debug("[struct]", "WeatherResponseDto", w)

	return &w, nil
}

func NewWeatherByServiceB(ctx context.Context, tracer trace.Tracer, cli *http.Client, z dto.ZipcodeDto) (*dto.LocalWeatherDto, error) {

	ctx, span := tracer.Start(ctx, "NewWeatherByServiceB")
	defer span.End()

	wcReq, err := webclient.NewWebclient(ctx, cli, http.MethodGet, "http://"+os.Getenv("SERVICE_B_HOST")+":"+os.Getenv("SERVICE_B_PORT")+"/zipcode/"+z.Zipcode, nil)
	if err != nil {
		slog.Error("[service b webclient]", "error", err.Error())
		return nil, err
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(wcReq.Request().Header))

	var l dto.LocalWeatherDto

	err = wcReq.Do(func(p []byte) error {
		err = json.Unmarshal(p, &l)
		if err != nil {
			slog.Error("[service b body unmarshal]", "error", err.Error())
		}
		return err
	})
	if err != nil {
		slog.Error("[service b webclient do]", "error", err.Error())
		return nil, err
	}

	return &l, nil
}
