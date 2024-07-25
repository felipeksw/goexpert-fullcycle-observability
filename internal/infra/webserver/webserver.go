package webserver

import (
	"log/slog"
	"net/http"
)

type WebServer struct {
	WebServerPort string
	Mux           *http.ServeMux
}

func NewWebServer(serverPort string) *WebServer {
	slog.Info("[webserver created]")

	return &WebServer{
		WebServerPort: serverPort,
		Mux:           http.NewServeMux(),
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Mux.Handle(path, handler)
	slog.Info("[route added]", "path", path)
}

func (s *WebServer) Start() error {
	slog.Info("[server listening]", "port", s.WebServerPort)

	err := http.ListenAndServe(":"+s.WebServerPort, s.Mux)
	if err != nil {
		return err
	}

	return nil
}
