package juniper

import (
	"encoding/xml"
	"github.com/Juniper/go-netconf/netconf"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"strings"
)

func ParseGetInterfaceInformation(data []byte) ([]*models.Interface, error) {
	var interfaces *junosInterfaces
	if err := xml.Unmarshal(data, &interfaces); err != nil {
		return nil, err
	}

	resp := make([]*models.Interface, 0)
	for _, physicalInterface := range interfaces.Interfaces {
		physicalInterfaceName := strings.TrimSpace(physicalInterface.Name)
		if !isValidInterface(physicalInterfaceName) {
			continue
		}
		mac := strings.TrimSpace(physicalInterface.MacAddress)
		resp = append(resp, &models.Interface{
			Name:       physicalInterfaceName,
			MacAddress: models.MacAddress(mac),
		})
	}
	return resp, nil
}

func (j *Juniper) ListInterfaces() ([]*models.Interface, error) {

	if err := j.updateVlanMap(); err != nil {
		return nil, err
	}
	reply, err := j.session.Exec(netconf.RawMethod("<get-interface-information><level>extensive</level></get-interface-information>"))
	if err != nil {
		return nil, err
	}

	interfaces, err := ParseGetInterfaceInformation([]byte(reply.RawReply))
	if err != nil {
		return nil, err
	}

	// Add VLANs to interfaces
	for i := range interfaces {
		if vlanCfg, ok := j.interfaceVlans[interfaces[i].Name]; ok {
			interfaces[i].TaggedVLANs = vlanCfg.TaggedVLANs
			if vlanCfg.UntaggedVLAN > 0 {
				interfaces[i].UntaggedVLAN = &vlanCfg.UntaggedVLAN
			}
		}
	}

	return interfaces, nil
}
