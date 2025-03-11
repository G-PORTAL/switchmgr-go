package arista_eos_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/arista_eos"
	"github.com/g-portal/switchmgr-go/pkg/vendors/arista_eos/utils"
	"testing"
)

func TestParseLLDPTable(t *testing.T) {
	var response arista_eos.AristaLLDPResponse
	err := utils.ReadTestData("show lldp neighbors", &response)
	if err != nil {
		t.Error(err)
	}

	table, err := arista_eos.ParseLLDPNeighbors(response)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(table) != 3 {
		t.Errorf("expected 3 entries, got %d", len(table))
	}

	if table[0].LocalInterface != "Ethernet51" {
		t.Errorf("expected port Ethernet51, got %s", table[0].LocalInterface)
	}

	if table[0].RemoteHostname != "core1.host.com" {
		t.Errorf("expected hostname core1.host.com, got %s", table[0].RemoteHostname)
	}

	if table[1].LocalInterface != "Ethernet52" {
		t.Errorf("expected port Ethernet52, got %s", table[1].LocalInterface)
	}

	if table[1].RemoteHostname != "core2.host.com" {
		t.Errorf("expected hostname core2.host.com, got %s", table[1].RemoteHostname)
	}

	if table[2].LocalInterface != "Management1" {
		t.Errorf("expected port Management1, got %s", table[2].LocalInterface)
	}

	if table[2].RemoteHostname != "switch2.host.com" {
		t.Errorf("expected hostname switch2.host.com, got %s", table[2].RemoteHostname)
	}
}
