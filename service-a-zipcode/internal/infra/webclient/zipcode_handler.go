package webclient

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/felipeksw/goexpert-fullcycle-observability/service-a-zipcode/internal/dto"
	"github.com/felipeksw/goexpert-fullcycle-observability/service-a-zipcode/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type LocaleHandler struct {
	localeFinder usecase.Finder
}

func NewLocaleHandler(localeFinder usecase.Finder) *LocaleHandler {
	return &LocaleHandler{
		localeFinder: localeFinder,
	}
}

func (h *LocaleHandler) Handle(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	tracer := otel.Tracer("service-a-zipcode")
	_, span := tracer.Start(ctx, "zipcode-handler")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")

	var input dto.LocaleInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(msg)
		return
	}

	var re = regexp.MustCompile(`^[0-9]{8}$`)
	if !re.MatchString(input.Cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(&dto.ErrorOutput{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "invalid zipcode",
		})
		return
	}

	output, err := h.localeFinder.Execute(ctx, input.Cep)
	if err != nil {
		slog.Error("[hanlder execute]", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&dto.ErrorOutput{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
		return
	}

	if output.Localidade == "" {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(&dto.ErrorOutput{
			StatusCode: http.StatusNotFound,
			Message:    "can not find zipcode",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("[struct]", "LocaleOutput", output)
	_ = json.NewEncoder(w).Encode(output)
}
