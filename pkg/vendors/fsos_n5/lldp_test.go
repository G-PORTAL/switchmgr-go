package fsos_n5_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_n5"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_n5/utils"
	"testing"
)

func TestParseLLDPTable(t *testing.T) {
	table, err := fsos_n5.ParseLLDPNeighbors(utils.ReadTestData("show lldp neighbors", nil))
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(table) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(table))
	}

	if table[0].LocalInterface != "TenGigabitEthernet 0/47" {
		t.Errorf("expected port TenGigabitEthernet 0/47, got %s", table[0].LocalInterface)
	}

	if table[0].RemoteHostname != "stl-ds-07" {
		t.Errorf("expected hostname AMI9C6B001D74D4, got %s", table[0].RemoteHostname)
	}

	if table[1].LocalInterface != "HundredGigabitEthernet 0/55" {
		t.Errorf("expected port HundredGigabitEthernet 0/55, got %s", table[1].LocalInterface)
	}

	if table[1].RemoteHostname != "openshift-core1" {
		t.Errorf("expected hostname core02, got %s", table[1].RemoteHostname)
	}
}
