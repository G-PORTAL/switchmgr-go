package fsos_s5_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s5"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s5/utils"
	"testing"
)

func TestParseHardwareInfo(t *testing.T) {
	hwInfo, err := fsos_s5.ParseHardwareInfo(utils.ReadTestData("show version", nil))
	if err != nil {
		t.Fatalf("Error parsing hardware info: %s", err.Error())
	}

	if hwInfo == nil {
		t.Fatalf("Hardware info is nil")
	}

	if hwInfo.Serial != "CG2112162735N00000" {
		t.Fatalf("Serial number is wrong: %s", hwInfo.Serial)
	}

	if hwInfo.Model != "S5800-48T4S" {
		t.Fatalf("Model is wrong: %s", hwInfo.Model)
	}

	if hwInfo.Vendor != "Fiberstore" {
		t.Fatalf("Vendor is wrong: %s", hwInfo.Vendor)
	}

	if hwInfo.FirmwareVersion != "FSOS 7.4.1.r1" {
		t.Fatalf("Firmware version is wrong: %s", hwInfo.FirmwareVersion)
	}

	if hwInfo.Hostname != "wup-as-test123" {
		t.Fatalf("Hostname is wrong: %s", hwInfo.Hostname)
	}

}

func TestParseHardwareInfo2(t *testing.T) {
	hwInfo, err := fsos_s5.ParseHardwareInfo(utils.ReadTestData("show version 2", nil))
	if err != nil {
		t.Fatalf("Error parsing hardware info: %s", err.Error())
	}

	if hwInfo == nil {
		t.Fatalf("Hardware info is nil")
	}

	if hwInfo.Serial != "CG2206232964N00000" {
		t.Fatalf("Serial number is wrong: %s", hwInfo.Serial)
	}

	if hwInfo.Model != "S5800-48T4S" {
		t.Fatalf("Model is wrong: %s", hwInfo.Model)
	}

	if hwInfo.Vendor != "Fiberstore" {
		t.Fatalf("Vendor is wrong: %s", hwInfo.Vendor)
	}

	if hwInfo.FirmwareVersion != "FSOS 7.4.3.r5" {
		t.Fatalf("Firmware version is wrong: %s", hwInfo.FirmwareVersion)
	}

	if hwInfo.Hostname != "stl-as-test123" {
		t.Fatalf("Hostname is wrong: %s", hwInfo.Hostname)
	}

}
