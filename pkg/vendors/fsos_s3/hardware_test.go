package fsos_s3_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s3"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s3/utils"
	"testing"
)

func TestParseHardwareInfo(t *testing.T) {
	hwInfo, err := fsos_s3.ParseHardwareInfo(utils.ReadTestData("show version", nil))
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

func TestParseHardwareInfo2(t *testing.T) {
	version := 2

	hwInfo, err := fsos_s3.ParseHardwareInfo(utils.ReadTestData("show version", &version))
	if err != nil {
		t.Fatalf("Error parsing hardware info: %s", err.Error())
	}

	if hwInfo == nil {
		t.Fatalf("Hardware info is nil")
	}

	if hwInfo.Serial != "CG22042XXXXXXXXXXX" {
		t.Fatalf("Serial number is wrong: %s", hwInfo.Serial)
	}

	if hwInfo.Model != "S3900-24T4S-R" {
		t.Fatalf("Model is wrong: %s", hwInfo.Model)
	}

	if hwInfo.Vendor != "Fiberstore" {
		t.Fatalf("Vendor is wrong: %s", hwInfo.Vendor)
	}

	if hwInfo.FirmwareVersion != "FSOS 2.2.0F" {
		t.Fatalf("Firmware version is wrong: %s", hwInfo.FirmwareVersion)
	}

	if hwInfo.Hostname != "switch" {
		t.Fatalf("Hostname is wrong: %s", hwInfo.Hostname)
	}

}
