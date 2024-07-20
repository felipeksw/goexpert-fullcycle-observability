package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/felipeksw/goexpert-fullcycle-observability/pkg/otel"
	"github.com/felipeksw/goexpert-fullcycle-observability/service-b-weather/internal/infra/web"
	"github.com/felipeksw/goexpert-fullcycle-observability/service-b-weather/internal/usecase"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := otel.SetupOTelSDK("service-b-weather", ctx)
	if err != nil {
		slog.Error("[setup OTEeL fail]", "error", err.Error())
		os.Exit(11)
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	srv := &http.Server{
		Addr:         ":" + os.Getenv("SERVICE_B_PORT"),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(),
	}
	srvErr := make(chan error, 1)
	go func() {
		slog.Info("[server listening...]", "port", srv.Addr)
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		slog.Error("[server listening fail]", "error", err.Error())
		os.Exit(11)
	case <-ctx.Done():
		stop()
	}

	err = srv.Shutdown(context.Background())
	if err != nil {
		slog.Error("[shutdown contenxt]", "error", err.Error())
	}
}

func newHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	weatherFinder := usecase.NewWeatherFinder(http.DefaultClient)
	localeFinder := usecase.NewLocaleFinder(http.DefaultClient)
	mux.HandleFunc("GET /zipcode/{cep}", web.NewWeatherHandler(weatherFinder, localeFinder).Handle)

	return otelhttp.NewHandler(mux, "/")
}
