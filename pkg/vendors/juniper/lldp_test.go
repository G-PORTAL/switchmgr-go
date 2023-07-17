package juniper_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper"
	"testing"
)

func TestListLLDPNeighbors(t *testing.T) {
	driver := juniper.NewMockDriver()
	neighbors, err := driver.ListLLDPNeighbors()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(neighbors) != 1 {
		t.Errorf("expected 1 neighbor, got %d", len(neighbors))
	}

	for _, n := range neighbors {
		if n.LocalInterface == "ge-0/0/0.0" {
			if n.RemoteHostname != "switch-host-name" {
				t.Errorf("interface %s: expected hostname switch73_Frankfurt, got %s", n.LocalInterface, n.RemoteHostname)
			}
		}
	}
}
