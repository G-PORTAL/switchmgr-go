package juniper_els_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper-els"
	"golang.org/x/exp/slices"
	"net"
	"testing"
)

func TestListInterfaces(t *testing.T) {
	driver := juniper_els.NewMockDriver()
	interfaces, err := driver.ListInterfaces()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(interfaces) != 2 {
		t.Errorf("expected 2 interface, got %d", len(interfaces))
	}

	for _, m := range interfaces {
		if m.MacAddress == "" {
			continue
		}

		if _, err = net.ParseMAC(string(m.MacAddress)); err != nil {
			t.Errorf("interface %s: invalid mac address: %s", m.Name, err.Error())
		}
	}
}
func TestGetInterface(t *testing.T) {
	driver := juniper_els.NewMockDriver()
	iface, err := driver.GetInterface("ge-0/0/0")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if iface.Name != "ge-0/0/0" {
		t.Errorf("expected interface ge-0/0/0, got %s", iface.Name)
	}
	if iface.MacAddress != "9c:8a:cb:18:f6:aa" {
		t.Errorf("expected mac address 9c:8a:cb:18:f6:aa, got %s", iface.MacAddress)
	}
	if *iface.UntaggedVLAN != 1 {
		t.Errorf("expected untagged VLAN to be 1, got %v", iface.UntaggedVLAN)
	}
	if !slices.Contains(iface.TaggedVLANs, 4) {
		t.Errorf("expected interface ge-0/0/0 being member of vlan 4, got %v", iface.TaggedVLANs)
	}
	if !iface.Enabled {
		t.Errorf("interface ge-0/0/0 should be enabled")
	}
	if iface.Management {
		t.Errorf("interface ge-0/0/0 should not be managed")
	}
	if iface.Speed != 1_000_000 {
		t.Errorf("Speed should be 1GBit, got %v", iface.Speed)
	}
	if iface.MTU != 1514 {
		t.Errorf("MTU should be 1514, got %v", iface.MTU)
	}
}
