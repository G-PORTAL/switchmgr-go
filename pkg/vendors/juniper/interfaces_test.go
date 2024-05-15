package juniper_test

import (
	"bytes"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/utils"
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/slices"
	"net"
	"testing"
	"text/template"
)

func TestListInterfaces(t *testing.T) {
	driver := juniper.NewMockDriver()
	interfaces, err := driver.ListInterfaces()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(interfaces) != 2 {
		t.Errorf("expected 2 interface, got %d", len(interfaces))
	}

	for _, m := range interfaces {
		if m.MacAddress == "" {
			continue
		}

		if _, err = net.ParseMAC(string(m.MacAddress)); err != nil {
			t.Errorf("interface %s: invalid mac address: %s", m.Name, err.Error())
		}
	}
}
func TestGetInterface(t *testing.T) {
	driver := juniper.NewMockDriver()
	iface, err := driver.GetInterface("ge-0/0/0")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if iface.Name != "ge-0/0/0" {
		t.Errorf("expected interface ge-0/0/0, got %s", iface.Name)
	}
	if iface.MacAddress != "28:c0:da:2c:07:b1" {
		t.Errorf("expected mac address 28:c0:da:2c:07:b1, got %s", iface.MacAddress)
	}
	if !slices.Contains(iface.TaggedVLANs, 4) {
		t.Errorf("expected interface ge-0/0/0 being member of vlan 4, got %v", iface.TaggedVLANs)
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
					<description>example interface 2</description>

					<unit>
						<name>0</name>
						<family>
							<ethernet-switching>
								<port-mode>trunk</port-mode>
								<vlan>
									<members>1</members><members>2</members><members>3</members>
								</vlan>
								<native-vlan-id>1337</native-vlan-id>
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
	tmpl, err := template.New("").Parse(juniper.EditPortConfigurationTemplate)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	interfaceDescription := "example interface 2"
	untaggedVlan := int32(1337)
	if err = tmpl.Execute(&tpl, &models.UpdateInterface{
		Name:         "eth0",
		Description:  &interfaceDescription,
		UntaggedVLAN: &untaggedVlan,
		TaggedVLANs:  []int32{1, 2, 3},
	}); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	isEqual, err := utils.CompareXMLIgnoreWhitespace(tpl.String(), EditPortConfigurationExpected)
	if err != nil {
		t.Errorf("error while comparing XML: %s", err.Error())
	}
	if !isEqual {
		diff := cmp.Diff(EditPortConfigurationExpected, tpl.String())
		t.Errorf("expected:\n%s\ngot:\n%s\n\ndiff:\n%s", EditPortConfigurationExpected, tpl.String(), diff)
	}
}
