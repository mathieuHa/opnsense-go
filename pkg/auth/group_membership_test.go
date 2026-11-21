package auth

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestGroupUpdateKeepsMembers(t *testing.T) {
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

	groupID, err := controller.AddGroup(ctx, &Group{Name: "tfmembergrp", Description: "membership"})
	if err != nil {
		t.Fatalf("AddGroup failed: %v", err)
	}
	defer controller.DeleteGroup(ctx, groupID)
	group, err := controller.GetGroup(ctx, groupID)
	if err != nil {
		t.Fatalf("GetGroup failed: %v", err)
	}

	userID, err := controller.AddUser(ctx, &User{
		Name:             "tfmemberuser",
		Disabled:         "0",
		Password:         "Membership-Test-1!",
		GroupMemberships: api.SelectedMapList{group.GroupId},
	})
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}
	defer controller.DeleteUser(ctx, userID)

	if err := controller.UpdateGroup(ctx, groupID, &Group{Name: "tfmembergrp", Description: "membership-updated"}); err != nil {
		t.Fatalf("UpdateGroup failed: %v", err)
	}

	user, err := controller.GetUser(ctx, userID)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if len(user.GroupMemberships) != 1 || user.GroupMemberships[0] != group.GroupId {
		t.Fatalf("updating the group removed its members: user group_memberships = %v, want [%s]", user.GroupMemberships, group.GroupId)
	}
}
