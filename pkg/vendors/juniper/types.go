package juniper

import (
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/models"
)

type junosNif struct {
	Name         string   `xml:"name"`
	UntaggedVLAN string   `xml:"unit>family>ethernet-switching>native-vlan-id"`
	TaggedVLANs  []string `xml:"unit>family>ethernet-switching>vlan>members"`
}

type junosVlanMapEntry struct {
	UntaggedVLAN int32
	TaggedVLANs  []int32
}

type junosInterfaces struct {
	Interfaces []junosPhysicalNif `xml:"interface-information>physical-interface"`
}
type junosPhysicalInterface struct {
	FPC  int32
	PIC  int32
	Port int32
}

func (p junosPhysicalInterface) String() string {
	return fmt.Sprintf("%d/%d/%d", p.FPC, p.PIC, p.Port)
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
}

// Requests and responses
type junosVlan struct {
	Name string `xml:"name"`
	ID   int32  `xml:"vlan-id"`
}

type junosInterfacesList struct {
	Interfaces []junosNif  `xml:"configuration>interfaces>interface"`
	Vlans      []junosVlan `xml:"configuration>vlans>vlan"`
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
