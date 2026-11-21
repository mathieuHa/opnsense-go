package unbound

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestDnsbl(t *testing.T) {
	controller := Controller{Api: api.NewClient(api.Options{
		Uri:           os.Getenv("OPNSENSE_URI"),
		APIKey:        os.Getenv("OPNSENSE_API_KEY"),
		APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})}
	ctx := context.Background()

	blocklist := &Dnsbl{
		Enabled:     "1",
		Type:        api.SelectedMapList{},
		Lists:       api.SelectedMapList{},
		Allowlists:  api.SelectedMapList{"allowed.example.com"},
		Blocklists:  api.SelectedMapList{"ads.example.com", "tracker.example.com"},
		Wildcards:   api.SelectedMapList{},
		SourceNets:  api.SelectedMapList{"192.0.2.0/24"},
		NXDomain:    "1",
		CacheTTL:    "72000",
		Description: "test-dnsbl",
	}

	id, err := controller.AddDnsbl(ctx, blocklist)
	if err != nil {
		t.Fatalf("AddDnsbl failed: %v", err)
	}
	t.Logf("AddDnsbl: uuid=%s", id)

	got, err := controller.GetDnsbl(ctx, id)
	if err != nil {
		t.Fatalf("GetDnsbl failed: %v", err)
	}
	t.Logf("GetDnsbl: %+v", got)
	if len(got.Blocklists) != 2 || len(got.SourceNets) != 1 || got.SourceNets[0] != "192.0.2.0/24" || got.NXDomain != "1" {
		t.Fatalf("GetDnsbl mismatch: %+v", got)
	}

	blocklist.Blocklists = api.SelectedMapList{"ads.example.com"}
	blocklist.Wildcards = api.SelectedMapList{"example.net"}
	if err := controller.UpdateDnsbl(ctx, id, blocklist); err != nil {
		t.Fatalf("UpdateDnsbl failed: %v", err)
	}
	if got, err = controller.GetDnsbl(ctx, id); err != nil || len(got.Blocklists) != 1 || len(got.Wildcards) != 1 {
		t.Fatalf("UpdateDnsbl not applied: %+v (err %v)", got, err)
	}

	if err := controller.DeleteDnsbl(ctx, id); err != nil {
		t.Fatalf("DeleteDnsbl failed: %v", err)
	}
}
