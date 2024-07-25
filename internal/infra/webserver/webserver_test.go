package webserver_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/felipeksw/goexpert-fullcycle-cloud-run/internal/infra/webserver"
	"github.com/stretchr/testify/assert"
)

func TestWebServer(t *testing.T) {

	ws := webserver.NewWebServer("9080")
	ws.AddHandler("GET /ping", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	go func() {
		err := ws.Start()
		assert.Nil(t, err)
	}()

	//---
	time.Sleep(30 * time.Millisecond)
	req, err := http.NewRequest(http.MethodGet, "http://localhost:9080/ping", nil)
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}
