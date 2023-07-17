package fscom_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fscom"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fscom/utils"
	"testing"
)

func TestParseHardwareInfo(t *testing.T) {
	hwInfo, err := fscom.ParseHardwareInfo(utils.ReadTestData("show version"))
	if err != nil {
		t.Fatalf("Error parsing hardware info: %s", err.Error())
	}

	if hwInfo == nil {
		t.Fatalf("Hardware info is nil")
	}

	if hwInfo.Serial != "CG22041XXXXXXXXXXX" {
		t.Fatalf("Serial number is wrong: %s", hwInfo.Serial)
	}

	if hwInfo.Model != "S3900-24T4S-R" {
		t.Fatalf("Model is wrong: %s", hwInfo.Model)
	}

	if hwInfo.Vendor != "Fiberstore" {
		t.Fatalf("Vendor is wrong: %s", hwInfo.Vendor)
	}

	if hwInfo.FirmwareVersion != "FSOS 2.2.0E" {
		t.Fatalf("Firmware version is wrong: %s", hwInfo.FirmwareVersion)
	}

	if hwInfo.Hostname != "switch" {
		t.Fatalf("Hostname is wrong: %s", hwInfo.Hostname)
	}

}
