package api

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientDefaultBackoffIsSeconds(t *testing.T) {
	c := NewClient(Options{})
	if c.client.RetryWaitMin != time.Second || c.client.RetryWaitMax != 30*time.Second {
		t.Fatalf("backoff = %s..%s, want 1s..30s", c.client.RetryWaitMin, c.client.RetryWaitMax)
	}
	if c.client.HTTPClient.Timeout != clientRequestTimeout {
		t.Fatalf("timeout = %s, want %s", c.client.HTTPClient.Timeout, clientRequestTimeout)
	}
}

func TestClientRequestTimeout(t *testing.T) {
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
	}))
	defer srv.Close()
	defer close(release)

	c := NewClient(Options{Uri: srv.URL, RequestTimeout: 1, MaxRetries: 1, MinBackoff: 1, MaxBackoff: 1, Logger: log.New(&bytes.Buffer{}, "", 0)})
	start := time.Now()
	err := c.doRequest(context.Background(), "GET", "/core/firmware/status", nil, &struct{}{})
	if err == nil {
		t.Fatal("expected a timeout error")
	}
	if d := time.Since(start); d > 10*time.Second {
		t.Fatalf("request took %s, the timeout was not applied", d)
	}
}

func TestClientDoesNotLogBodies(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"privkey":"response-secret"}`))
	}))
	defer srv.Close()

	var buf bytes.Buffer
	c := NewClient(Options{Uri: srv.URL, APIKey: "key", APISecret: "secret", Logger: log.New(&buf, "", 0)})
	resp := map[string]string{}
	if err := c.doRequest(context.Background(), "POST", "/wireguard/server/add_server", map[string]string{"privkey": "request-secret"}, &resp); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"request-secret", "response-secret", "secret"} {
		if strings.Contains(buf.String(), s) {
			t.Fatalf("log output contains %q: %s", s, buf.String())
		}
	}
	if !strings.Contains(buf.String(), "POST /api/wireguard/server/add_server -> 200") {
		t.Fatalf("unexpected log output: %s", buf.String())
	}
}
