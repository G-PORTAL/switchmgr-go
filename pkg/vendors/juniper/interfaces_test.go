package juniper_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper"
	"golang.org/x/exp/slices"
	"net"
	"testing"
)

func TestListInterfaces(t *testing.T) {
	driver := juniper.NewMockDriver()
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
	driver := juniper.NewMockDriver()
	iface, err := driver.GetInterface("ge-0/0/0")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if iface.Name != "ge-0/0/0" {
		t.Errorf("expected interface ge-0/0/0, got %s", iface.Name)
	}
	if iface.MacAddress != "28:c0:da:2c:07:b1" {
		t.Errorf("expected mac address 28:c0:da:2c:07:b1, got %s", iface.MacAddress)
	}
	if !slices.Contains(iface.TaggedVLANs, 4) {
		t.Errorf("expected interface ge-0/0/0 being member of vlan 4, got %v", iface.TaggedVLANs)
	}
}
