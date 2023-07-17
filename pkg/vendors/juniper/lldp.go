package juniper

import (
	"encoding/xml"
	"github.com/Juniper/go-netconf/netconf"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"strings"
)

func (j *Juniper) ListLLDPNeighbors() ([]models.LLDPNeighbor, error) {
	reply, err := j.session.Exec(netconf.RawMethod("<get-lldp-neighbors-information/>"))
	if err != nil {
		return nil, err
	}
	var interfaces *junosLLDPInformation
	if err := xml.Unmarshal([]byte(reply.RawReply), &interfaces); err != nil {
		return nil, err
	}

	// TODO: add mtu, speed, management, mode
	entries := make([]models.LLDPNeighbor, 0)
	for _, entry := range interfaces.Entries {
		localInterface1 := strings.TrimSuffix(strings.TrimSpace(entry.LocalInterface), ".0")
		localInterface2 := strings.TrimSuffix(strings.TrimSpace(entry.LocalPortID), ".0")

		if localInterface1 != "" {
			entries = append(entries, models.LLDPNeighbor{
				LocalInterface: localInterface1,
				RemoteHostname: strings.TrimSpace(entry.RemoteHostname),
			})
		} else if localInterface2 != "" {
			entries = append(entries, models.LLDPNeighbor{
				LocalInterface: localInterface2,
				RemoteHostname: strings.TrimSpace(entry.RemoteHostname),
			})
		}
	}

	return entries, nil
}
