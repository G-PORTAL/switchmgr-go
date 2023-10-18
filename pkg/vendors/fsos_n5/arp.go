package fsos_n5

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s3/utils"
	"golang.org/x/exp/slices"
	"regexp"
	"strings"
)

// Example
// 1        649d.99d2.502d       MLAG     AggregatePort 10                       32d 21:31:31
var arpLineRgx = regexp.MustCompile(`^([0-9]+)\s+([0-9a-f]{4}\.[0-9a-f]{4}\.[0-9a-f]{4})\s+([A-Z]+)\s+([a-zA-Z]+\s[0-9]+)\s+([0-9a-z]+)\s([0-9:]+)$`)

func (fs *FSComN5) ListArpTable() ([]models.ArpEntry, error) {
	output, err := fs.SendCommands("show mac-address-table")
	if err != nil {
		return nil, err
	}

	return ParseArpTable(output)
}

func ParseArpTable(output string) ([]models.ArpEntry, error) {
	var table []models.ArpEntry

	portsWithMac := map[string][]models.MacAddress{}

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if !arpLineRgx.MatchString(line) {
			continue
		}

		matches := arpLineRgx.FindStringSubmatch(line)
		if len(matches) != 7 {
			continue
		}

		matches[4] = utils.ConvertInterface(matches[4])
		if _, ok := portsWithMac[matches[4]]; !ok {
			portsWithMac[matches[4]] = []models.MacAddress{}
		}

		mac := models.MacAddress(strings.TrimSpace(matches[2]))
		if !mac.Valid() {
			continue
		}

		// reformat mac address
		mac = models.MacAddress(mac.String())

		// avoid duplicate entries
		if slices.Contains(portsWithMac[matches[4]], mac) {
			continue
		}

		portsWithMac[matches[4]] = append(portsWithMac[matches[4]], mac)
	}

	for portName, macs := range portsWithMac {
		table = append(table, models.ArpEntry{
			SwitchPort:   portName,
			MacAddresses: macs,
		})
	}

	return table, nil
}
