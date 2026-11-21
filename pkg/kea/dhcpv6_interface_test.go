package kea

import (
	"context"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func withDhcpv6Interface(t *testing.T, c *Controller, iface string) {
	t.Helper()
	ctx := context.Background()
	cur, err := c.Dhcpv6ConfigGet(ctx)
	if err != nil {
		t.Fatalf("Dhcpv6ConfigGet: %v", err)
	}
	orig := cur.Monad
	set := orig
	set.General.Interfaces = api.SelectedMapList{iface}
	if _, err := c.Dhcpv6ConfigSet(ctx, &set); err != nil {
		t.Fatalf("Dhcpv6ConfigSet(%s): %v", iface, err)
	}
	t.Cleanup(func() {
		if _, err := c.Dhcpv6ConfigSet(context.Background(), &orig); err != nil {
			t.Logf("restoring the DHCPv6 interfaces: %v", err)
		}
	})
}
