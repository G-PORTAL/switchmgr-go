package iosconfig

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
)

// Vlans List of all vlans
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

func (cfg Config) VlanIDs() []int32 {
	vlanIDs := make([]int32, 0)
	for id := range cfg.Vlans() {
		vlanIDs = append(vlanIDs, id)
	}

	slices.Sort(vlanIDs)

	return vlanIDs
}

func (cfg Config) Vlan(id int32) (ConfigValues, error) {
	for vlanID, config := range cfg.Vlans() {
		if vlanID == id {
			return config, nil
		}
	}

	return nil, fmt.Errorf("vlan %d not found", id)
}
