package juniper

import (
	"encoding/xml"
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/openshift-telco/go-netconf-client/netconf/message"
	"strings"
)

func (j *Juniper) GetHardwareInfo() (*models.HardwareInfo, error) {
	reply, err := j.session.SyncRPC(message.NewRPC("<get-system-information/>"), 10)
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
