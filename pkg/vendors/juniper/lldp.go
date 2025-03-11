package juniper

import (
	"encoding/xml"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/openshift-telco/go-netconf-client/netconf/message"
	"strings"
)

func (j *Juniper) ListLLDPNeighbors() ([]models.LLDPNeighbor, error) {
	reply, err := j.session.SyncRPC(message.NewRPC("<get-lldp-neighbors-information/>"), Timeout)
	if err != nil {
		return nil, err
	}
	var interfaces *junosLLDPInformation
	if err := xml.Unmarshal([]byte(reply.RawReply), &interfaces); err != nil {
		return nil, err
	}

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
