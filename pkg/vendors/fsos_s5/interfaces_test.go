package fsos_s5_test

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s5"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s5/utils"
	"testing"
)

func TestListInterfaces(t *testing.T) {
	iosConfig := fsos_s5.ParseConfiguration(utils.ReadTestData("show running-config", nil))

	cfg := fsos_s5.Configuration(iosConfig)

	nics, err := cfg.ListInterfaces()
	if err != nil {
		t.Fatal(err)
	}

	if len(nics) != 55 {
		t.Fatal("Expected 55 interfaces, got", len(nics))
	}

	_, err = cfg.GetInterface("GigaEthernet0/124")
	if err == nil {
		t.Fatal("Expected error for non-existing interface")
	}

	nic, err := cfg.GetInterface("eth-0-5")
	if err != nil {
		t.Fatal(err)
	}

	if !nic.Enabled {
		t.Fatalf("Expected interface %s to be enabled", nic.Name)
	}

	if nic.Name != "eth-0-5" {
		t.Fatalf("Expected interface %s name to be eth-0-5", nic.Name)
	}

	if nic.UntaggedVLAN == nil {
		t.Fatalf("Expected untagged VLAN to be set on %s", nic.Name)
	}

	nic, err = cfg.GetInterface("eth-0-9")
	if err != nil {
		t.Fatal(err)
	}

	if nic.Enabled {
		t.Fatalf("Expected interface %s to be disabled", nic.Name)
	}

	nic, err = cfg.GetInterface("vlan20")
	if err != nil {
		t.Fatal(err)
	}

	if !nic.Enabled {
		t.Fatalf("Expected interface %s to be enabled", nic.Name)
	}

	if nic.Mode != models.InterfaceModeAccess {
		t.Fatalf("Expected interface %s to be in access mode", nic.Name)
	}

	if nic.UntaggedVLAN == nil {
		t.Fatalf("Expected untagged VLAN to be set on %s", nic.Name)
	}

	untaggedVLAN := int32(20)
	if *nic.UntaggedVLAN != untaggedVLAN {
		t.Fatalf("Expected untagged VLAN to be %d, got %d", untaggedVLAN, *nic.UntaggedVLAN)
	}

}

func TestParseInterfaces(t *testing.T) {
	nics, err := fsos_s5.ParseInterfaces(utils.ReadTestData("show interface", nil))
	if err != nil {
		t.Fatal(err)
	}

	if len(nics) != 54 {
		t.Fatalf("expected 54 nics, got %d", len(nics))
	}

	if info, ok := nics["vlan20"]; !ok {
		t.Fatalf("interface vlan20 not found")
	} else {
		if info.MacAddress.String() != "64:9d:99:06:ff:33" {
			t.Fatalf("expected vlan20 mac 64:9d:99:06:ff:33, got %s", info.MacAddress.String())
		}

		if info.MTU != 1300 {
			t.Fatalf("expected vlan20 MTU 1300, got %d", info.MTU)
		}

		if info.Speed != 10000000 {
			t.Fatalf("expected vlan20 Speed 10000000, got %d", info.Speed)
		}
	}

	if info, ok := nics["eth-0-51"]; !ok {
		t.Fatalf("interface eth-0-51 not found")
	} else {
		if info.MacAddress.String() != "64:9d:99:06:ff:66" {
			t.Fatalf("expected eth-0-51 mac 64:9d:99:06:ff:66, got %s", info.MacAddress.String())
		}

		if info.MTU != 1500 {
			t.Fatalf("expected eth-0-51 MTU 1500, got %d", info.MTU)
		}

		if info.Speed != 10000000 {
			t.Fatalf("expected eth-0-51 Speed 10000000, got %d", info.Speed)
		}
	}

	if _, ok := nics["TGigaEthernet0/123"]; ok {
		t.Fatalf("interface TGigaEthernet0/123 found, but should not exist")
	}

}
