package dnsmasq

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestGeneralSettingsLargeResponse(t *testing.T) {
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

	for i := 1; i <= 4; i++ {
		id, err := controller.AddOption(ctx, &Option{
			Type:        api.SelectedMap("set"),
			OptionV4:    api.SelectedMap("1"),
			Value:       fmt.Sprintf("255.255.255.%d", i),
			Description: fmt.Sprintf("large-response-%d", i),
		})
		if err != nil {
			t.Fatalf("AddOption %d failed: %v", i, err)
		}
		defer controller.DeleteOption(ctx, id)
	}

	if _, err := controller.GeneralSettingsGet(ctx); err != nil {
		t.Fatalf("GeneralSettingsGet with a response over 64 KiB failed: %v", err)
	}
}
