package fsos_s5

import (
	"errors"
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"regexp"
)

var serialRegex = regexp.MustCompile(`System\sserial\snumber\sis\s([0-9A-Z]+)\r\n`)
var seriesRegex = regexp.MustCompile(`FSOS\sSoftware,\s([0-9A-Z]+),\sVersion`)
var modelRegex = regexp.MustCompile(`Hardware\sType\sis\s([0-9A-Z-]+)\r\n`)
var versionRegex = regexp.MustCompile(`Current\sWeb\sVersion\sis\s([0-9A-Za-z.]+)`)
var version2Regex = regexp.MustCompile(`,\sVersion\s([0-9A-Za-z.]+)`)
var hostnameRegex = regexp.MustCompile(`(.+)\suptime\sis\s`)

func (fs *FSComS5) GetHardwareInfo() (*models.HardwareInfo, error) {
	output, err := fs.SendCommands("show version")
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
	matches = seriesRegex.FindStringSubmatch(output)
	if len(matches) != 2 {
		return nil, errors.New("could not parse model series")
	}
	hwInfo.Model = matches[1]
	matches = modelRegex.FindStringSubmatch(output)
	if len(matches) != 2 {
		return nil, errors.New("could not parse model")
	}
	hwInfo.Model = fmt.Sprintf("%s-%s", hwInfo.Model, matches[1])

	// firmware version
	matches = versionRegex.FindStringSubmatch(output)
	if len(matches) != 2 {
		matches = version2Regex.FindStringSubmatch(output)
		if len(matches) != 2 {
			return nil, errors.New("could not parse firmware version")
		}
	}
	hwInfo.FirmwareVersion = fmt.Sprintf("FSComS5 %s", matches[1])

	// hostname
	matches = hostnameRegex.FindStringSubmatch(output)
	if len(matches) != 2 {
		return nil, errors.New("could not parse hostname")
	}
	hwInfo.Hostname = matches[1]

	return hwInfo, nil
}
