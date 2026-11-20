package core

import (
	"context"
	"testing"
)

func TestTunable(t *testing.T) {
	controller := newController()
	ctx := context.Background()

	tunable := &Tunable{
		Tunable:     "net.inet.tcp.delayed_ack",
		Value:       "0",
		Description: "test-tunable",
	}

	id, err := controller.AddTunable(ctx, tunable)
	if err != nil {
		t.Fatalf("AddTunable failed: %v", err)
	}
	t.Logf("AddTunable: uuid=%s", id)

	got, err := controller.GetTunable(ctx, id)
	if err != nil {
		t.Fatalf("GetTunable failed: %v", err)
	}
	t.Logf("GetTunable: %+v", got)
	if got.Tunable != "net.inet.tcp.delayed_ack" || got.Value != "0" {
		t.Fatalf("GetTunable mismatch: %+v", got)
	}

	tunable.Value = "1"
	if err := controller.UpdateTunable(ctx, id, tunable); err != nil {
		t.Fatalf("UpdateTunable failed: %v", err)
	}
	if got, err = controller.GetTunable(ctx, id); err != nil || got.Value != "1" {
		t.Fatalf("UpdateTunable not applied: %+v (err %v)", got, err)
	}

	if err := controller.DeleteTunable(ctx, id); err != nil {
		t.Fatalf("DeleteTunable failed: %v", err)
	}
}
