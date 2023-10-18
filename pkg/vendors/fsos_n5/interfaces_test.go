package fsos_n5_test

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_n5"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_n5/utils"
	"testing"
)

func TestListInterfaces(t *testing.T) {
	iosConfig := fsos_n5.ParseConfiguration(utils.ReadTestData("show running-config", nil))

	cfg := fsos_n5.Configuration(iosConfig)

	nics, err := cfg.ListInterfaces()
	if err != nil {
		t.Fatal(err)
	}

	if len(nics) != 72 {
		t.Fatal("Expected 72 interfaces, got", len(nics))
	}

	_, err = cfg.GetInterface("TenGigabitEthernet 0/124")
	if err == nil {
		t.Fatal("Expected error for non-existing interface")
	}

	nic, err := cfg.GetInterface("TenGigabitEthernet 0/34")
	if err != nil {
		t.Fatal(err)
	}

	if !nic.Enabled {
		t.Fatalf("Expected interface %s to be enabled", nic.Name)
	}

	if nic.Name != "TenGigabitEthernet 0/34" {
		t.Fatalf("Expected interface %s name to be TenGigabitEthernet 0/34", nic.Name)
	}

	if nic.UntaggedVLAN == nil {
		t.Fatalf("Expected untagged VLAN to be set on %s", nic.Name)
	}

	untaggedVLAN := int32(4)
	if *nic.UntaggedVLAN != untaggedVLAN {
		t.Fatalf("Expected untagged VLAN to be %d, got %d", untaggedVLAN, *nic.UntaggedVLAN)
	}

	nic, err = cfg.GetInterface("AggregatePort 4")
	if err != nil {
		t.Fatal(err)
	}

	if !nic.Enabled {
		t.Fatalf("Expected interface %s to be disabled", nic.Name)
	}

	if nic.Mode != models.InterfaceModeTrunk {
		t.Fatalf("Expected interface %s to be in trunk mode", nic.Name)
	}

	if nic.UntaggedVLAN == nil {
		t.Fatalf("Expected untagged VLAN to be set on %s", nic.Name)
	}

	untaggedVLAN = int32(1900)
	if *nic.UntaggedVLAN != untaggedVLAN {
		t.Fatalf("Expected untagged VLAN to be %d, got %d", untaggedVLAN, *nic.UntaggedVLAN)
	}

	nic, err = cfg.GetInterface("VLAN 4")
	if err != nil {
		t.Fatal(err)
	}

	if !nic.Enabled {
		t.Fatalf("Expected interface %s to be disabled", nic.Name)
	}

	if nic.Mode != models.InterfaceModeAccess {
		t.Fatalf("Expected interface %s to be in access mode", nic.Name)
	}

	if nic.UntaggedVLAN == nil {
		t.Fatalf("Expected untagged VLAN to be set on %s", nic.Name)
	}

	untaggedVLAN = int32(4)
	if *nic.UntaggedVLAN != untaggedVLAN {
		t.Fatalf("Expected untagged VLAN to be %d, got %d", untaggedVLAN, *nic.UntaggedVLAN)
	}

}

func TestParseInterfaces(t *testing.T) {
	nics, err := fsos_n5.ParseInterfaces(utils.ReadTestData("show interfaces", nil))
	if err != nil {
		t.Fatal(err)
	}

	if len(nics) != 71 {
		t.Fatalf("expected 71 nics, got %d", len(nics))
	}

	if info, ok := nics["VLAN 4"]; !ok {
		t.Fatalf("interface VLAN 4 not found")
	} else {
		if info.MacAddress.String() != "64:9d:99:d2:50:2e" {
			t.Fatalf("expected vlan20 mac 64:9d:99:d2:50:2e, got %s", info.MacAddress.String())
		}

		if info.MTU != 1500 {
			t.Fatalf("expected VLAN 4 MTU 1500, got %d", info.MTU)
		}

		if info.Speed != 1000000 {
			t.Fatalf("expected VLAN 4 Speed 1000000, got %d", info.Speed)
		}
	}

	if info, ok := nics["TenGigabitEthernet 0/1"]; !ok {
		t.Fatalf("interface TenGigabitEthernet 0/1 not found")
	} else {
		if info.MacAddress.String() != "64:9d:99:d2:50:2d" {
			t.Fatalf("expected eth-0-51 mac 64:9d:99:d2:50:2d, got %s", info.MacAddress.String())
		}

		if info.MTU != 1500 {
			t.Fatalf("expected eth-0-51 MTU 1500, got %d", info.MTU)
		}

		if info.Speed != 10000000 {
			t.Fatalf("expected eth-0-51 Speed 10000000, got %d", info.Speed)
		}
	}

	if _, ok := nics["TenGigabitEthernet 0/100"]; ok {
		t.Fatalf("interface TenGigabitEthernet 0/100 found, but should not exist")
	}

}
