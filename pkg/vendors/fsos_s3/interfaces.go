package fsos_s3

import (
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type InterfaceMode string

const (
	InterfaceModeAccess InterfaceMode = "access"
	InterfaceModeTrunk  InterfaceMode = "trunk"
	InterfaceModeHybrid InterfaceMode = "hybrid"
)

func (fs *FSComS3) ListInterfaces() ([]*models.Interface, error) {
	config, err := fs.GetConfiguration()
	if err != nil {
		return nil, err
	}

	ports, err := config.ListInterfaces()
	if err != nil {
		return nil, err
	}

	interfaceInfo, err := fs.getInterfaceInfo()
	if err != nil {
		return nil, err
	}

	for i, nic := range ports {
		if info, ok := interfaceInfo[nic.Name]; ok {
			ports[i].MacAddress = models.MacAddress(info.MacAddress.String())
			ports[i].MTU = info.MTU
			ports[i].Speed = info.Speed
		}
	}

	return ports, nil
}

func (fs *FSComS3) GetInterface(name string) (*models.Interface, error) {
	nics, err := fs.ListInterfaces()
	if err != nil {
		return nil, err
	}

	for _, nic := range nics {
		if nic.Name == name {
			return nic, nil
		}
	}

	return nil, fmt.Errorf("interface %s not found", name)
}

func (fs *FSComS3) ConfigureInterface(update *models.UpdateInterface) (bool, error) {
	nic, err := fs.GetInterface(update.Name)
	if err != nil {
		return false, err
	}

	// Check if the interface is already configured as requested
	if !nic.Differs(update) {
		fs.Logger().Debugf("interface %s is already configured as requested", update.Name)
		return false, nil
	}

	commands := []string{
		"config",                                 // enter config mode
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

	if update.TaggedVLANs != nil {
		taggedVLANs := make([]string, 0)
		if update.UntaggedVLAN != nil {
			taggedVLANs = append(taggedVLANs, strconv.Itoa(int(*update.UntaggedVLAN)))
		}

		for _, vlan := range update.TaggedVLANs {
			taggedVLANs = append(taggedVLANs, strconv.Itoa(int(vlan)))
		}

		commands = append(commands, fmt.Sprintf("switchport trunk vlan-allowed %s", strings.Join(taggedVLANs, ",")))
	}

	if update.UntaggedVLAN != nil {
		commands = append(commands,
			"no switchport trunk vlan-untagged",
			fmt.Sprintf("switchport pvid %d", *update.UntaggedVLAN)) // set untagged vlan
	}

	// exit interface config mode
	commands = append(commands, "exit", "exit")
	_, err = fs.SendCommands(commands...)
	if err != nil {
		return false, fmt.Errorf("failed to configure interface: %w", err)
	}

	nic, err = fs.GetInterface(update.Name)
	if err != nil {
		return false, err
	}

	if nic.Differs(update) {
		return false, fmt.Errorf("interface differs, original %v, updated %v", nic, update)
	}

	err = fs.Save()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (fs *FSComS3) getInterfaceInfo() (map[string]fscomInterface, error) {
	output, err := fs.SendCommands("show interface")
	if err != nil {
		return nil, err
	}

	return ParseInterfaces(output)
}

var interfaceRgx = regexp.MustCompile(`([a-zA-Z0-9\/]+) is (down|up),.+\n(?:\s+.+\n)+\s+.+\s+.+[a|A]ddress\sis\s([a-z0-9]{4}\.[a-z0-9]{4}\.[a-z0-9]{4}).+\n\s+(.+\n\s\s)?MTU\s([0-9]+)\s.+BW\s([0-9]+)`)

type fscomInterface struct {
	MacAddress net.HardwareAddr
	MTU        uint32
	Speed      uint32
}

func ParseInterfaces(output string) (map[string]fscomInterface, error) {
	interfaceInfo := make(map[string]fscomInterface)

	matches := interfaceRgx.FindAllStringSubmatch(output, -1)
	for _, match := range matches {
		nic := match[1]
		mac := match[3]

		macAddress, err := net.ParseMAC(mac)
		if err != nil {
			return nil, err
		}

		mtu := uint32(1500)
		if match[5] != "" {
			mtuInt, err := strconv.Atoi(match[5])
			if err == nil {
				mtu = uint32(mtuInt)
			}
		}

		speed := uint32(1000000)
		if match[6] != "" {
			speedInt, err := strconv.Atoi(match[6])
			if err == nil {
				speed = uint32(speedInt)
			}
		}

		interfaceInfo[nic] = fscomInterface{
			MacAddress: macAddress,
			MTU:        mtu,
			Speed:      speed,
		}
	}

	return interfaceInfo, nil
}
