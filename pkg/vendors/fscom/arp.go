package fscom

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fscom/utils"
	"regexp"
	"strings"
)

var arpLineRgx = regexp.MustCompile(`^(\d+)\s+([0-9a-fA-F]{4}\.[0-9a-fA-F]{4}\.[0-9a-fA-F]{4})\s+([A-Z]+)\s+(.*)$`)

func (fs *FSCom) ListArpTable() ([]models.ArpEntry, error) {
	output, err := fs.sendCommands("show mac address-table")
	if err != nil {
		return nil, err
	}

	return ParseArpTable(output)
}

func ParseArpTable(output string) ([]models.ArpEntry, error) {
	var table []models.ArpEntry

	portsWithMac := map[string][]models.MacAddress{}

	for _, line := range strings.Split(output, "\n") {
		if !arpLineRgx.MatchString(line) {
			continue
		}

		matches := arpLineRgx.FindStringSubmatch(line)
		if len(matches) != 5 {
			continue
		}

		matches[4] = utils.ConvertInterface(matches[4])

		if _, ok := portsWithMac[matches[4]]; !ok {
			portsWithMac[matches[4]] = []models.MacAddress{}
		}

		mac := models.MacAddress(matches[2])
		if !mac.Valid() {
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
