package fscom

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
	"net"
	"regexp"
)

type InterfaceMode string

const (
	InterfaceModeAccess InterfaceMode = "access"
	InterfaceModeTrunk  InterfaceMode = "trunk"
	InterfaceModeHybrid InterfaceMode = "hybrid"
)

func (fs *FSCom) ListInterfaces() ([]*models.Interface, error) {
	config, err := fs.getConfiguration()
	if err != nil {
		return nil, err
	}

	ports, err := config.ListInterfaces()
	if err != nil {
		return nil, err
	}

	macToInterface, err := fs.getInterfaceToMacMap()
	if err != nil {
		return nil, err
	}

	for i, nic := range ports {
		if macAddress, ok := macToInterface[nic.Name]; ok {
			ports[i].MacAddress = models.MacAddress(macAddress.String())
		}
	}

	return ports, nil
}

func (fs *FSCom) getInterfaceToMacMap() (map[string]net.HardwareAddr, error) {
	output, err := fs.sendCommands("show interface")
	if err != nil {
		return nil, err
	}

	return ParseInterfaces(output)
}

var interfaceRgx = regexp.MustCompile("([a-zA-Z0-9\\/]+) is (down|up),.+\\n\\s+.+\\n\\s+.+\\s+.+[address|Adddress]\\sis\\s([a-z0-9]{4}\\.[a-z0-9]{4}\\.[a-z0-9]{4})")

func ParseInterfaces(output string) (map[string]net.HardwareAddr, error) {
	macToInterface := make(map[string]net.HardwareAddr)

	matches := interfaceRgx.FindAllStringSubmatch(output, -1)
	for _, match := range matches {
		nic := match[1]
		mac := match[3]

		macAddress, err := net.ParseMAC(mac)
		if err != nil {
			return nil, err
		}

		macToInterface[nic] = macAddress
	}

	return macToInterface, nil
}
