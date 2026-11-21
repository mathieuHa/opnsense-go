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

func TestNestedMonad(t *testing.T) {
	var posted string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			b := new(bytes.Buffer)
			b.ReadFrom(r.Body)
			posted = b.String()
			w.Write([]byte(`{"result":"saved","uuid":"u1"}`))
		default:
			w.Write([]byte(`{"pac":{"proxy":{"name":"p1"}}}`))
		}
	}))
	defer srv.Close()

	type item struct {
		Name string `json:"name"`
	}
	c := NewClient(Options{Uri: srv.URL, Logger: log.New(&bytes.Buffer{}, "", 0)})
	opts := ReqOpts{AddEndpoint: "/proxy/settings/add_pac_proxy", GetEndpoint: "/proxy/settings/get_pac_proxy", Monad: "pac.proxy"}
	if _, err := Add(c, context.Background(), opts, &item{Name: "p1"}); err != nil {
		t.Fatal(err)
	}
	if posted != `{"proxy":{"name":"p1"}}` {
		t.Fatalf("posted %s", posted)
	}
	got, err := Get(c, context.Background(), opts, &item{}, "u1")
	if err != nil || got.Name != "p1" {
		t.Fatalf("got %+v, %v", got, err)
	}
}
