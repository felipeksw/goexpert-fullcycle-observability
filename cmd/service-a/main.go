package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/infra/webserver"
	otelpkg "github.com/felipeksw/goexpert-fullcycle-cloud-run/pkg/otel"
)

func main() {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	chSo := make(chan os.Signal, 1)
	signal.Notify(chSo, os.Interrupt, syscall.SIGINT)

	ctx, shutdownSo := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer shutdownSo()

	ShutdownProvider, err := otelpkg.InitProvider(ctx, "service-a", "otel-collector:4317")
	if err != nil {
		slog.Error("[InitProvider]", "error", err.Error())
		os.Exit(5)
	}
	defer func() {
		if err := ShutdownProvider(ctx); err != nil {
			slog.Error("[ShutdownProvider]", "error", err.Error())
			os.Exit(5)
		}
	}()

	ws := webserver.NewWebServer(os.Getenv("SERVICE_A_PORT"))
	ws.AddHandler("POST /zipcode/", webserver.GetZipcodeHandler)
	errWs := ws.Start()
	if errWs != nil {
		slog.Error("could not start the webserver:" + errWs.Error())
	}

	select {
	case <-chSo:
		slog.Info("Shutting down gracefully wbc2, CTRL+C pressed...")
	case <-ctx.Done():
		slog.Info("Shutting down gracefully wbc2, interrupet system...")
	}
}
