package els_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper/els"
	"golang.org/x/exp/slices"
	"testing"
)

var arpResponse = `<rpc-reply
	xmlns:junos="http://xml.juniper.net/junos/18.2R3/junos" message-id="913ee354-577d-4215-a819-8b02b027c7ee"
	xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
	<l2ng-l2ald-rtb-macdb>
		<l2ng-l2ald-mac-entry-vlan junos:style="brief-rtb">
			<mac-count-global>1006</mac-count-global>
			<learnt-mac-count>1006</learnt-mac-count>
			<l2ng-l2-mac-routing-instance>default-switch</l2ng-l2-mac-routing-instance>
			<l2ng-l2-vlan-id>10</l2ng-l2-vlan-id>
			<l2ng-mac-entry>
				<l2ng-l2-mac-vlan-name>bm-backend</l2ng-l2-mac-vlan-name>
				<l2ng-l2-mac-address>72:3c:3c:b4:49:eb</l2ng-l2-mac-address>
				<l2ng-l2-mac-flags>D</l2ng-l2-mac-flags>
				<l2ng-l2-mac-age>-</l2ng-l2-mac-age>
				<l2ng-l2-mac-logical-interface>ae0.0</l2ng-l2-mac-logical-interface>
				<l2ng-l2-mac-fwd-next-hop>0</l2ng-l2-mac-fwd-next-hop>
				<l2ng-l2-mac-rtr-id>0</l2ng-l2-mac-rtr-id>
			</l2ng-mac-entry>
			<l2ng-mac-entry>
				<l2ng-l2-mac-vlan-name>bm-kvm</l2ng-l2-mac-vlan-name>
				<l2ng-l2-mac-address>f8:f2:1e:55:68:80</l2ng-l2-mac-address>
				<l2ng-l2-mac-flags>D</l2ng-l2-mac-flags>
				<l2ng-l2-mac-age>-</l2ng-l2-mac-age>
				<l2ng-l2-mac-logical-interface>ae0.0</l2ng-l2-mac-logical-interface>
				<l2ng-l2-mac-fwd-next-hop>0</l2ng-l2-mac-fwd-next-hop>
				<l2ng-l2-mac-rtr-id>0</l2ng-l2-mac-rtr-id>
			</l2ng-mac-entry>
			<l2ng-mac-entry>
				<l2ng-l2-mac-vlan-name>switch</l2ng-l2-mac-vlan-name>
				<l2ng-l2-mac-address>fc:96:43:c0:c5:f6</l2ng-l2-mac-address>
				<l2ng-l2-mac-flags>D</l2ng-l2-mac-flags>
				<l2ng-l2-mac-age>-</l2ng-l2-mac-age>
				<l2ng-l2-mac-logical-interface>ae0.0</l2ng-l2-mac-logical-interface>
				<l2ng-l2-mac-fwd-next-hop>0</l2ng-l2-mac-fwd-next-hop>
				<l2ng-l2-mac-rtr-id>0</l2ng-l2-mac-rtr-id>
			</l2ng-mac-entry>
			<l2ng-mac-entry>
				<l2ng-l2-mac-vlan-name>bm-public-produktiv</l2ng-l2-mac-vlan-name>
				<l2ng-l2-mac-address>2c:ea:7f:d6:9d:6c</l2ng-l2-mac-address>
				<l2ng-l2-mac-flags>D</l2ng-l2-mac-flags>
				<l2ng-l2-mac-age>-</l2ng-l2-mac-age>
				<l2ng-l2-mac-logical-interface>ge-0/0/20.0</l2ng-l2-mac-logical-interface>
				<l2ng-l2-mac-fwd-next-hop>0</l2ng-l2-mac-fwd-next-hop>
				<l2ng-l2-mac-rtr-id>0</l2ng-l2-mac-rtr-id>
			</l2ng-mac-entry>
		</l2ng-l2ald-mac-entry-vlan>
	</l2ng-l2ald-rtb-macdb>
</rpc-reply>
`

func TestParseArpTable(t *testing.T) {
	entries, err := els.ParseArpTable([]byte(arpResponse))
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
	if !slices.Contains(ports, "ae0") {
		t.Errorf("expected port ae0 being present, got %v", ports)
	}

	if !slices.Contains(ports, "ge-0/0/20") {
		t.Errorf("expected port ge-0/0/20 being present, got %v", ports)
	}

	for _, entry := range entries {
		if entry.SwitchPort == "ae0" {
			if len(entry.MacAddresses) != 3 {
				t.Errorf("expected 3 macs on ae0, got %v", len(entry.MacAddresses))
			}
		}
		if entry.SwitchPort == "ge-0/0/20" {
			if len(entry.MacAddresses) != 1 {
				t.Errorf("expected 3 macs on ge-0/0/20, got %v", len(entry.MacAddresses))
			}
		}
	}
}
