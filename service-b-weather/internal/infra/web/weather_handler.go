package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/felipeksw/goexpert-fullcycle-observability/service-b-weather/internal/dto"
	"github.com/felipeksw/goexpert-fullcycle-observability/service-b-weather/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type WeatherHandler struct {
	weatherFinder usecase.Finder
	localeFinder  usecase.Finder
}

func NewWeatherHandler(weatherFinder usecase.Finder, localeFinder usecase.Finder) *WeatherHandler {
	return &WeatherHandler{
		weatherFinder: weatherFinder,
		localeFinder:  localeFinder,
	}
}

func (h *WeatherHandler) Handle(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	tracer := otel.Tracer("service-b-weather")
	_, span := tracer.Start(ctx, "weather-handler")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")

	zipcode := r.PathValue("cep")

	var re = regexp.MustCompile(`^[0-9]{8}$`)
	if !re.MatchString(zipcode) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(&dto.ErrorOutput{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "invalid zipcode",
		})
		return
	}

	localeOutputRaw, err := h.localeFinder.Execute(ctx, zipcode)
	if err != nil {
		slog.Error("[weather handler localeFinder]", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&dto.ErrorOutput{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
		return
	}

	localeOutput := localeOutputRaw.(*dto.LocaleOutput)
	if localeOutput.Localidade == "" {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(&dto.ErrorOutput{
			StatusCode: http.StatusNotFound,
			Message:    "can not find zipcode",
		})
		return
	}

	weatherOutputRaw, err := h.weatherFinder.Execute(ctx, localeOutput.Localidade)
	if err != nil {
		slog.Error("[weather handler weatherFinder]", "error", err.Error())
		if err.Error() == "API key is invalid" || err.Error() == "API key is not provided" {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(&dto.ErrorOutput{
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
			})
			return
		}

		if err.Error() == "can not find zipcode" {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(&dto.ErrorOutput{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&dto.ErrorOutput{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
		return
	}

	weatherOutput := weatherOutputRaw.(*dto.WeatherOutput)

	w.WriteHeader(http.StatusOK)
	result := dto.ResultOutput{
		City:  localeOutput.Localidade,
		TempC: weatherOutput.Current.TempC,
		TempF: weatherOutput.Current.TempF,
		TempK: weatherOutput.Current.TempC + 273.15,
	}

	_ = json.NewEncoder(w).Encode(result)
}
