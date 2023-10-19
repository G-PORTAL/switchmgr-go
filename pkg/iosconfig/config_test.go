package iosconfig_test

import (
	"github.com/g-portal/switchmgr-go/pkg/iosconfig"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s3/utils"
	"testing"
)

func TestConfig_Interfaces(t *testing.T) {
	cfg := iosconfig.Parse(utils.ReadTestData("show running-config", nil))

	if len(cfg) != 46 {
		t.Errorf("Config length is not 46, it is %d", len(cfg))
	}

	if len(cfg.Interfaces()) != 30 {
		t.Errorf("Interface count is not 30, it is %d", len(cfg.Interfaces()))
	}

	nicConfig, err := cfg.Interface("GigaEthernet0/12")
	if err != nil {
		t.Error(err)
	}

	nicConfigMode := nicConfig.GetStringValue("switchport mode", "")
	if nicConfigMode != "trunk" {
		t.Errorf("switchport mode is not trunk, it is %q", nicConfigMode)
	}

	nicConfigPVID := nicConfig.GetInt32Value("switchport pvid", 0)
	if nicConfigPVID != 6 {
		t.Errorf("switchport pvid is not 6, it is %d", nicConfigPVID)
	}
}

func TestParse(t *testing.T) {
	cfg := iosconfig.Parse("")
	if len(cfg) > 0 {
		t.Errorf("Config length is not 0, it is %d", len(cfg))
	}

	cfg = iosconfig.Parse("!!\n!\n\n\n!")
	if len(cfg) > 0 {
		t.Errorf("Config length is not 0, it is %d", len(cfg))
	}

	cfg = iosconfig.Parse("!\na 1\n!\nb 2\n!\nc 3\n")
	if len(cfg) != 3 {
		t.Errorf("Config length is not 3 it is %d", len(cfg))
	}

	cfg = iosconfig.Parse("\n!\ninterface eth0\n a 1\nb 5\nc 3\n!\nb 2\n!\nc 3\n!")
	if len(cfg) != 3 {
		t.Errorf("Config length is not 3, it is %d", len(cfg))
	}
}
