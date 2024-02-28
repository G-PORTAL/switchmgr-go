package fsos_s5

import (
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/models"
)

func (fs *FSComS5) ListVlanMappings() ([]models.VlanMapping, error) {
	output, err := fs.SendCommands("show running-config")
	if err != nil {
		return nil, err
	}

	cfg := ParseConfiguration(output)
	if err != nil {
		return nil, err
	}

	vlanMaps := cfg.VlanMappings()
	res := make([]models.VlanMapping, 0)
	for group, mapping := range vlanMaps {
		res = append(res, models.VlanMapping{
			GroupName: group,
			Mapping: map[int32]int32{
				mapping[0]: mapping[1],
			},
		})
	}
	return res, nil
}

func (fs *FSComS5) ConfigureVlanMapping(mapping *models.VlanMapping) (bool, error) {
	commands := []string{
		"configure terminal",
	}
	for _, out := range mapping.Mapping {
		commands = append(
			commands,
			[]string{
				fmt.Sprintf("ethernet evc %s_%d", mapping.GroupName, out),
				fmt.Sprintf("dot1q mapped-vlan %d", out),
				"exit",
			}...,
		)
	}

	commands = append(commands, fmt.Sprintf("vlan mapping table %s", mapping.GroupName))

	for in, out := range mapping.Mapping {
		commands = append(
			commands,
			[]string{
				fmt.Sprintf("raw-vlan %d evc %s_%d", in, mapping.GroupName, out),
				"exit",
			}...,
		)
	}
	if _, err := fs.SendCommands(commands...); err != nil {
		return false, err
	}
	return true, nil
}

func (fs *FSComS5) DeleteVLANMapping(name string) {
	//TODO implement me
	panic("implement me")
}
