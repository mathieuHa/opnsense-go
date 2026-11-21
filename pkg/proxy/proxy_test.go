package proxy

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func testController() Controller {
	return Controller{
		Api: api.NewClient(api.Options{
			Uri:           os.Getenv("OPNSENSE_URI"),
			APIKey:        os.Getenv("OPNSENSE_API_KEY"),
			APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
			AllowInsecure: true,
		}),
	}
}

func TestPacProxy(t *testing.T) {
	c := testController()
	ctx := context.Background()
	if _, err := c.SettingsGet(ctx); err != nil {
		t.Skipf("proxy plugin not installed: %s", err)
	}

	id, err := c.AddPacProxy(ctx, &PacProxy{Name: "sdk-test", ProxyType: "PROXY", URL: "proxy.example.com:3128", Description: "sdk test"})
	if err != nil {
		t.Fatal(err)
	}
	defer c.DeletePacProxy(ctx, id)

	got, err := c.GetPacProxy(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "sdk-test" || got.ProxyType.String() != "PROXY" || got.URL != "proxy.example.com:3128" {
		t.Fatalf("unexpected PAC proxy: %+v", got)
	}

	got.URL = "proxy2.example.com:3128"
	if err := c.UpdatePacProxy(ctx, id, got); err != nil {
		t.Fatal(err)
	}
	if got, _ = c.GetPacProxy(ctx, id); got.URL != "proxy2.example.com:3128" {
		t.Fatalf("update not applied: %+v", got)
	}
}

func TestRemoteBlacklist(t *testing.T) {
	c := testController()
	ctx := context.Background()
	if _, err := c.SettingsGet(ctx); err != nil {
		t.Skipf("proxy plugin not installed: %s", err)
	}

	id, err := c.AddRemoteBlacklist(ctx, &RemoteBlacklist{Enabled: "0", Filename: "sdktest", URL: "http://192.0.2.10/list.txt", Filter: api.SelectedMapList{}, SSLNoVerify: "0", Description: "sdk test"})
	if err != nil {
		t.Fatal(err)
	}
	defer c.DeleteRemoteBlacklist(ctx, id)

	got, err := c.GetRemoteBlacklist(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if got.URL != "http://192.0.2.10/list.txt" || got.Filename != "sdktest" {
		t.Fatalf("unexpected blacklist: %+v", got)
	}
}

func TestSettingsRoundTrip(t *testing.T) {
	c := testController()
	ctx := context.Background()
	s, err := c.SettingsGet(ctx)
	if err != nil {
		t.Skipf("proxy plugin not installed: %s", err)
	}
	orig := s.Proxy.General.VisibleHostname
	s.Proxy.General.VisibleHostname = "sdk-test.example"
	if _, err := c.SettingsSet(ctx, &s.Proxy); err != nil {
		t.Fatal(err)
	}
	after, err := c.SettingsGet(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if after.Proxy.General.VisibleHostname != "sdk-test.example" || after.Proxy.Forward.Port != s.Proxy.Forward.Port {
		t.Fatalf("round trip changed settings: %+v", after.Proxy.General)
	}
	after.Proxy.General.VisibleHostname = orig
	if _, err := c.SettingsSet(ctx, &after.Proxy); err != nil {
		t.Fatal(err)
	}
}
