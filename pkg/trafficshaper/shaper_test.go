package trafficshaper

import (
	"context"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestShaperChain(t *testing.T) {
	controller := newController()
	ctx := context.Background()

	pipe := &Pipe{
		Enabled:         "1",
		Bandwidth:       "10",
		BandwidthMetric: api.SelectedMap("Mbit"),
		Mask:            api.SelectedMap("none"),
		Scheduler:       api.SelectedMap("fq_codel"),
		CodelEnable:     "0",
		CodelECNEnable:  "0",
		PIEEnable:       "0",
		Description:     "test-pipe",
	}

	pipeID, err := controller.AddPipe(ctx, pipe)
	if err != nil {
		t.Fatalf("AddPipe failed: %v", err)
	}
	t.Logf("AddPipe: uuid=%s", pipeID)
	defer controller.DeletePipe(ctx, pipeID)

	gotPipe, err := controller.GetPipe(ctx, pipeID)
	if err != nil {
		t.Fatalf("GetPipe failed: %v", err)
	}
	t.Logf("GetPipe: %+v", gotPipe)
	if gotPipe.Number == "" || gotPipe.Bandwidth != "10" || gotPipe.BandwidthMetric.String() != "Mbit" || gotPipe.Scheduler.String() != "fq_codel" {
		t.Fatalf("GetPipe mismatch: %+v", gotPipe)
	}

	pipe.Number = gotPipe.Number
	pipe.Bandwidth = "20"
	if err := controller.UpdatePipe(ctx, pipeID, pipe); err != nil {
		t.Fatalf("UpdatePipe failed: %v", err)
	}
	if gotPipe, err = controller.GetPipe(ctx, pipeID); err != nil || gotPipe.Bandwidth != "20" {
		t.Fatalf("UpdatePipe not applied: %+v (err %v)", gotPipe, err)
	}

	queue := &Queue{
		Enabled:        "1",
		Pipe:           api.SelectedMap(pipeID),
		Weight:         "50",
		Mask:           api.SelectedMap("none"),
		CodelEnable:    "0",
		CodelECNEnable: "0",
		PIEEnable:      "0",
		Description:    "test-queue",
	}

	queueID, err := controller.AddQueue(ctx, queue)
	if err != nil {
		t.Fatalf("AddQueue failed: %v", err)
	}
	t.Logf("AddQueue: uuid=%s", queueID)
	defer controller.DeleteQueue(ctx, queueID)

	gotQueue, err := controller.GetQueue(ctx, queueID)
	if err != nil {
		t.Fatalf("GetQueue failed: %v", err)
	}
	t.Logf("GetQueue: %+v", gotQueue)
	if gotQueue.Pipe.String() != pipeID || gotQueue.Weight != "50" {
		t.Fatalf("GetQueue mismatch: %+v", gotQueue)
	}

	queue.Number = gotQueue.Number
	queue.Weight = "75"
	if err := controller.UpdateQueue(ctx, queueID, queue); err != nil {
		t.Fatalf("UpdateQueue failed: %v", err)
	}
	if gotQueue, err = controller.GetQueue(ctx, queueID); err != nil || gotQueue.Weight != "75" {
		t.Fatalf("UpdateQueue not applied: %+v (err %v)", gotQueue, err)
	}

	rule := &Rule{
		Enabled:         "1",
		Sequence:        "10",
		Interface:       api.SelectedMap("wan"),
		Protocol:        api.SelectedMap("ip"),
		Source:          "any",
		SourceNot:       "0",
		SourcePort:      "any",
		Destination:     "192.0.2.0/24",
		DestinationNot:  "0",
		DestinationPort: "any",
		DSCP:            api.SelectedMapList{},
		Direction:       api.SelectedMap("out"),
		Target:          api.SelectedMap(queueID),
		Description:     "test-rule",
	}

	ruleID, err := controller.AddRule(ctx, rule)
	if err != nil {
		t.Fatalf("AddRule failed: %v", err)
	}
	t.Logf("AddRule: uuid=%s", ruleID)

	gotRule, err := controller.GetRule(ctx, ruleID)
	if err != nil {
		t.Fatalf("GetRule failed: %v", err)
	}
	t.Logf("GetRule: %+v", gotRule)
	if gotRule.Target.String() != queueID || gotRule.Destination != "192.0.2.0/24" || gotRule.Direction.String() != "out" {
		t.Fatalf("GetRule mismatch: %+v", gotRule)
	}

	rule.DSCP = api.SelectedMapList{"ef", "af41"}
	rule.Target = api.SelectedMap(pipeID)
	if err := controller.UpdateRule(ctx, ruleID, rule); err != nil {
		t.Fatalf("UpdateRule failed: %v", err)
	}
	if gotRule, err = controller.GetRule(ctx, ruleID); err != nil || gotRule.Target.String() != pipeID || len(gotRule.DSCP) != 2 {
		t.Fatalf("UpdateRule not applied: %+v (err %v)", gotRule, err)
	}

	if err := controller.DeleteRule(ctx, ruleID); err != nil {
		t.Fatalf("DeleteRule failed: %v", err)
	}
	t.Log("DeleteRule: deleted")
}
