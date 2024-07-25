package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/dto"
	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/entity"
	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func GetWeatherByZipcodeHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))

	tracer := otel.Tracer("weatherByZipcode-tracer")

	w.Header().Add("Content-Type", "application/json")

	httpClient := http.DefaultClient

	zipcodeDto, err := entity.NewZipcode(r.PathValue("zipcode"))
	if err != nil {

		code := http.StatusInternalServerError
		if err.Error() == "zip code must be 8 numeric digits" {
			code = http.StatusNotFound
		}

		w.WriteHeader(code)
		json.NewEncoder(w).Encode(&dto.ErroDto{Msg: err.Error()})
		w.Write([]byte(err.Error()))
		return
	}

	addressDto, err := usecase.NewAddressByZipcode(ctx, tracer, *zipcodeDto, httpClient)
	if err != nil {

		code := http.StatusInternalServerError
		if err.Error() == "zip code not found" {
			code = http.StatusNotFound
		}

		w.WriteHeader(code)
		json.NewEncoder(w).Encode(&dto.ErroDto{Msg: err.Error()})
		w.Write([]byte(err.Error()))
		return
	}

	weatherDto, err := usecase.NewWeatherByAddress(ctx, tracer, *addressDto, httpClient)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&dto.ErroDto{Msg: err.Error()})
		return
	}

	localeWeatherDto, err := entity.NewLocaleWeather(addressDto.Localidade, weatherDto.Current.TempC)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&dto.ErroDto{Msg: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(localeWeatherDto)
}
