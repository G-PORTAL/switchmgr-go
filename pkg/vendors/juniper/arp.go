package juniper

import (
	"encoding/xml"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/openshift-telco/go-netconf-client/netconf/message"
	"strings"
)

func ParseArpTable(data []byte) ([]models.ArpEntry, error) {
	var table junosArpTable
	if err := xml.Unmarshal(data, &table); err != nil {
		return nil, err
	}

	entries := make(map[string][]models.MacAddress)

	for _, entry := range table.Entries {
		entryName := strings.TrimSuffix(strings.TrimSpace(entry.Interfaces[0].Interface), ".0")
		if strings.HasPrefix(entryName, "esi.") {
			continue
		}

		if len(entry.Interfaces) == 0 || !isValidInterface(entryName) {
			continue
		}
		if _, ok := entries[entryName]; !ok {
			entries[entryName] = make([]models.MacAddress, 0)
		}

		mac := entry.MacAddress
		if !mac.Valid() {
			continue
		}

		entries[entryName] = append(entries[entryName], mac)
	}

	result := make([]models.ArpEntry, 0)
	for port, macs := range entries {
		result = append(result, models.ArpEntry{
			SwitchPort:   port,
			MacAddresses: macs,
		})
	}

	return result, nil
}

func (j *Juniper) ListArpTable() ([]models.ArpEntry, error) {
	reply, err := j.session.SyncRPC(message.NewRPC("<get-ethernet-switching-table-information/>"), Timeout)
	if err != nil {
		return nil, err
	}
	return ParseArpTable([]byte(reply.RawReply))

}
