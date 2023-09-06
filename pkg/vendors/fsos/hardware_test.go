package fsos_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos/utils"
	"testing"
)

func TestParseHardwareInfo(t *testing.T) {
	hwInfo, err := fsos.ParseHardwareInfo(utils.ReadTestData("show version", nil))
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
