package fsos_n5

import (
	"errors"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"regexp"
)

var serialRegex = regexp.MustCompile(`System\sserial\snumber\s+:\s([0-9A-Z]+)\r\n`)
var modelRegex = regexp.MustCompile(`System\sdescription\s+:\s.+\(([0-9A-Z-]+)\).+\r\n`)
var versionRegex = regexp.MustCompile(`System\ssoftware\sversion\s+:\s(.+)\r\n`)
var hostnameRegex = regexp.MustCompile(`hostname\s(.+)\r\n`)

func (fs *FSComN5) GetHardwareInfo() (*models.HardwareInfo, error) {
	output, err := fs.SendCommands("show version", "show running-config | include hostname")
	if err != nil {
		return nil, err
	}

	return ParseHardwareInfo(output)
}

func ParseHardwareInfo(output string) (*models.HardwareInfo, error) {
	hwInfo := &models.HardwareInfo{
		Vendor: "Fiberstore",
	}

	// serial
	matches := serialRegex.FindStringSubmatch(output)
	if len(matches) != 2 {
		return nil, errors.New("could not parse serial")
	}
	hwInfo.Serial = matches[1]

	// model
	matches = modelRegex.FindStringSubmatch(output)
	if len(matches) != 2 {
		return nil, errors.New("could not parse model")
	}
	hwInfo.Model = matches[1]

	// firmware version
	matches = versionRegex.FindStringSubmatch(output)
	if len(matches) != 2 {
		return nil, errors.New("could not parse firmware version")
	}
	hwInfo.FirmwareVersion = matches[1]

	// hostname
	matches = hostnameRegex.FindStringSubmatch(output)
	if len(matches) != 2 {
		return nil, errors.New("could not parse hostname")
	}
	hwInfo.Hostname = matches[1]

	return hwInfo, nil
}
