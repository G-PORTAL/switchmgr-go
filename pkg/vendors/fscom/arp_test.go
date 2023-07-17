package fscom_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fscom"
	"testing"
)

var arpResponse = `switch05_tw#show mac address-table dynamic          
        Mac Address Table (Total 70)
------------------------------------------

Vlan    Mac Address       Type       Ports
----    -----------       ----       -----
4	32f3.614e.cd02	  DYNAMIC    tg0/25
6	507c.6f2a.bc2a	  DYNAMIC    g0/6
1	80ac.acbd.8585	  DYNAMIC    tg0/25
1	80ac.acbd.8580	  DYNAMIC    tg0/25
1	80ac.acbd.8582	  DYNAMIC    tg0/25
6	921e.bca0.8f7f	  DYNAMIC    tg0/25`

func TestParseArpTable(t *testing.T) {
	entries, err := fscom.ParseArpTable(arpResponse)
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
