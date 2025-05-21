package fsos_n5

import (
	"errors"
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"regexp"
	"strings"
)

var serialRegex = regexp.MustCompile(`System\sserial\snumber\s+:\s([0-9A-Z]+)`)
var seria2lRegex = regexp.MustCompile(`Serial\snumber\s+:\s([0-9A-Z]+)`)
var modelRegex = regexp.MustCompile(`System\sdescription\s+:\s.+\(([0-9A-Z-]+)\).+`)
var versionRegex = regexp.MustCompile(`System\ssoftware\sversion\s+:\s.+_FSOS\s(.+)`)
var hostnameRegex = regexp.MustCompile(`hostname\s(.+)`)

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
		matches = seria2lRegex.FindStringSubmatch(output)
		if len(matches) != 2 {
			return nil, errors.New("could not parse serial")
		}
	}
	hwInfo.Serial = strings.TrimSpace(matches[1])

	// model
	matches = modelRegex.FindStringSubmatch(output)
	if len(matches) != 2 {
		return nil, errors.New("could not parse model")
	}
	hwInfo.Model = strings.TrimSpace(matches[1])

	// firmware version
	matches = versionRegex.FindStringSubmatch(output)
	if len(matches) != 2 {
		return nil, errors.New("could not parse firmware version")
	}
	hwInfo.FirmwareVersion = fmt.Sprintf("FSOS %s", strings.TrimSpace(matches[1]))

	// hostname
	matches = hostnameRegex.FindStringSubmatch(output)
	if len(matches) != 2 {
		return nil, errors.New("could not parse hostname")
	}
	hwInfo.Hostname = strings.TrimSpace(matches[1])

	return hwInfo, nil
}
