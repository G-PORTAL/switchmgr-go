package fsos_n5

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_n5/utils"
	"regexp"
	"strings"
)

var lldpLineRgx = regexp.MustCompile(`^([a-zA-Z0-9_-]+)\s+([A-Za-z0-9]+/[0-9]+).*$`)

func (fs *FSComN5) ListLLDPNeighbors() ([]models.LLDPNeighbor, error) {
	output, err := fs.SendCommands("show lldp neighbors")
	if err != nil {
		return nil, err
	}

	return ParseLLDPNeighbors(output)
}

func ParseLLDPNeighbors(output string) ([]models.LLDPNeighbor, error) {
	var table []models.LLDPNeighbor

	for _, line := range strings.Split(output, "\n") {
		if !lldpLineRgx.MatchString(line) {
			continue
		}

		matches := lldpLineRgx.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}

		table = append(table, models.LLDPNeighbor{
			LocalInterface: utils.ConvertInterface(matches[2]),
			RemoteHostname: matches[1],
		})
	}

	return table, nil
}
