package juniper_els_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper-els"
	"testing"
)

func TestListLLDPNeighbors(t *testing.T) {
	driver := juniper_els.NewMockDriver()
	neighbors, err := driver.ListLLDPNeighbors()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(neighbors) != 2 {
		t.Errorf("expected 2 neighbor, got %d", len(neighbors))
	}

	found := false
	for _, n := range neighbors {
		if n.LocalInterface == "xe-0/2/1" {
			found = true
			if n.RemoteHostname != "switch-system-name" {
				t.Errorf("expected %s being connected to switch-system-name, got %s", n.LocalInterface, n.RemoteHostname)
			}
		}
	}

	if !found {
		t.Errorf("expected xe-0/2/1 being present, got %v", neighbors)
	}
}
