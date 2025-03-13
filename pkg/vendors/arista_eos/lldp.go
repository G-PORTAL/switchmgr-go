package arista_eos

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
)

type AristaLLDPResponse struct {
	LldpNeighbors []struct {
		Port           string `json:"port"`
		NeighborDevice string `json:"neighborDevice"`
		NeighborPort   string `json:"neighborPort"`
		TTL            int    `json:"ttl"`
	} `json:"lldpNeighbors"`
}

func (arista *AristaEOS) ListLLDPNeighbors() ([]models.LLDPNeighbor, error) {
	var response AristaLLDPResponse
	err := arista.GetJsonResponse(&response, "show lldp neighbors")
	if err != nil {
		return nil, err
	}

	return ParseLLDPNeighbors(response)
}

func ParseLLDPNeighbors(output AristaLLDPResponse) ([]models.LLDPNeighbor, error) {
	var table []models.LLDPNeighbor
	for _, neighbor := range output.LldpNeighbors {
		table = append(table, models.LLDPNeighbor{
			LocalInterface: neighbor.Port,
			RemoteHostname: neighbor.NeighborDevice,
		})
	}

	return table, nil
}
