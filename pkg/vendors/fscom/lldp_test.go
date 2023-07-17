package fscom_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fscom"
	"testing"
)

var lldpResponse = `switch05_tw#show lldp neighbors         
Capability Codes:
       (R)Router,(B)Bridge,(C)DOCSIS Cable Device,(T)Telephone
       (W)WLAN Access Point, (P)Repeater,(S)Station,(O)Other

Device-ID       Local-Intf      Hldtme       Port-ID      Capability
core01_tw       TGi0/25         117          518          R B 
core02_tw       TGi0/26         117          518          R B 

Total entries displayed: 1
`

func TestParseLLDPTable(t *testing.T) {
	table, err := fscom.ParseLLDPTable(lldpResponse)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(table) != 2 {
		t.Errorf("expected 2 entries, got %d", len(table))
	}

	if table[0].LocalInterface != "TGigaEthernet0/25" {
		t.Errorf("expected port TGigaEthernet0/25, got %s", table[0].LocalInterface)
	}

	if table[0].RemoteHostname != "core01_tw" {
		t.Errorf("expected hostname core01_tw, got %s", table[0].RemoteHostname)
	}

	if table[1].LocalInterface != "TGigaEthernet0/26" {
		t.Errorf("expected port TGigaEthernet0/26, got %s", table[1].LocalInterface)
	}

	if table[1].RemoteHostname != "core02_tw" {
		t.Errorf("expected hostname core02_tw, got %s", table[1].RemoteHostname)
	}
}
