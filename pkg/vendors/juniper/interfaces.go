package juniper

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Juniper/go-netconf/netconf"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"strconv"
	"strings"
	"text/template"
)

const defaultMTU = 1500

func (j *Juniper) ListInterfaces() ([]*models.Interface, error) {
	if err := j.updateVlanMap(); err != nil {
		return nil, err
	}
	reply, err := j.session.Exec(netconf.RawMethod("<get-interface-information><level>extensive</level></get-interface-information>"))
	if err != nil {
		return nil, err
	}

	var interfaces *junosInterfaces
	if err := xml.Unmarshal([]byte(reply.RawReply), &interfaces); err != nil {
		return nil, err
	}

	// TODO: add management, mode
	resp := make([]*models.Interface, 0)
	for _, physicalInterface := range interfaces.Interfaces {
		physicalInterfaceName := strings.TrimSpace(physicalInterface.Name)
		nic := &models.Interface{
			Name:        physicalInterfaceName,
			Description: strings.TrimSpace(physicalInterface.Description),
			Enabled:     strings.TrimSpace(physicalInterface.AdminStatus) == "up",
			MTU:         defaultMTU,
			Speed:       physicalInterface.GetSpeed(),
		}

		mtu, err := strconv.Atoi(physicalInterface.MTU)
		if err == nil && mtu > 0 {
			nic.MTU = uint32(mtu)
		}

		for _, logicalInterface := range physicalInterface.LogicalInterfaces {
			if len(logicalInterface.AddressFamily.Addresses) > 0 {
				nic.Management = true
				break
			}
		}

		mac := models.MacAddress(strings.TrimSpace(physicalInterface.MacAddress))
		if mac.Valid() {
			nic.MacAddress = mac
		}

		resp = append(resp, nic)
	}

	// Add VLANs to interfaces
	for i := range resp {
		if vlanCfg, ok := j.interfaceVlans[resp[i].Name]; ok {
			resp[i].TaggedVLANs = vlanCfg.TaggedVLANs
			if vlanCfg.UntaggedVLAN > 0 {
				resp[i].UntaggedVLAN = &vlanCfg.UntaggedVLAN
			}
		}
	}

	return resp, nil
}

func (j *Juniper) GetInterface(name string) (*models.Interface, error) {
	ifaces, err := j.ListInterfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Name == name {
			return iface, nil
		}
	}
	return nil, fmt.Errorf("interface %s not found", name)
}

func (s *Juniper) isUplink(inter string) bool {
	inter = strings.TrimSpace(inter)

	// Skip et, xe, ae interfaces for now
	if strings.HasPrefix(inter, "xe-") ||
		strings.HasPrefix(inter, "te-") ||
		strings.HasPrefix(inter, "ae") {
		return true
	}
	return false
}

func (j *Juniper) ConfigureInterface(update *models.UpdateInterface) (bool, error) {
	configMutex.Lock(j.identifier)
	defer configMutex.Unlock(j.identifier)

	if j.isUplink(update.Name) {
		return false, fmt.Errorf("never configure the uplink port %s", update.Name)
	}

	swport, err := j.GetInterface(update.Name)
	if err != nil {
		return false, err
	}
	if !swport.Differs(update) {
		return false, nil
	}

	config := `<edit-config>
	<target>
		<candidate/>
	</target>
    <default-operation>merge</default-operation>
    <test-option>test-then-set</test-option>
	<config>
		<configuration>
			<interfaces>
				<interface operation="replace">
					<name>{{ .Name }}</name>
					<description>{{ .Description }}</description>
					<unit>
						<name>0</name>
						<family>
							<ethernet-switching>
								<port-mode>trunk</port-mode>
								{{if gt (len .TaggedVLANs) 0}}<vlan>
									{{range .TaggedVLANs }}<members>{{ . }}</members>{{end}}
								</vlan>{{end}}
								{{if .UntaggedVLAN}}<native-vlan-id>{{ .UntaggedVLAN }}</native-vlan-id>{{end}}
							</ethernet-switching>
						</family>       
					</unit>             
				</interface>   
			</interfaces>
		</configuration>
	</config>
</edit-config>
`

	if len(update.TaggedVLANs) == 0 && update.UntaggedVLAN != nil && *update.UntaggedVLAN == 0 {
		return false, errors.New("switch port has no vlans to configure")
	}

	var tpl bytes.Buffer
	tmpl, err := template.New("").Parse(config)
	if err != nil {
		return false, err
	}

	if err = tmpl.Execute(&tpl, update); err != nil {
		return false, err
	}

	_, err = j.session.Exec(netconf.RawMethod(tpl.String()))
	if err != nil {
		return false, err
	}

	_, err = j.session.Exec(netconf.RawMethod("<commit/>"))
	if err != nil {
		return false, err
	}

	return true, nil
}
