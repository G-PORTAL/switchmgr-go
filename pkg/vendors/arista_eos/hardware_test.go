package arista_eos_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/arista_eos"
	"github.com/g-portal/switchmgr-go/pkg/vendors/arista_eos/utils"
	"testing"
)

func TestParseHardwareInfo(t *testing.T) {
	var version arista_eos.AristaVersionResponse
	err := utils.ReadTestData("show version", &version)
	if err != nil {
		t.Error(err)
	}

	var hostname arista_eos.AristaHostnameResponse
	err = utils.ReadTestData("show hostname", &hostname)
	if err != nil {
		t.Error(err)
	}

	hwInfo, err := arista_eos.ParseHardwareInfo(version, hostname)
	if err != nil {
		t.Fatalf("Error parsing hardware info: %s", err.Error())
	}

	if hwInfo == nil {
		t.Fatalf("Hardware info is nil")
	}

	if hwInfo.Serial != "WTW22400000" {
		t.Fatalf("Serial number is wrong: %s", hwInfo.Serial)
	}

	if hwInfo.Model != "CCS-720DT-48S-2F" {
		t.Fatalf("Model is wrong: %s", hwInfo.Model)
	}

	if hwInfo.Vendor != "Arista" {
		t.Fatalf("Vendor is wrong: %s", hwInfo.Vendor)
	}

	if hwInfo.FirmwareVersion != "EOS 4.31.3M" {
		t.Fatalf("Firmware version is wrong: %s", hwInfo.FirmwareVersion)
	}

	if hwInfo.Hostname != "switch1" {
		t.Fatalf("Hostname is wrong: %s", hwInfo.Hostname)
	}
}
