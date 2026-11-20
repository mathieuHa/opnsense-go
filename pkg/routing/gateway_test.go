package routing

import (
	"context"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func newTestGateway(name, address string) *Gateway {
	return &Gateway{
		Disabled:       "0",
		Name:           name,
		Description:    "test-gateway",
		Interface:      api.SelectedMap("wan"),
		IPProtocol:     api.SelectedMap("inet"),
		Gateway:        address,
		FarGateway:     "1",
		MonitorDisable: "1",
		Priority:       "255",
		Weight:         "1",
	}
}

func TestGateway(t *testing.T) {
	controller := newController()
	ctx := context.Background()

	gateway := newTestGateway("test_gw", "192.0.2.1")

	id, err := controller.AddGateway(ctx, gateway)
	if err != nil {
		t.Fatalf("AddGateway failed: %v", err)
	}
	t.Logf("AddGateway: uuid=%s", id)

	got, err := controller.GetGateway(ctx, id)
	if err != nil {
		t.Fatalf("GetGateway failed: %v", err)
	}
	t.Logf("GetGateway: %+v", got)
	if got.Name != "test_gw" || got.Gateway != "192.0.2.1" || got.Interface.String() != "wan" {
		t.Fatalf("GetGateway mismatch: %+v", got)
	}

	gateway.Description = "test-gateway-updated"
	gateway.Priority = "10"
	err = controller.UpdateGateway(ctx, id, gateway)
	if err != nil {
		t.Fatalf("UpdateGateway failed: %v", err)
	}

	got, err = controller.GetGateway(ctx, id)
	if err != nil {
		t.Fatalf("GetGateway after update failed: %v", err)
	}
	if got.Description != "test-gateway-updated" || got.Priority != "10" {
		t.Fatalf("UpdateGateway not applied: %+v", got)
	}

	err = controller.DeleteGateway(ctx, id)
	if err != nil {
		t.Fatalf("DeleteGateway failed: %v", err)
	}
	t.Log("DeleteGateway: deleted")
}
