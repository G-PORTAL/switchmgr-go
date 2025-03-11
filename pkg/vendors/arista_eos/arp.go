package arista_eos

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
	"sort"
)

type AristaMacAddressTableResponse struct {
	UnicastTable struct {
		TableEntries []struct {
			VlanID     int     `json:"vlanId"`
			MacAddress string  `json:"macAddress"`
			EntryType  string  `json:"entryType"`
			Interface  string  `json:"interface"`
			Moves      int     `json:"moves,omitempty"`
			LastMove   float64 `json:"lastMove,omitempty"`
		} `json:"tableEntries"`
	} `json:"unicastTable"`
}

func (arista *AristaEOS) ListArpTable() ([]models.ArpEntry, error) {
	var response AristaMacAddressTableResponse
	err := arista.GetJsonResponse(&response, "show mac address-table")
	if err != nil {
		return nil, err
	}

	return ParseArpTable(response)
}

func ParseArpTable(output AristaMacAddressTableResponse) ([]models.ArpEntry, error) {
	var table []models.ArpEntry

	portsWithMac := map[string][]models.MacAddress{}
	for _, entry := range output.UnicastTable.TableEntries {
		if _, ok := portsWithMac[entry.Interface]; !ok {
			portsWithMac[entry.Interface] = []models.MacAddress{}
		}

		mac := models.MacAddress(entry.MacAddress)
		if !mac.Valid() {
			continue
		}

		portsWithMac[entry.Interface] = append(portsWithMac[entry.Interface], mac)
	}

	for portName, macs := range portsWithMac {
		table = append(table, models.ArpEntry{
			SwitchPort:   portName,
			MacAddresses: macs,
		})
	}

	// sort by switch port name
	sort.Slice(table, func(i, j int) bool {
		return table[i].SwitchPort < table[j].SwitchPort
	})

	return table, nil
}
