package fsos_s3_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s3"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s3/utils"
	"testing"
)

func TestParseLLDPTable(t *testing.T) {
	table, err := fsos_s3.ParseLLDPNeighbors(utils.ReadTestData("show lldp neighbors", nil))
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(table) != 2 {
		t.Errorf("expected 2 entries, got %d", len(table))
	}

	if table[0].LocalInterface != "TGigaEthernet0/25" {
		t.Errorf("expected port TGigaEthernet0/25, got %s", table[0].LocalInterface)
	}

	if table[0].RemoteHostname != "core01" {
		t.Errorf("expected hostname core01, got %s", table[0].RemoteHostname)
	}

	if table[1].LocalInterface != "TGigaEthernet0/26" {
		t.Errorf("expected port TGigaEthernet0/26, got %s", table[1].LocalInterface)
	}

	if table[1].RemoteHostname != "core02" {
		t.Errorf("expected hostname core02, got %s", table[1].RemoteHostname)
	}
}
