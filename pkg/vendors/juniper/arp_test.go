package juniper_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper"
	"golang.org/x/exp/slices"
	"testing"
)

func TestListArpTable(t *testing.T) {
	driver := juniper.NewMockDriver()
	entries, err := driver.ListArpTable()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(entries) != 1 {
		t.Errorf("expected 1 entries, got %d", len(entries))
	}

	var ports []string
	for _, entry := range entries {
		ports = append(ports, entry.SwitchPort)
	}
	if !slices.Contains(ports, "ge-0/0/0") {
		t.Errorf("expected port ge-0/0/0 being present, got %v", ports)
	}

	if !slices.Contains(ports, "ge-0/0/0") {
		t.Errorf("expected port ge-0/0/0 being present, got %v", ports)
	}

	for _, entry := range entries {
		if entry.SwitchPort == "ge-0/0/0" {
			if len(entry.MacAddresses) != 1 {
				t.Errorf("expected 1 macs on ge-0/0/0, got %v", len(entry.MacAddresses))
			}
		}
		if entry.SwitchPort == "ge-0/0/5" {
			if len(entry.MacAddresses) != 1 {
				t.Errorf("expected 1 macs on ge-0/0/5, got %v", len(entry.MacAddresses))
			}
		}
	}
}
