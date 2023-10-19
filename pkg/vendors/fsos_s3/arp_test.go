package fsos_s3_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s3"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s3/utils"
	"testing"
)

func TestParseArpTable(t *testing.T) {
	entries, err := fsos_s3.ParseArpTable(utils.ReadTestData("show mac address table", nil))
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(entries) != 2 {
		t.Errorf("expected 4 entries, got %d", len(entries))
	}

	if entries[0].SwitchPort != "TGigaEthernet0/25" {
		t.Errorf("expected port TGigaEthernet0/25, got %s", entries[0].SwitchPort)
	}

	if len(entries[0].MacAddresses) != 5 {
		t.Errorf("expected 5 macs on tg0/25, got %v", len(entries[0].MacAddresses))
	}

	if entries[1].SwitchPort != "GigaEthernet0/6" {
		t.Errorf("expected port GigaEthernet0/6, got %s", entries[1].SwitchPort)
	}

	if len(entries[1].MacAddresses) != 1 {
		t.Errorf("expected 5 macs on g0/6, got %v", len(entries[1].MacAddresses))
	}
}
