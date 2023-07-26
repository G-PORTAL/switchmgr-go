package iosconfig

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
)

// Vlans list all VLANs and their ids. The VLANs are taken from the config
// command.
func (cfg Config) Vlans() map[int32]ConfigValues {
	vlans := make(map[int32]ConfigValues)
	for key, values := range cfg {
		if strings.HasPrefix(key, "vlan ") {
			vlansIDStrings := strings.Split(strings.TrimSpace(strings.TrimPrefix(key, "vlan ")), ",")
			for _, d := range vlansIDStrings {
				vlanID, err := strconv.Atoi(d)
				if err != nil {
					continue
				}

				vlans[int32(vlanID)] = values
			}
		}
	}

	return vlans
}

// VlanIDs returns a list of all VLAN IDs, specified in the config. The list is
// sorted.
func (cfg Config) VlanIDs() []int32 {
	vlanIDs := make([]int32, 0)
	for id := range cfg.Vlans() {
		vlanIDs = append(vlanIDs, id)
	}

	slices.Sort(vlanIDs)

	return vlanIDs
}

// Vlan returns the configuration of a specific VLAN. If the VLAN is not found,
// an error is returned.
func (cfg Config) Vlan(id int32) (ConfigValues, error) {
	for vlanID, config := range cfg.Vlans() {
		if vlanID == id {
			return config, nil
		}
	}

	return nil, fmt.Errorf("vlan %d not found", id)
}
