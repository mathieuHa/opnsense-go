package routing

import (
	"context"
	"strings"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestGatewayGroup(t *testing.T) {
	controller := newController()
	ctx := context.Background()

	if _, err := controller.GetGatewayGroup(ctx, "00000000-0000-0000-0000-000000000000"); err != nil && strings.Contains(err.Error(), "status code 404") {
		t.Skip("gateway group API not available; requires OPNsense 26.7 or later")
	}

	primaryID, err := controller.AddGateway(ctx, newTestGateway("test_gw_primary", "192.0.2.1"))
	if err != nil {
		t.Fatalf("AddGateway primary failed: %v", err)
	}
	defer controller.DeleteGateway(ctx, primaryID)

	backupID, err := controller.AddGateway(ctx, newTestGateway("test_gw_backup", "192.0.2.2"))
	if err != nil {
		t.Fatalf("AddGateway backup failed: %v", err)
	}
	defer controller.DeleteGateway(ctx, backupID)

	group := &GatewayGroup{
		Name:        "test_gw_group",
		Tier1:       api.SelectedMapList{"test_gw_primary"},
		Tier2:       api.SelectedMapList{"test_gw_backup"},
		Trigger:     api.SelectedMap("down"),
		Description: "test-gateway-group",
	}

	id, err := controller.AddGatewayGroup(ctx, group)
	if err != nil {
		t.Fatalf("AddGatewayGroup failed: %v", err)
	}
	t.Logf("AddGatewayGroup: uuid=%s", id)

	got, err := controller.GetGatewayGroup(ctx, id)
	if err != nil {
		t.Fatalf("GetGatewayGroup failed: %v", err)
	}
	t.Logf("GetGatewayGroup: %+v", got)
	if got.Name != "test_gw_group" || len(got.Tier1) != 1 || got.Tier1[0] != "test_gw_primary" ||
		len(got.Tier2) != 1 || got.Tier2[0] != "test_gw_backup" || got.Trigger.String() != "down" {
		t.Fatalf("GetGatewayGroup mismatch: %+v", got)
	}

	group.Tier1 = api.SelectedMapList{"test_gw_primary", "test_gw_backup"}
	group.Tier2 = api.SelectedMapList{}
	group.Trigger = api.SelectedMap("downloss")
	err = controller.UpdateGatewayGroup(ctx, id, group)
	if err != nil {
		t.Fatalf("UpdateGatewayGroup failed: %v", err)
	}

	got, err = controller.GetGatewayGroup(ctx, id)
	if err != nil {
		t.Fatalf("GetGatewayGroup after update failed: %v", err)
	}
	if len(got.Tier1) != 2 || len(got.Tier2) != 0 || got.Trigger.String() != "downloss" {
		t.Fatalf("UpdateGatewayGroup not applied: %+v", got)
	}

	err = controller.DeleteGatewayGroup(ctx, id)
	if err != nil {
		t.Fatalf("DeleteGatewayGroup failed: %v", err)
	}
	t.Log("DeleteGatewayGroup: deleted")
}
