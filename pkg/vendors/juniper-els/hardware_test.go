package juniper_els_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper-els"
	"testing"
)

func TestGetHardwareInfo(t *testing.T) {
	driver := juniper_els.NewMockDriver()
	system, err := driver.GetHardwareInfo()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if system.Model != "EX3400-48T" {
		t.Errorf("expected model EX3400-48T, got %s", system.Model)
	}
	if system.Serial != "XXX" {
		t.Errorf("expected serial XXX, got %s", system.Serial)
	}
	if system.Hostname != "switch-host-name" {
		t.Errorf("expected hostname switch-host-name, got %s", system.Hostname)
	}
	if system.FirmwareVersion != "Junos 18.2R3-S8.5" {
		t.Errorf("expected firmware version Junos 18.2R3-S8.5, got %s", system.FirmwareVersion)
	}
	if system.Vendor != "Juniper" {
		t.Errorf("expected vendor Juniper, got %s", system.Vendor)
	}
}
