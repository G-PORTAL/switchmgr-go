package fsos_n5_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_n5"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_n5/utils"
	"testing"
)

func TestParseHardwareInfo(t *testing.T) {
	hwInfo, err := fsos_n5.ParseHardwareInfo(utils.ReadTestData("show version", nil))
	if err != nil {
		t.Fatalf("Error parsing hardware info: %s", err.Error())
	}

	if hwInfo == nil {
		t.Fatalf("Hardware info is nil")
	}

	if hwInfo.Serial != "G1RL71S00XXXX" {
		t.Fatalf("Serial number is wrong: %s", hwInfo.Serial)
	}

	if hwInfo.Model != "N5860-48SC" {
		t.Fatalf("Model is wrong: %s", hwInfo.Model)
	}

	if hwInfo.Vendor != "Fiberstore" {
		t.Fatalf("Vendor is wrong: %s", hwInfo.Vendor)
	}

	if hwInfo.FirmwareVersion != "FSOS 11.0(5)B9P66S2" {
		t.Fatalf("Firmware version is wrong: %s", hwInfo.FirmwareVersion)
	}

	if hwInfo.Hostname != "openshift-core2" {
		t.Fatalf("Hostname is wrong: %s", hwInfo.Hostname)
	}

}
