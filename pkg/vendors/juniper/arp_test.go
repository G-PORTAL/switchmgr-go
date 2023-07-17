package juniper_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper"
	"golang.org/x/exp/slices"
	"testing"
)

var arpResponse = `<rpc-reply
	xmlns:junos="http://xml.juniper.net/junos/14.1X53/junos" message-id="b3deff5e-55a5-41fa-838d-6045810b0b95"
	xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
	<ethernet-switching-table-information
		xmlns="http://xml.juniper.net/junos/14.1X53/junos-esw" junos:style="brief">
		<ethernet-switching-table junos:style="brief">
			<mac-table-count>217</mac-table-count>
			<mac-table-learned>206</mac-table-learned>
			<mac-table-persistent>0</mac-table-persistent>
			<mac-table-entry junos:style="brief">
				<mac-vlan>intern_daemon</mac-vlan>
				<mac-address>*</mac-address>
				<mac-type>Flood</mac-type>
				<mac-age>-</mac-age>
				<mac-interfaces-list>
					<mac-interfaces>All-members</mac-interfaces>
				</mac-interfaces-list>
			</mac-table-entry>
			<mac-table-entry junos:style="brief">
				<mac-vlan>intern_daemon</mac-vlan>
				<mac-address>00:1b:21:74:34:20</mac-address>
				<mac-type>Learn</mac-type>
				<mac-age junos:seconds="0">0</mac-age>
				<mac-interfaces-list>
					<mac-interfaces>ge-0/0/0.0</mac-interfaces>
				</mac-interfaces-list>
			</mac-table-entry>
			<mac-table-entry junos:style="brief">
				<mac-vlan>platform-mgmt</mac-vlan>
				<mac-address>d2:99:fe:b0:56:11</mac-address>
				<mac-type>Learn</mac-type>
				<mac-age junos:seconds="0">0</mac-age>
				<mac-interfaces-list>
					<mac-interfaces>ge-0/0/0.0</mac-interfaces>
				</mac-interfaces-list>
			</mac-table-entry>
			<mac-table-entry junos:style="brief">
				<mac-vlan>platform-mgmt</mac-vlan>
				<mac-address>f4:02:70:f9:5f:ee</mac-address>
				<mac-type>Learn</mac-type>
				<mac-age junos:seconds="0">0</mac-age>
				<mac-interfaces-list>
					<mac-interfaces>ge-0/0/1.0</mac-interfaces>
				</mac-interfaces-list>
			</mac-table-entry>
		</ethernet-switching-table>
	</ethernet-switching-table-information>
</rpc-reply>
`

func TestParseArpTable(t *testing.T) {
	entries, err := juniper.ParseArpTable([]byte(arpResponse))
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(entries) != 2 {
		t.Errorf("expected 4 entries, got %d", len(entries))
	}

	ports := []string{}
	for _, entry := range entries {
		ports = append(ports, entry.SwitchPort)
	}
	if !slices.Contains(ports, "ge-0/0/0") {
		t.Errorf("expected port ge-0/0/0 being present, got %v", ports)
	}

	if !slices.Contains(ports, "ge-0/0/1") {
		t.Errorf("expected port ge-0/0/1 being present, got %v", ports)
	}

	for _, entry := range entries {
		if entry.SwitchPort == "ge-0/0/0" {
			if len(entry.MacAddresses) != 2 {
				t.Errorf("expected 2 macs on ge-0/0/0, got %v", len(entry.MacAddresses))
			}
		}
		if entry.SwitchPort == "ge-0/0/1" {
			if len(entry.MacAddresses) != 1 {
				t.Errorf("expected 3 macs on ge-0/0/1, got %v", len(entry.MacAddresses))
			}
		}
	}
}
