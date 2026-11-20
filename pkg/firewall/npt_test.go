package firewall

import (
	"context"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestNPT(t *testing.T) {
	controller := newController(t)
	ctx := context.Background()

	rule := &NPT{
		Enabled:        "1",
		Log:            "0",
		Sequence:       "100",
		Interface:      api.SelectedMap("wan"),
		SourceNet:      "fd00:1::/64",
		DestinationNet: "2001:db8:1::/64",
		TrackInterface: api.SelectedMap(""),
		Categories:     api.SelectedMapList{},
		Description:    "test-npt",
	}

	id, err := controller.AddNPT(ctx, rule)
	if err != nil {
		t.Fatalf("AddNPT failed: %v", err)
	}
	t.Logf("AddNPT: uuid=%s", id)

	got, err := controller.GetNPT(ctx, id)
	if err != nil {
		t.Fatalf("GetNPT failed: %v", err)
	}
	t.Logf("GetNPT: %+v", got)
	if got.SourceNet != "fd00:1::/64" || got.DestinationNet != "2001:db8:1::/64" || got.Interface.String() != "wan" {
		t.Fatalf("GetNPT mismatch: %+v", got)
	}

	rule.DestinationNet = "2001:db8:2::/64"
	rule.Log = "1"
	if err := controller.UpdateNPT(ctx, id, rule); err != nil {
		t.Fatalf("UpdateNPT failed: %v", err)
	}
	if got, err = controller.GetNPT(ctx, id); err != nil || got.DestinationNet != "2001:db8:2::/64" || got.Log != "1" {
		t.Fatalf("UpdateNPT not applied: %+v (err %v)", got, err)
	}

	if err := controller.DeleteNPT(ctx, id); err != nil {
		t.Fatalf("DeleteNPT failed: %v", err)
	}
}
