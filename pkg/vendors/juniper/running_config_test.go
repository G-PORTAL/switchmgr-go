package juniper_test

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper"
	"testing"
)

func TestGetRunningConfig(t *testing.T) {
	driver := juniper.NewMockDriver()
	cfg, err := driver.GetRunningConfig()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(cfg.VLANs) != 2 {
		t.Errorf("expected model 2 vlans, got %d", len(cfg.VLANs))
	}

	// Ensure vlan data to be correct
	expectedVlans := map[string]int32{
		"vlan1": 1,
		"vlan2": 2,
	}

	for _, vlan := range cfg.VLANs {
		if desiredState, ok := expectedVlans[vlan.Name]; ok {
			if vlan.ID != desiredState {
				t.Errorf("vlan %s has wrong vlan-id, expected %d, got %d", vlan.Name, desiredState, vlan.ID)
			}
		} else {
			t.Errorf("vlan %s (%d) not found in expected vlans", vlan.Name, vlan.ID)
		}
	}

	for name, id := range expectedVlans {
		if vlan, err := cfg.GetVlanIDByName(name); err == nil {
			if vlan != id {
				t.Errorf("vlan %s has wrong vlan-id, expected %d, got %d", name, id, vlan)
			}
		} else {
			t.Errorf("vlan %s not found in config", name)
		}
	}

	// Test if the GetInterfaceMode utility works
	expectedModes := map[string]models.InterfaceMode{
		"ge-0/0/0": models.InterfaceModeTrunk,
		"ge-0/0/1": models.InterfaceModeTrunk,
		"ge-0/0/2": models.InterfaceModeAccess,
	}
	for iface, mode := range expectedModes {
		if actualMode, err := cfg.GetInterfaceMode(iface); err == nil {
			if actualMode != mode {
				t.Errorf("interface %s has wrong mode, expected %s, got %s", iface, mode, actualMode)
			}
		} else {
			t.Errorf("interface %s not found in config", iface)
		}
	}
}
