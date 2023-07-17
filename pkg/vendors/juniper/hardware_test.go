package juniper_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper"
	"testing"
)

var hardwareResponse = `<rpc-reply xmlns:junos="http://xml.juniper.net/junos/14.1X53/junos" message-id="d9aa76b3-3580-42ed-b0e6-94190c15ae3a" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
<system-information>
<hardware-model>ex3200-48t</hardware-model>
<os-name>junos-ex</os-name>
<os-version>14.1X53-D40.8</os-version>
<serial-number>BK0210468776</serial-number>
<host-name>switch15_Frankfurt</host-name>
</system-information>
</rpc-reply>
`

func TestParseHardware(t *testing.T) {
	system, err := juniper.ParseHardwareInfo([]byte(hardwareResponse))
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if system.Model != "ex3200-48t" {
		t.Errorf("expected model ex3200-48t, got %s", system.Model)
	}
	if system.Serial != "BK0210468776" {
		t.Errorf("expected serial BK0210468776, got %s", system.Serial)
	}
	if system.Hostname != "switch15_Frankfurt" {
		t.Errorf("expected hostname switch15_Frankfurt, got %s", system.Hostname)
	}
	if system.FirmwareVersion != "14.1X53-D40.8" {
		t.Errorf("expected firmware version 14.1X53-D40.8, got %s", system.FirmwareVersion)
	}
}
