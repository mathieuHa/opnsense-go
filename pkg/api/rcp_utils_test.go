package api

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCallReturnsValidationFailures(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"result":"failed","validations":{"general.threads":"Invalid integer value."}}`))
	}))
	defer srv.Close()

	c := NewClient(Options{Uri: srv.URL, Logger: log.New(&bytes.Buffer{}, "", 0)})
	opts := RPCOpts{BaseEndpoint: "/ipsec/settings/set", Method: "POST", BodyParameters: map[string]any{"ipsec": map[string]string{"threads": "x"}}}
	_, err := Call(c, context.Background(), opts, &struct {
		Result string `json:"result"`
	}{})
	if err == nil || !strings.Contains(err.Error(), "Invalid integer value.") {
		t.Fatalf("expected the validation message in the error, got %v", err)
	}
}

func TestCallAcceptsArrays(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"result":"failed"}]`))
	}))
	defer srv.Close()

	c := NewClient(Options{Uri: srv.URL, Logger: log.New(&bytes.Buffer{}, "", 0)})
	var out []map[string]string
	if _, err := Call(c, context.Background(), RPCOpts{BaseEndpoint: "/x", Method: "GET"}, &out); err != nil {
		t.Fatalf("array responses must not be treated as failures: %s", err)
	}
}
