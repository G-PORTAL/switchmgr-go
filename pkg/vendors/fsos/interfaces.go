package fsos

import (
	"bufio"
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

func (fs *FSOS) ListInterfaces() ([]*models.Interface, error) {
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

func (fs *FSOS) GetInterface(name string) (*models.Interface, error) {
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

func (fs *FSOS) ConfigureInterface(update *models.UpdateInterface) (bool, error) {
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

	if update.Description != nil {
		commands = append(commands, fmt.Sprintf("description %s", *update.Description)) // set interface description
	}

	if update.UntaggedVLAN != nil {
		commands = append(commands, fmt.Sprintf("switchport trunk native vlan %d", *update.UntaggedVLAN)) // set untagged vlan
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

func (fs *FSOS) getInterfaceInfo() (map[string]fscomInterface, error) {
	output, err := fs.SendCommands("show interface")
	if err != nil {
		return nil, err
	}

	return ParseInterfaces(output)
}

type fscomInterface struct {
	MacAddress net.HardwareAddr
	MTU        uint32
	Speed      uint32
}

var interfaceMacAddressRegex = regexp.MustCompile(`address is ([0-9a-fA-F.]+)`)
var interfaceBandwithRegex = regexp.MustCompile(`Bandwidth (\d+) kbits`)
var interfaceMTURegex = regexp.MustCompile(`The maximum transmit unit \(MTU\) is (\d+) bytes`)

func ParseInterfaces(output string) (map[string]fscomInterface, error) {
	interfaces := map[string]*fscomInterface{}
	reader := strings.NewReader(output)
	scanner := bufio.NewScanner(reader)
	currentInterface := ""
	var currentInterfaceConfig *fscomInterface
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Interface ") {
			// New interface definition reached, saving old one
			currentInterface = strings.TrimSpace(strings.TrimPrefix(line, "Interface "))
		}
		if _, ok := interfaces[currentInterface]; !ok {
			currentInterfaceConfig = &fscomInterface{
				Speed: uint32(1000000),
				MTU:   uint32(1500),
			}
			interfaces[currentInterface] = currentInterfaceConfig
			continue
		}
		if currentInterface != "" {
			if match := interfaceMacAddressRegex.FindStringSubmatch(line); len(match) > 1 {
				if mac, err := net.ParseMAC(match[1]); err == nil {
					currentInterfaceConfig.MacAddress = mac
				}
			}
			if match := interfaceBandwithRegex.FindStringSubmatch(line); len(match) > 1 {
				if bw, err := strconv.Atoi(match[1]); err == nil {
					currentInterfaceConfig.Speed = uint32(bw)
				}
			}
			if match := interfaceMTURegex.FindStringSubmatch(line); len(match) > 1 {
				if mtu, err := strconv.Atoi(match[1]); err == nil {
					currentInterfaceConfig.MTU = uint32(mtu)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
		return nil, err
	}
	var interfaces2 = map[string]fscomInterface{}
	for k, v := range interfaces {
		interfaces2[k] = *v
	}
	return interfaces2, nil
}
