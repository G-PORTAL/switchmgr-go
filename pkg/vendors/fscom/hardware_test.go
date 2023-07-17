package fscom_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fscom"
	"testing"
)

var versionOutput = `switch05_tw#show version
Fiberstore Co., Limited Internetwork Operating System Software
S3900-24T4S-R Series Software, Version 2.2.0E Build 99013, RELEASE SOFTWARE
Copyright (c) 2021 by FS.COM All Rights Reserved
Compiled: 2022-6-10 13:12:54 by SYS, Image text-base: 0x80010000
ROM: System Bootstrap, Version 0.1.3,hardware version:A
Serial num:CG22041XXXXXN00116, ID num:20057011962
System image file is "Switch.bin"
FS S3900-24T4S-R RISC
262144K bytes of memory,16384K bytes of flash
Base ethernet MAC Address: 64:9d:99:c5:72:fb
PCB version:B 
snmp info:
  vend_ID:52642   product_ID:446   system_ID:1.3.6.1.4.1.52642.1.446.0
switch05_tw uptime is 139:03:46:01, The current time: 2000-5-19 5:39:21
 Reboot history information:
  No. 1: System is rebooted by power-on
`

func TestParseHardwareInfo(t *testing.T) {
	hwInfo, err := fscom.ParseHardwareInfo(versionOutput)
	if err != nil {
		t.Fatalf("Error parsing hardware info: %s", err.Error())
	}

	if hwInfo == nil {
		t.Fatalf("Hardware info is nil")
	}

	if hwInfo.Serial != "CG22041XXXXXN00116" {
		t.Fatalf("Serial number is wrong: %s", hwInfo.Serial)
	}

	if hwInfo.Model != "S3900-24T4S-R" {
		t.Fatalf("Model is wrong: %s", hwInfo.Model)
	}

	if hwInfo.Vendor != "Fiberstore" {
		t.Fatalf("Vendor is wrong: %s", hwInfo.Vendor)
	}

	if hwInfo.FirmwareVersion != "2.2.0E" {
		t.Fatalf("Firmware versio is wrong: %s", hwInfo.FirmwareVersion)
	}

	if hwInfo.Hostname != "switch05_tw" {
		t.Fatalf("Hostname is wrong: %s", hwInfo.Hostname)
	}

}
