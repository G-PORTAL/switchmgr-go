package els

import (
	"encoding/xml"
	"github.com/Juniper/go-netconf/netconf"
	"github.com/g-portal/switchmgr-go/pkg/models"
)

func ParseHardwareInfo(data []byte) (*models.HardwareInfo, error) {
	var system *junosSystemInformation
	if err := xml.Unmarshal(data, &system); err != nil {
		return nil, err
	}

	hwInfo := models.HardwareInfo{
		Hostname:        system.HostName,
		Vendor:          "Juniper",
		Model:           system.HardwareModel,
		FirmwareVersion: system.OsVersion,
		Serial:          system.SerialNumber,
	}
	return &hwInfo, nil
}

func (j *JuniperELS) GetHardwareInfo() (*models.HardwareInfo, error) {
	reply, err := j.session.Exec(netconf.RawMethod("<get-system-information/>"))
	if err != nil {
		return nil, err
	}
	return ParseHardwareInfo([]byte(reply.RawReply))
}
