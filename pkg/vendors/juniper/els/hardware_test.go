package els_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper/els"
	"testing"
)

var hardwareResponse = `<rpc-reply xmlns:junos="http://xml.juniper.net/junos/18.2R3/junos" message-id="de5bef26-4415-46e3-8f9c-dcbd93d29970" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
<system-information>
<hardware-model>ex3400-48t</hardware-model>
<os-name>junos</os-name>
<os-version>18.2R3-S8.5</os-version>
<serial-number>NX0219440002</serial-number>
<host-name>switch64_Frankfurt</host-name>
</system-information>
</rpc-reply>
`

func TestParseHardware(t *testing.T) {
	system, err := els.ParseHardwareInfo([]byte(hardwareResponse))
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if system.Model != "ex3400-48t" {
		t.Errorf("expected model ex3400-48t, got %s", system.Model)
	}
	if system.Serial != "NX0219440002" {
		t.Errorf("expected serial NX0219440002, got %s", system.Serial)
	}
	if system.Hostname != "switch64_Frankfurt" {
		t.Errorf("expected hostname switch64_Frankfurt, got %s", system.Hostname)
	}
	if system.FirmwareVersion != "18.2R3-S8.5" {
		t.Errorf("expected firmware version 18.2R3-S8.5, got %s", system.FirmwareVersion)
	}
}
