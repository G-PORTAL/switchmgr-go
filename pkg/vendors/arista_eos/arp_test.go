package arista_eos_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/arista_eos"
	"github.com/g-portal/switchmgr-go/pkg/vendors/arista_eos/utils"
	"testing"
)

func TestParseArpTable(t *testing.T) {
	var table arista_eos.AristaMacAddressTableResponse
	err := utils.ReadTestData("show mac address-table", &table)
	if err != nil {
		t.Error(err)
	}

	entries, err := arista_eos.ParseArpTable(table)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(entries) != 3 {
		t.Errorf("expected 4 entries, got %d", len(entries))
	}

	if entries[0].SwitchPort != "Ethernet11" {
		t.Errorf("expected port Ethernet11, got %s", entries[0].SwitchPort)
	}

	if len(entries[0].MacAddresses) != 1 {
		t.Errorf("expected 1 macs on Ethernet11, got %v", len(entries[0].MacAddresses))
	}

	if entries[1].SwitchPort != "Ethernet2" {
		t.Errorf("expected port Ethernet2, got %s", entries[1].SwitchPort)
	}

	if len(entries[1].MacAddresses) != 1 {
		t.Errorf("expected 1 macs on Ethernet2, got %v", len(entries[1].MacAddresses))
	}

	if entries[2].SwitchPort != "Ethernet31" {
		t.Errorf("expected port Ethernet31, got %s", entries[2].SwitchPort)
	}

	if len(entries[2].MacAddresses) != 2 {
		t.Errorf("expected 2 macs on Ethernet31, got %v", len(entries[2].MacAddresses))
	}
}
