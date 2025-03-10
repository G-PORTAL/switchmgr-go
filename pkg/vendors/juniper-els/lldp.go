package juniper_els

import (
	"encoding/xml"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/openshift-telco/go-netconf-client/netconf/message"
	"strings"
)

func (j *JuniperELS) ListLLDPNeighbors() ([]models.LLDPNeighbor, error) {
	output, err := j.session.SyncRPC(message.NewRPC("<get-lldp-neighbors-information/>"), 10)
	if err != nil {
		return nil, err
	}

	neighbors, err := ParseLLDPNeighbors([]byte(output.RawReply))
	if err != nil {
		return nil, err
	}

	lldpNeighbors := make([]models.LLDPNeighbor, 0)
	for _, neighbor := range neighbors.Neighbors {
		localInterface := strings.TrimSuffix(strings.TrimSpace(neighbor.LocalPortID), ".0")

		lldpNeighbors = append(lldpNeighbors, models.LLDPNeighbor{
			LocalInterface: localInterface,
			RemoteHostname: strings.TrimSpace(neighbor.RemoteHostname),
		})
	}

	return lldpNeighbors, nil
}

func ParseLLDPNeighbors(data []byte) (*junosLLDPNeighbors, error) {
	var neighbors junosLLDPNeighbors
	if err := xml.Unmarshal(data, &neighbors); err != nil {
		return nil, err
	}

	return &neighbors, nil
}
