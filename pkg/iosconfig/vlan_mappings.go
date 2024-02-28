package iosconfig

import (
	"strconv"
	"strings"
)

// VlanMappings list all VLANs and their ids. The VLANs are taken from the config
// command.
func (cfg Config) VlanMappings() map[string]map[int32]int32 {
	vlans := make(map[string]map[int32]int32)
	for key, values := range cfg {
		if strings.HasPrefix(key, "vlan mapping table ") {
			groupName := strings.TrimSpace(strings.TrimPrefix(key, "vlan mapping table "))
			if _, ok := vlans[groupName]; !ok {
				vlans[groupName] = map[int32]int32{}
			}
			for _, value := range values {
				if strings.HasPrefix(value, "raw-vlan ") {
					configParts := strings.Split(strings.TrimSpace(value), " ")
					if configParts[0] == "raw-vlan" && len(configParts) == 4 {
						sourceVlanID, err := strconv.Atoi(configParts[1])
						if err != nil {
							continue
						}
						destinationEvc := cfg["ethernet evc "+configParts[3]]
						if destinationEvc == nil {
							continue
						}
						destinationVlanID, err := strconv.Atoi(destinationEvc.GetStringValue("dot1q mapped-vlan", ""))
						if err != nil {
							continue
						}
						vlans[groupName][int32(sourceVlanID)] = int32(destinationVlanID)

					}
				}
			}
		}
	}

	return vlans
}
