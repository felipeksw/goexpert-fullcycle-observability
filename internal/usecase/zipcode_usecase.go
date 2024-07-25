package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/dto"
	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/infra/webclient"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func NewAddressByZipcode(ctx context.Context, tracer trace.Tracer, z dto.ZipcodeDto, client *http.Client) (*dto.AddressDto, error) {

	ctx, span := tracer.Start(ctx, "NewAddressByZipcode")
	defer span.End()

	wcReq, err := webclient.NewWebclient(ctx, client, http.MethodGet, "https://viacep.com.br/ws/"+z.Zipcode+"/json/", nil)
	if err != nil {
		slog.Error("[viacep NewWebclient failed]", "error", err.Error())
		return nil, err
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(wcReq.Request().Header))

	var a dto.AddressDto

	err = wcReq.Do(func(p []byte) error {
		err = json.Unmarshal(p, &a)
		if err != nil {
			slog.Error("[zipcode body unmarshal]", "error", err.Error())
		}
		return err
	})
	if err != nil {
		slog.Error("[webclient do]", "error", err.Error())

	}
	slog.Debug("[zipcode body]", "body", a)

	if a.Error != "" {
		return nil, errors.New("zip code not found")
	}

	return &a, err
}
