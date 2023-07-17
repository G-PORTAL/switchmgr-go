package juniper_els

import (
	"encoding/xml"
	"github.com/Juniper/go-netconf/netconf"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"strings"
)

func (j *JuniperELS) ListArpTable() ([]models.ArpEntry, error) {
	reply, err := j.session.Exec(netconf.RawMethod("<get-ethernet-switching-table-information/>"))
	if err != nil {
		return nil, err
	}
	var table junosArpTable
	if err := xml.Unmarshal([]byte(reply.RawReply), &table); err != nil {
		return nil, err
	}

	entries := make(map[string][]models.MacAddress)
	for _, entry := range table.Entries {
		entryName := strings.TrimSuffix(strings.TrimSpace(entry.Interface), ".0")
		if strings.HasPrefix(entryName, "esi.") {
			continue
		}

		if _, ok := entries[entryName]; !ok {
			entries[entryName] = make([]models.MacAddress, 0)
		}

		entries[entryName] = append(entries[entryName], entry.MacAddress)
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
