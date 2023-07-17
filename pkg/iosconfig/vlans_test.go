package iosconfig_test

import (
	"github.com/g-portal/switchmgr-go/pkg/iosconfig"
	"testing"
)

func TestConfig_Vlans(t *testing.T) {
	//cfg := iosconfig.
	cfg, err := iosconfig.Parse(plainConfig)
	if err != nil {
		t.Error(err)
	}

	vlans := cfg.Vlans()

	if len(vlans) != 3 {
		t.Errorf("expected 3 vlans, got %d", len(vlans))
	}

	ids := cfg.VlanIDs()
	if len(ids) != 3 {
		t.Errorf("expected 3 vlan ids, got %d", len(ids))
	}

	if ids[0] != 1 {
		t.Errorf("expected vlan id 1, got %d", ids[0])
	}

	if ids[1] != 4 {
		t.Errorf("expected vlan id 4, got %d", ids[1])
	}

	if ids[2] != 6 {
		t.Errorf("expected vlan id 6, got %d", ids[2])
	}

}
