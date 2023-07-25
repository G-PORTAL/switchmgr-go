package juniper

import (
	"errors"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"strconv"
	"strings"
)

type junosNif struct {
	Name         string   `xml:"name"`
	UntaggedVLAN string   `xml:"unit>family>ethernet-switching>native-vlan-id"`
	TaggedVLANs  []string `xml:"unit>family>ethernet-switching>vlan>members"`
	PortMode     string   `xml:"unit>family>ethernet-switching>port-mode"`
}

type junosVlanMapEntry struct {
	UntaggedVLAN int32
	TaggedVLANs  []int32
}

type junosInterfaces struct {
	Interfaces []junosPhysicalNif `xml:"interface-information>physical-interface"`
}

type junosPhysicalNif struct {
	Name              string `xml:"name"`
	Description       string `xml:"description"`
	AdminStatus       string `xml:"admin-status"`
	MTU               string `xml:"mtu"`
	MacAddress        string `xml:"hardware-physical-address"`
	LogicalInterfaces []struct {
		Name          string `xml:"name"`
		AddressFamily struct {
			Name  string `xml:"address-family-name"`
			Flags struct {
				Trunk *bool `xml:"ifff-port-mode-trunk,omitempty"`
			} `xml:"address-family-flags"`
			Addresses []struct {
				IP string `xml:"ifa-local"`
			} `xml:"interface-address"`
		} `xml:"address-family"`
		LagTrafficStatistics struct {
			Links []struct {
				Name string `xml:"name"`
			} `xml:"lag-link"`
		} `xml:"lag-traffic-statistics"`
	} `xml:"logical-interface"`
	Speed                   string `xml:"speed"`
	EthernetAutonegotiation struct {
		Text                  string `xml:",chardata"`
		AutonegotiationStatus string `xml:"autonegotiation-status"`
		LinkPartnerStatus     string `xml:"link-partner-status"`
		LinkPartnerDuplexity  string `xml:"link-partner-duplexity"`
		LinkPartnerSpeed      string `xml:"link-partner-speed"`
		FlowControl           string `xml:"flow-control"`
		LocalInfo             struct {
			Text             string `xml:",chardata"`
			LocalFlowControl string `xml:"local-flow-control"`
			LocalRemoteFault string `xml:"local-remote-fault"`
		} `xml:"local-info"`
	} `xml:"ethernet-autonegotiation"`
}

// Requests and responses
type junosVlan struct {
	Name string `xml:"name"`
	ID   int32  `xml:"vlan-id"`
}

type junosConfiguration struct {
	VLANs      []junosVlan `xml:"data>configuration>vlans>vlan"`
	Interfaces []junosNif  `xml:"data>configuration>interfaces>interface"`
}

type junosArpTable struct {
	Entries []struct {
		MacAddress models.MacAddress `xml:"mac-address"`
		Interfaces []struct {
			Interface string `xml:"mac-interfaces"`
		} `xml:"mac-interfaces-list"`
	} `xml:"ethernet-switching-table-information>ethernet-switching-table>mac-table-entry"`
}

type junosSystemInformation struct {
	HardwareModel string `xml:"system-information>hardware-model"`
	OsName        string `xml:"system-information>os-name"`
	OsVersion     string `xml:"system-information>os-version"`
	SerialNumber  string `xml:"system-information>serial-number"`
	HostName      string `xml:"system-information>host-name"`
}

type junosLLDPInformation struct {
	Entries []struct {
		// LocalPortID Junos 16 or higher
		LocalPortID string `xml:"lldp-local-port-id"`
		// LocalPortID Junos 15 or lower
		LocalInterface string `xml:"lldp-local-interface"`
		RemoteHostname string `xml:"lldp-remote-system-name"`
	} `xml:"lldp-neighbors-information>lldp-neighbor-information"`
}

// GetSpeed Tries to convert the speed of the interface if set, otherwise it uses the speed of the link partner
func (p *junosPhysicalNif) GetSpeed() uint32 {
	speedValue := strings.TrimSpace(strings.ToLower(p.Speed))
	parsedSpeed, err := p.parseSpeed(speedValue)
	if err == nil {
		return parsedSpeed
	}
	speedValue = strings.TrimSpace(strings.ToLower(p.EthernetAutonegotiation.LinkPartnerSpeed))
	parsedSpeed, err = p.parseSpeed(speedValue)
	if err == nil {
		return parsedSpeed
	}
	return 0
}

// parseSpeed parses the speed string and returns the speed in kbit/s if possible
func (p *junosPhysicalNif) parseSpeed(speed string) (uint32, error) {
	if strings.HasSuffix(speed, "gbps") {
		gpps, err := strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(speed, "gbps")))
		if err == nil {
			return uint32(gpps * 1000000), nil
		}
	}
	if strings.HasSuffix(speed, "mbps") {
		mbps, err := strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(speed, "mbps")))
		if err == nil {
			return uint32(mbps * 1000), nil
		}
	}
	return 0, errors.New("invalid speed format")
}
