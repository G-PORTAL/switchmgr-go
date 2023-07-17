package fscom_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fscom"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fscom/utils"
	"testing"
)

func TestListInterfaces(t *testing.T) {
	iosConfig, err := fscom.ParseConfiguration(utils.ReadTestData("show running-config"))
	if err != nil {
		t.Fatal(err)
	}

	cfg := fscom.Configuration(iosConfig)

	nics, err := cfg.ListInterfaces()
	if err != nil {
		t.Fatal(err)
	}

	if len(nics) != 30 {
		t.Fatal("Expected 30 interfaces, got", len(nics))
	}

	_, err = cfg.GetInterface("GigaEthernet0/124")
	if err == nil {
		t.Fatal("Expected error for non-existing interface")
	}

	nic, err := cfg.GetInterface("GigaEthernet0/19")
	if err != nil {
		t.Fatal(err)
	}

	if !nic.Enabled {
		t.Fatal("Expected interface to be enabled")
	}

	if nic.Name != "GigaEthernet0/19" {
		t.Fatal("Expected interface name to be GigaEthernet0/19, got", nic.Name)
	}

	if nic.UntaggedVLAN == nil {
		t.Fatal("Expected untagged VLAN to be set")
	}

	if *nic.UntaggedVLAN != 6 {
		t.Fatal("Expected untagged VLAN to be 6, got", *nic.UntaggedVLAN)
	}

	nic, err = cfg.GetInterface("GigaEthernet0/17")
	if err != nil {
		t.Fatal(err)
	}

	if nic.Enabled {
		t.Fatal("Expected interface to be disabled")
	}

	nic, err = cfg.GetInterface("GigaEthernet0/22")
	if err != nil {
		t.Fatal(err)
	}

	if nic.UntaggedVLAN == nil || *nic.UntaggedVLAN != 1 {
		t.Fatal("Expected untagged VLAN to be 1")
	}

	nic, err = cfg.GetInterface("GigaEthernet0/23")
	if err != nil {
		t.Fatal(err)
	}

	if nic.UntaggedVLAN == nil || *nic.UntaggedVLAN != 6 {
		t.Fatal("Expected untagged VLAN to be 6")
	}

	nic, err = cfg.GetInterface("GigaEthernet0/18")
	if err != nil {
		t.Fatal(err)
	}

	if nic.UntaggedVLAN == nil || *nic.UntaggedVLAN != 6 {
		t.Fatal("Expected untagged VLAN to be 6")
	}

	if len(nic.TaggedVLANs) != 2 {
		t.Fatal("Expected 2 tagged VLANs, got", len(nic.TaggedVLANs))
	}

	nic, err = cfg.GetInterface("GigaEthernet0/24")
	if err != nil {
		t.Fatal(err)
	}

	if nic.UntaggedVLAN == nil || *nic.UntaggedVLAN != 21 {
		t.Fatal("Expected untagged VLAN to be 21")
	}

	if len(nic.TaggedVLANs) != 6 {
		t.Fatalf("Expected 2 tagged VLANs, got %+v", nic.TaggedVLANs)
	}

	if nic.TaggedVLANs[0] != 2 {
		t.Fatalf("Expected tagged VLAN 0 to be 2, got %d", nic.TaggedVLANs[0])
	}

	if nic.TaggedVLANs[1] != 4 {
		t.Fatalf("Expected tagged VLAN 1 to be 4, got %d", nic.TaggedVLANs[1])
	}

	if nic.TaggedVLANs[2] != 22 {
		t.Fatalf("Expected tagged VLAN 2 to be 22, got %d", nic.TaggedVLANs[2])
	}

	if nic.Management != false {
		t.Fatalf("Expected management to be false, got %t", nic.Management)
	}

	nic, err = cfg.GetInterface("VLAN4")
	if err != nil {
		t.Fatal(err)
	}

	if nic.Management != true {
		t.Fatalf("Expected management to be true, got %t", nic.Management)
	}

}

func TestParseInterfaces(t *testing.T) {
	nics, err := fscom.ParseInterfaces(utils.ReadTestData("show interfaces"))
	if err != nil {
		t.Fatal(err)
	}

	if len(nics) != 29 {
		t.Fatalf("expected 29 nics, got %d", len(nics))
	}

	info, ok := nics["TGigaEthernet0/25"]
	if !ok {
		t.Fatalf("interface TGigaEthernet0/25 not found")
	}

	if info.MacAddress.String() != "64:9d:99:c5:aa:aa" {
		t.Fatalf("expected mac 64:9d:99:c5:aa:aa, got %s", info.MacAddress.String())
	}

	if info.MTU != 1500 {
		t.Fatalf("expected MTU 1500, got %d", info.MTU)
	}

	if info.Speed != 10000000 {
		t.Fatalf("expected Speed 10000000, got %d", info.Speed)
	}

	_, ok = nics["TGigaEthernet0/123"]
	if ok {
		t.Fatalf("interface TGigaEthernet0/123 found, but should not exist")
	}

}
