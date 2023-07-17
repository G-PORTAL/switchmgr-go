package juniper_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper"
	"testing"
)

func TestGetHardwareInfo(t *testing.T) {
	driver := juniper.NewMockDriver()
	system, err := driver.GetHardwareInfo()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if system.Model != "EX3200-48T" {
		t.Errorf("expected model EX3200-48T, got %s", system.Model)
	}
	if system.Serial != "xxx" {
		t.Errorf("expected serial xxx, got %s", system.Serial)
	}
	if system.Hostname != "switch-host-name" {
		t.Errorf("expected hostname switch-host-name, got %s", system.Hostname)
	}
	if system.FirmwareVersion != "Junos 14.1X53" {
		t.Errorf("expected firmware version Junos 14.1X53, got %s", system.FirmwareVersion)
	}
	if system.Vendor != "Juniper" {
		t.Errorf("expected vendor Juniper, got %s", system.Vendor)
	}
}
