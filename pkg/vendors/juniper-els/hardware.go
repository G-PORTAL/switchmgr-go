package juniper_els

import (
	"encoding/xml"
	"fmt"
	"github.com/Juniper/go-netconf/netconf"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"strings"
)

func (j *JuniperELS) GetHardwareInfo() (*models.HardwareInfo, error) {
	reply, err := j.session.Exec(netconf.RawMethod("<get-system-information/>"))
	if err != nil {
		return nil, err
	}
	var system *junosSystemInformation
	if err := xml.Unmarshal([]byte(reply.RawReply), &system); err != nil {
		return nil, err
	}

	hwInfo := models.HardwareInfo{
		Hostname:        system.HostName,
		Vendor:          "Juniper",
		Model:           strings.ToUpper(system.HardwareModel),
		FirmwareVersion: fmt.Sprintf("Junos %s", system.OsVersion),
		Serial:          system.SerialNumber,
	}
	return &hwInfo, nil
}
