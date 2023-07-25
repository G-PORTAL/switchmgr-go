package juniper_els_test

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper-els"
	"testing"
)

func TestGetRunningConfig(t *testing.T) {
	driver := juniper_els.NewMockDriver()
	cfg, err := driver.GetRunningConfig()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(cfg.VLANs) != 2 {
		t.Errorf("expected model 2 vlans, got %d", len(cfg.VLANs))
	}
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

	// Test if the GetInterfaceMode utility works
	expectedModes := map[string]models.InterfaceMode{
		"ge-0/0/0":  models.InterfaceModeTrunk,
		"ge-0/0/1":  models.InterfaceModeAccess,
		"pfe-0/0/0": models.InterfaceModeAccess,
	}
	for iface, mode := range expectedModes {
		if actualMode, err := cfg.GetInterfaceMode(iface); err == nil {
			if actualMode != mode {
				t.Errorf("interface %s has wrong mode, expected %s, got %s", iface, mode, actualMode)
			}
		} else {
			if iface == "pfe-0/0/0" {
				continue
			}
			t.Errorf("interface %s not found in config", iface)
		}
	}
}
