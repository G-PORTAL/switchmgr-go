package fsos_s5

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
	"regexp"
	"strings"
)

var lldpLineRgx = regexp.MustCompile(`(?s)Local\sPort\s+:\s+([A-Za-z0-9-]+)\r\n.+System\sName\s+:\s+(.+)\r\n`)
var lldpGroupRgx = regexp.MustCompile(`(-+)\r\n`)

func (fs *FSComS5) ListLLDPNeighbors() ([]models.LLDPNeighbor, error) {
	output, err := fs.SendCommands("show lldp neighbor brief")
	if err != nil {
		return nil, err
	}

	return ParseLLDPNeighbors(output)
}

func ParseLLDPNeighbors(output string) ([]models.LLDPNeighbor, error) {
	groups := lldpGroupRgx.Split(output, -1)

	var table []models.LLDPNeighbor
	for _, group := range groups {
		lldpEntries := lldpLineRgx.FindAllStringSubmatch(group, -1)
		for _, entry := range lldpEntries {
			if len(entry) != 3 {
				continue
			}

			table = append(table, models.LLDPNeighbor{
				LocalInterface: strings.TrimSpace(entry[1]),
				RemoteHostname: strings.TrimSpace(entry[2]),
			})
		}
	}

	return table, nil
}
