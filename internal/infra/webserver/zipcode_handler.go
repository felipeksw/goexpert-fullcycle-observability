package webserver

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/dto"
	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/entity"
	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func GetZipcodeHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))

	trc := otel.Tracer("weatherByZipcode-tracer")

	slog.Debug("[struct]", "r.Body", r.Body)

	var z dto.ZipcodeBodyDto

	err := json.NewDecoder(r.Body).Decode(&z)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&dto.ErroDto{Msg: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")

	httpClient := http.DefaultClient

	slog.Debug("[struct]", "z.Cep", z.Cep)

	zipcodeDto, err := entity.NewZipcode(z.Cep)
	if err != nil {
		stsCod := http.StatusInternalServerError
		stsMsg := err.Error()

		if strings.Contains(strings.ToLower(err.Error()), "invalid zipcode") {
			stsCod = http.StatusNotFound
			stsMsg = "invalid zipcode"
		}

		w.WriteHeader(stsCod)
		json.NewEncoder(w).Encode(&dto.ErroDto{Msg: stsMsg})
		return
	}

	slog.Debug("[struct]", "zipcodeDto", zipcodeDto)

	localeWeatherDto, err := usecase.NewWeatherByServiceB(ctx, trc, httpClient, *zipcodeDto)
	if err != nil {
		stsCod := http.StatusInternalServerError
		stsMsg := err.Error()

		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			stsCod = http.StatusNotFound
			stsMsg = "can not find zipcode"
		}

		w.WriteHeader(stsCod)
		json.NewEncoder(w).Encode(&dto.ErroDto{Msg: stsMsg})
		return
	}

	slog.Debug("[struct]", "localeWeatherDto", localeWeatherDto)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(localeWeatherDto)
}
