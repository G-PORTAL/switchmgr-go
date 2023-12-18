package juniper

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Juniper/go-netconf/netconf"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"strconv"
	"strings"
)

const defaultMTU = 1500

func (j *Juniper) ListInterfaces() ([]*models.Interface, error) {
	cfg, err := j.GetRunningConfig()
	if err != nil {
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

		// Define type of interface
		switch strings.Split(physicalInterfaceName, ".")[0] {
		case "xe-":
			nic.Type = models.InterfaceType10GSFPPlus
		case "et-":
			nic.Type = models.InterfaceType40QGSFPPlus
		case "ae":
			nic.Type = models.InterfaceTypeLAG
		default:
			nic.Type = models.InterfaceType1GE
		}

		// Get mode of the interface based on switch configuration
		if portMode, err := cfg.GetInterfaceMode(physicalInterfaceName); err == nil {
			nic.Mode = portMode
		}

		// Set VLANs of the interface based on switch configuration
		if vlanCfg, err := cfg.GetVlansByInterface(physicalInterfaceName); err == nil {
			nic.TaggedVLANs = vlanCfg.TaggedVLANs
			if vlanCfg.UntaggedVLAN > 0 {
				nic.UntaggedVLAN = &vlanCfg.UntaggedVLAN
			}
		}

		// Set MTU of the interface
		mtu, err := strconv.Atoi(strings.TrimSpace(physicalInterface.MTU))
		if err == nil && mtu > 0 {
			nic.MTU = uint32(mtu)
		}

		nic.IPAddresses = make([]string, 0)
		for _, logicalInterface := range physicalInterface.LogicalInterfaces {
			// Check if interface is used for management
			if len(logicalInterface.AddressFamily.Addresses) > 0 {
				nic.Management = true
				break
			}

			// Set IP addresses of the interface
			for _, address := range logicalInterface.AddressFamily.Addresses {
				nic.IPAddresses = append(nic.IPAddresses, strings.TrimSpace(address.IP))
			}

			// Add LAG interfaces to the interface if it is a LAG
			if nic.Type == models.InterfaceTypeLAG {
				nic.LagInterfaces = make([]string, 0)
				if strings.TrimSpace(logicalInterface.Name) == fmt.Sprintf("%s.0", physicalInterfaceName) {
					for _, lagLink := range logicalInterface.LagTrafficStatistics.Links {
						nic.LagInterfaces = append(nic.LagInterfaces, strings.TrimSuffix(strings.TrimSpace(lagLink.Name), ".0"))
					}
					break
				}
			}
		}

		// Set MAC address of the interface
		mac := models.MacAddress(strings.TrimSpace(physicalInterface.MacAddress))
		if mac.Valid() {
			nic.MacAddress = mac
		}

		resp = append(resp, nic)
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

func (j *Juniper) isUplink(inter string) bool {
	inter = strings.TrimSpace(inter)

	// Skip et, xe, ae interfaces for now
	if strings.HasPrefix(inter, "xe-") ||
		strings.HasPrefix(inter, "te-") ||
		strings.HasPrefix(inter, "ae") {
		return true
	}
	return false
}

const EditPortConfigurationTemplate = `<edit-config>
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
{{if .Description}}					<description>{{ .Description }}</description>
{{end}}{{if .Disabled}}					<disable/>
{{end}}					<unit>
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

	update.Fill(swport)
	if !swport.Differs(update) {
		return false, nil
	}

	if len(update.TaggedVLANs) == 0 && update.UntaggedVLAN != nil && *update.UntaggedVLAN == 0 {
		return false, errors.New("switch port has no vlans to configure")
	}

	tpl, err := update.Template(EditPortConfigurationTemplate)
	if err != nil {
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
