package fsos_n5

import (
	"bufio"
	"errors"
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

func (fs *FSComN5) ListInterfaces() ([]*models.Interface, error) {
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

func (fs *FSComN5) GetInterface(name string) (*models.Interface, error) {
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

func (fs *FSComN5) ConfigureInterface(update *models.UpdateInterface) (bool, error) {
	// TODO: Currently not implemented since it's not required at this point.e
	return false, errors.New("not implemented")
}

func (fs *FSComN5) getInterfaceInfo() (map[string]fscomInterface, error) {
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
var interfaceBandwithRegex = regexp.MustCompile(`,\sBW\s([0-9]+)\sKbit`)
var interfaceMTURegex = regexp.MustCompile(`TMTU\s([0-9]+)\sbytes`)

func ParseInterfaces(output string) (map[string]fscomInterface, error) {
	interfaces := map[string]*fscomInterface{}
	reader := strings.NewReader(output)
	scanner := bufio.NewScanner(reader)
	currentInterface := ""
	var currentInterfaceConfig *fscomInterface
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "==========================") {
			// New interface definition reached, saving old one
			currentInterface = strings.TrimSpace(strings.ReplaceAll(line, "=", " "))
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
		return nil, err
	}
	var interfaces2 = map[string]fscomInterface{}
	for k, v := range interfaces {
		interfaces2[k] = *v
	}
	return interfaces2, nil
}
