package arista_eos

import (
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/models"
)

type AristaVersionResponse struct {
	MfgName            string  `json:"mfgName"`
	ModelName          string  `json:"modelName"`
	HardwareRevision   string  `json:"hardwareRevision"`
	SerialNumber       string  `json:"serialNumber"`
	SystemMacAddress   string  `json:"systemMacAddress"`
	HwMacAddress       string  `json:"hwMacAddress"`
	ConfigMacAddress   string  `json:"configMacAddress"`
	Version            string  `json:"version"`
	Architecture       string  `json:"architecture"`
	InternalVersion    string  `json:"internalVersion"`
	InternalBuildID    string  `json:"internalBuildId"`
	ImageFormatVersion string  `json:"imageFormatVersion"`
	ImageOptimization  string  `json:"imageOptimization"`
	BootupTimestamp    float64 `json:"bootupTimestamp"`
	Uptime             float64 `json:"uptime"`
	MemTotal           int     `json:"memTotal"`
	MemFree            int     `json:"memFree"`
	IsIntlVersion      bool    `json:"isIntlVersion"`
}

type AristaHostnameResponse struct {
	Hostname string `json:"hostname"`
	Fqdn     string `json:"fqdn"`
}

func (arista *AristaEOS) GetHardwareInfo() (*models.HardwareInfo, error) {
	var versionResponse AristaVersionResponse
	var hostnameResponse AristaHostnameResponse

	err := arista.GetJsonResponse(&versionResponse, "show version")
	if err != nil {
		return nil, err
	}

	err = arista.GetJsonResponse(&hostnameResponse, "show hostname")
	if err != nil {
		return nil, err
	}

	return ParseHardwareInfo(versionResponse, hostnameResponse)
}

func ParseHardwareInfo(version AristaVersionResponse, hostname AristaHostnameResponse) (*models.HardwareInfo, error) {
	hwInfo := &models.HardwareInfo{
		Hostname:        hostname.Hostname,
		Vendor:          version.MfgName,
		Model:           version.ModelName,
		FirmwareVersion: fmt.Sprintf("EOS %s", version.Version),
		Serial:          version.SerialNumber,
	}

	return hwInfo, nil
}
