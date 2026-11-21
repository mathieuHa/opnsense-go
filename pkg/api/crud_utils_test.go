package api

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteMissingItemSucceeds(t *testing.T) {
	reconfigured := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/routes/routes/reconfigure" {
			reconfigured = true
			w.Write([]byte(`{"status":"ok"}`))
			return
		}
		w.Write([]byte(`{"result":"not found"}`))
	}))
	defer srv.Close()

	c := NewClient(Options{Uri: srv.URL, Logger: log.New(&bytes.Buffer{}, "", 0)})
	opts := ReqOpts{DeleteEndpoint: "/routes/routes/delroute", ReconfigureEndpoint: "/routes/routes/reconfigure"}
	if err := Delete(c, context.Background(), opts, "00000000-0000-0000-0000-000000000000"); err != nil {
		t.Fatalf("deleting a missing item failed: %s", err)
	}
	if reconfigured {
		t.Fatal("reconfigure was called although nothing was deleted")
	}
}
