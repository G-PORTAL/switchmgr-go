package fsos_s5_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s5"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s5/utils"
	"testing"
)

func TestParseLLDPTable(t *testing.T) {
	table, err := fsos_s5.ParseLLDPNeighbors(utils.ReadTestData("show lldp neighbor brief", nil))
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(table) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(table))
	}

	if table[0].LocalInterface != "eth-0-17" {
		t.Errorf("expected port eth-0-17, got %s", table[0].LocalInterface)
	}

	if table[0].RemoteHostname != "AMI9C6B001D74D4" {
		t.Errorf("expected hostname AMI9C6B001D74D4, got %s", table[0].RemoteHostname)
	}

	if table[2].LocalInterface != "eth-0-52" {
		t.Errorf("expected port TGigaEthernet0/26, got %s", table[1].LocalInterface)
	}

	if table[2].RemoteHostname != "wup-ds-08" {
		t.Errorf("expected hostname core02, got %s", table[1].RemoteHostname)
	}
}
