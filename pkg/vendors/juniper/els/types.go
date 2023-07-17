package els

import "github.com/g-portal/switchmgr-go/pkg/models"

type junosNif struct {
	Name          string   `xml:"name"`
	Description   string   `xml:"description"`
	NativeVlanID  string   `xml:"native-vlan-id"`
	TaggedVlanIDs []string `xml:"unit>family>ethernet-switching>vlan>members"`
}

// Requests and responses
type junosVlan struct {
	Name   string `xml:"name"`
	VlanID int32  `xml:"vlan-id"`
}

type junosInterfacesList struct {
	Interfaces []junosNif  `xml:"configuration>interfaces>interface"`
	Vlans      []junosVlan `xml:"configuration>vlans>vlan"`
}

type junosConfiguration struct {
	VLANs []junosVlan `xml:"data>configuration>vlans>vlan"`
}

type junosArpTable struct {
	Entries []struct {
		MacAddress models.MacAddress `xml:"l2ng-l2-mac-address"`
		Interface  string            `xml:"l2ng-l2-mac-logical-interface"`
	} `xml:"l2ng-l2ald-rtb-macdb>l2ng-l2ald-mac-entry-vlan>l2ng-mac-entry"`
}

type junosSystemInformation struct {
	HardwareModel string `xml:"system-information>hardware-model"`
	OsName        string `xml:"system-information>os-name"`
	OsVersion     string `xml:"system-information>os-version"`
	SerialNumber  string `xml:"system-information>serial-number"`
	HostName      string `xml:"system-information>host-name"`
}
