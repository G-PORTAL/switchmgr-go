package juniper_els

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/Juniper/go-netconf/netconf"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/utils"
	"strconv"
	"strings"
	"text/template"
)

func (j *JuniperELS) ListInterfaces() ([]*models.Interface, error) {
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
		// Name
		physicalInterfaceName := strings.TrimSpace(physicalInterface.Name)

		// IP addresses
		ips := make([]string, 0)
		for _, logicalInterface := range physicalInterface.LogicalInterfaces {
			for _, address := range logicalInterface.AddressFamily.Addresses {
				ips = append(ips, strings.TrimSpace(address.IP))
			}
		}

		// Mode
		mode := models.InterfaceModeAccess
		if detectedMode, err := cfg.GetInterfaceMode(physicalInterfaceName); err == nil {
			mode = detectedMode
		}

		// MTU
		mtu := 0
		if m, err := strconv.Atoi(strings.TrimSpace(physicalInterface.MTU)); err == nil {
			mtu = m
		}

		// Type
		t := models.InterfaceType1GE
		if strings.HasPrefix(physicalInterfaceName, "xe-") {
			t = models.InterfaceType10GSFPPlus
		}
		if strings.HasPrefix(physicalInterfaceName, "et-") {
			t = models.InterfaceType40QGSFPPlus
		}

		// LAG ports
		lagPorts := make([]string, 0)
		if strings.HasPrefix(physicalInterfaceName, "ae") {
			t = models.InterfaceTypeLAG
			// If it is a LAG, the mode is trunk
			for _, logicalInterface := range physicalInterface.LogicalInterfaces {
				if strings.TrimSpace(logicalInterface.Name) == fmt.Sprintf("%s.0", physicalInterfaceName) {
					for _, lagLink := range logicalInterface.LagTrafficStatistics.Links {
						lagPorts = append(lagPorts, strings.TrimSuffix(strings.TrimSpace(lagLink.Name), ".0"))
					}
					break
				}
			}
		}

		nic := &models.Interface{
			Name:          physicalInterfaceName,
			Mode:          mode,
			Type:          t,
			MTU:           uint32(mtu),
			Speed:         physicalInterface.GetSpeed(),
			Description:   strings.TrimSpace(physicalInterface.Description),
			Enabled:       strings.TrimSpace(physicalInterface.AdminStatus) == "up",
			LagInterfaces: lagPorts,
			IPAddresses:   ips,
		}

		// Is it a logical interface?
		for _, logicalInterface := range physicalInterface.LogicalInterfaces {
			if len(logicalInterface.AddressFamily.Addresses) > 0 {
				nic.Management = true
				break
			}
		}

		// If there is a MAC address, parse and add it
		mac := models.MacAddress(strings.TrimSpace(physicalInterface.MacAddress))
		if mac.Valid() {
			nic.MacAddress = mac
		}

		resp = append(resp, nic)
	}
	if err != nil {
		return nil, err
	}

	// Add VLANs to interfaces
	for i := range resp {
		if vlanCfg, err := cfg.GetVlansByInterface(resp[i].Name); err == nil {
			resp[i].TaggedVLANs = vlanCfg.TaggedVLANs
			if vlanCfg.UntaggedVLAN > 0 {
				resp[i].UntaggedVLAN = &vlanCfg.UntaggedVLAN
			}
		}
	}

	return resp, nil
}

func (s *JuniperELS) isUplink(inter string) bool {
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
					<description>{{ .Description }}</description>
					{{if not .Enabled}}<disable/>{{end}}
					{{if .UntaggedVLAN}}<native-vlan-id>{{ .UntaggedVLAN }}</native-vlan-id>{{end}}
					<unit>
						<name>0</name>
						<family>
							<ethernet-switching>
								<interface-mode>trunk</interface-mode>
								{{if gt (len .TaggedVLANs) 0}}<vlan>
									{{range .TaggedVLANs }}<members>{{ . }}</members>{{end}}
								</vlan>{{end}}
							</ethernet-switching>
						</family>       
					</unit>             
				</interface>   
			</interfaces>
		</configuration>
	</config>
</edit-config>
`

// ConfigureInterface configures a single interface. It returns true if the
// configuration has changed. If the interface is an uplink, it will return
// false and an error.
func (j *JuniperELS) ConfigureInterface(update *models.UpdateInterface) (bool, error) {
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

	// Add the untagged VLAN id to tagged VLANs, otherwise it won't work
	if *update.UntaggedVLAN != 0 {
		update.TaggedVLANs = append(update.TaggedVLANs, *update.UntaggedVLAN)
	}

	// Unique tagged VLAN ids
	update.TaggedVLANs = utils.UniqueVlanIDs(update.TaggedVLANs)

	if len(update.TaggedVLANs) == 0 && update.UntaggedVLAN != nil && *update.UntaggedVLAN == 0 {
		return false, fmt.Errorf("switch port has no vlans to configure")
	}

	var tpl bytes.Buffer
	tmpl, err := template.New("").Parse(EditPortConfigurationTemplate)
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

func (j *JuniperELS) GetInterface(name string) (*models.Interface, error) {
	nics, err := j.ListInterfaces()
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
