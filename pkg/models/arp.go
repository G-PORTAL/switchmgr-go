package models

type ArpEntry struct {
	// SwitchPort Name of the switch port interface
	SwitchPort string

	// MacAddresses List of mac addresses seen on the switch port
	MacAddresses []MacAddress
}
