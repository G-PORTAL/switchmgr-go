package juniper_els_test

import (
	"bytes"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper-els"
	"golang.org/x/exp/slices"
	"net"
	"testing"
	"text/template"
)

func TestListInterfaces(t *testing.T) {
	driver := juniper_els.NewMockDriver()
	interfaces, err := driver.ListInterfaces()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(interfaces) != 3 {
		t.Errorf("expected 3 interface, got %d", len(interfaces))
	}

	expectedMTUs := map[string]uint32{
		"ge-0/0/0": 1514,
		"ge-0/0/1": 1337,
	}

	for _, m := range interfaces {
		if m.MacAddress == "" {
			continue
		}
		if desiredMTU, ok := expectedMTUs[m.Name]; ok {
			if m.MTU != desiredMTU {
				t.Errorf("interface %s has wrong MTU, expected %d, got %d", m.Name, desiredMTU, m.MTU)
			}
		}

		if _, err = net.ParseMAC(string(m.MacAddress)); err != nil {
			t.Errorf("interface %s: invalid mac address: %s", m.Name, err.Error())
		}
	}
}
func TestGetInterface(t *testing.T) {
	driver := juniper_els.NewMockDriver()
	iface, err := driver.GetInterface("ge-0/0/0")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if iface.Name != "ge-0/0/0" {
		t.Errorf("expected interface ge-0/0/0, got %s", iface.Name)
	}
	if iface.MacAddress != "9c:8a:cb:18:f6:aa" {
		t.Errorf("expected mac address 9c:8a:cb:18:f6:aa, got %s", iface.MacAddress)
	}
	if *iface.UntaggedVLAN != 1 {
		t.Errorf("expected untagged VLAN to be 1, got %v", iface.UntaggedVLAN)
	}
	if !slices.Contains(iface.TaggedVLANs, 4) {
		t.Errorf("expected interface ge-0/0/0 being member of vlan 4, got %v", iface.TaggedVLANs)
	}
	if !iface.Enabled {
		t.Errorf("interface ge-0/0/0 should be enabled")
	}
	if iface.Management {
		t.Errorf("interface ge-0/0/0 should not be managed")
	}
	if iface.Speed != 1_000_000 {
		t.Errorf("Speed should be 1GBit, got %v", iface.Speed)
	}
	if iface.MTU != 1514 {
		t.Errorf("MTU should be 1514, got %v", iface.MTU)
	}
}

const EditPortConfigurationExpected = `<edit-config>
	<target>
		<candidate/>
	</target>
    <default-operation>merge</default-operation>
    <test-option>test-then-set</test-option>
	<config>
		<configuration>
			<interfaces>
				<interface operation="replace">
					<name>eth0</name>
					<description>example interface</description>
					<native-vlan-id>1337</native-vlan-id>
					<unit>
						<name>0</name>
						<family>
							<ethernet-switching>
								<interface-mode>trunk</interface-mode>
								<vlan>
									<members>1</members><members>2</members><members>3</members>
								</vlan>
							</ethernet-switching>
						</family>       
					</unit>             
				</interface>   
			</interfaces>
		</configuration>
	</config>
</edit-config>
`

func TestConfigureInterfaceTemplate(t *testing.T) {
	var tpl bytes.Buffer
	tmpl, err := template.New("").Parse(juniper_els.EditPortConfigurationTemplate)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	interfaceDescription := "example interface"
	untaggedVlan := int32(1337)
	if err = tmpl.Execute(&tpl, &models.UpdateInterface{
		Name:         "eth0",
		Description:  &interfaceDescription,
		UntaggedVLAN: &untaggedVlan,
		TaggedVLANs:  []int32{1, 2, 3},
	}); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if tpl.String() != EditPortConfigurationExpected {
		t.Errorf("expected:\n%s\ngot:\n%s", EditPortConfigurationExpected, tpl.String())
	}
}
