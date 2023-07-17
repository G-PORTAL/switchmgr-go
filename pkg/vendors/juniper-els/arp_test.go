package juniper_els_test

import (
	juniper_els "github.com/g-portal/switchmgr-go/pkg/vendors/juniper-els"
	"golang.org/x/exp/slices"
	"testing"
)

func TestListArpTable(t *testing.T) {
	driver := juniper_els.NewMockDriver()
	entries, err := driver.ListArpTable()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	var ports []string
	for _, entry := range entries {
		ports = append(ports, entry.SwitchPort)
	}
	if !slices.Contains(ports, "ae0") {
		t.Errorf("expected port ae0 being present, got %v", ports)
	}

	if !slices.Contains(ports, "ge-0/0/20") {
		t.Errorf("expected port ge-0/0/20 being present, got %v", ports)
	}

	for _, entry := range entries {
		if entry.SwitchPort == "ae0" {
			if len(entry.MacAddresses) != 1 {
				t.Errorf("expected 1 macs on ae0, got %v", len(entry.MacAddresses))
			}
		}
		if entry.SwitchPort == "ge-0/0/20" {
			if len(entry.MacAddresses) != 1 {
				t.Errorf("expected 1 macs on ge-0/0/20, got %v", len(entry.MacAddresses))
			}
		}
	}
}
