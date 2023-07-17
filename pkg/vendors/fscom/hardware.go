package fscom

import (
	"errors"
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"regexp"
)

var serialRegex = regexp.MustCompile(`Serial num:([0-9A-Z]+), ID num:`)
var modelRegex = regexp.MustCompile(`([0-9A-Z-]+)\sSeries Software, Version`)
var versionRegex = regexp.MustCompile(`Series\sSoftware,\sVersion\s([0-9A-Z.]+)`)
var hostnameRegex = regexp.MustCompile(`(.+)\suptime\sis\s`)

func (fs *FSCom) GetHardwareInfo() (*models.HardwareInfo, error) {
	output, err := fs.sendCommands("show version")
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
	hwInfo.FirmwareVersion = fmt.Sprintf("FSOS %s", matches[1])

	// hostname
	matches = hostnameRegex.FindStringSubmatch(output)
	if len(matches) != 2 {
		return nil, errors.New("could not parse hostname")
	}
	hwInfo.Hostname = matches[1]

	return hwInfo, nil
}
