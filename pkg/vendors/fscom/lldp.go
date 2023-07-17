package fscom

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fscom/utils"
	"regexp"
	"strings"
)

var lldpLineRgx = regexp.MustCompile(`([a-zA-Z0-9\_]+)\s+([a-zA-Z0-9\_]+\/[0-9]+)\s`)

func (fs *FSCom) ListLLDPNeighbors() ([]models.LLDPNeighbor, error) {
	output, err := fs.sendCommands("show lldp neighbors")
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
