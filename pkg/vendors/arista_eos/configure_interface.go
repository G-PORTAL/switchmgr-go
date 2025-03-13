package arista_eos

import (
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"strconv"
	"strings"
)

type InterfaceMode string

const (
	InterfaceModeAccess InterfaceMode = "access"
	InterfaceModeTrunk  InterfaceMode = "trunk"
)

func (arista *AristaEOS) ConfigureInterface(update *models.UpdateInterface) (bool, error) {
	nic, err := arista.GetInterface(update.Name)
	if err != nil {
		return false, err
	}

	// Check if the interface is already configured as requested
	if !nic.Differs(update) {
		arista.Logger().Debugf("interface %s is already configured as requested", update.Name)
		return false, nil
	}

	commands := []string{
		"configure",                              // enter config mode
		fmt.Sprintf("interface %s", update.Name), // enter interface config mode,
		fmt.Sprintf("switchport mode %s", InterfaceModeTrunk), // set interface mode
	}

	if update.Enabled != nil && !*update.Enabled {
		commands = append(commands, "shutdown")
	} else if update.Enabled != nil && *update.Enabled {
		commands = append(commands, "no shutdown")
	}

	if update.Description != nil {
		commands = append(commands, fmt.Sprintf("description %s", *update.Description)) // set interface description
	}

	if update.UntaggedVLAN != nil && *update.UntaggedVLAN != 1 {
		commands = append(commands, fmt.Sprintf("switchport trunk native vlan %d", *update.UntaggedVLAN)) // set untagged vlan
	} else if update.UntaggedVLAN != nil && *update.UntaggedVLAN == 1 {
		commands = append(commands, "no switchport trunk native vlan")
	}

	if update.TaggedVLANs != nil {
		taggedVLANs := make([]string, 0)
		if update.UntaggedVLAN != nil {
			taggedVLANs = append(taggedVLANs, strconv.Itoa(int(*update.UntaggedVLAN)))
		}

		for _, vlan := range update.TaggedVLANs {
			taggedVLANs = append(taggedVLANs, strconv.Itoa(int(vlan)))
		}

		commands = append(commands,
			"switchport trunk allowed vlan none",
			fmt.Sprintf("switchport trunk allowed vlan add %s", strings.Join(taggedVLANs, ",")))
	}

	// VLAN translation
	if update.VlanMapping != nil {
		for originalID, newID := range nic.VlanMapping {
			commands = append(commands, fmt.Sprintf("no switchport vlan translation %d %d", originalID, newID))
		}

		for originalVLAN, newVLAN := range update.VlanMapping {
			commands = append(commands, fmt.Sprintf("switchport vlan translation %d %d", originalVLAN, newVLAN))
		}
	}

	// exit interface config mode
	commands = append(commands, "exit", "exit")

	_, err = arista.SendCommands(commands...)
	if err != nil {
		return false, fmt.Errorf("failed to configure interface: %w", err)
	}

	nic, err = arista.GetInterface(update.Name)
	if err != nil {
		return false, err
	}

	if nic.Differs(update) {
		return false, fmt.Errorf("interface differs, original %v, updated %v", nic, update)
	}

	err = arista.Save()
	if err != nil {
		return false, err
	}

	return true, nil
}
