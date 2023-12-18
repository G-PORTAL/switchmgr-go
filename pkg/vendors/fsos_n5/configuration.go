package fsos_n5

import (
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/iosconfig"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/utils"
	"strconv"
	"strings"
)

type Configuration iosconfig.Config

func (cfg Configuration) ListInterfaces() ([]*models.Interface, error) {
	iosConfig := iosconfig.Config(cfg)

	vlanIDs := make([]int32, 0)
	for id := range iosConfig.Vlans() {
		vlanIDs = append(vlanIDs, id)
	}

	interfaces := make([]*models.Interface, 0)
	for nic, config := range iosConfig.Interfaces() {
		mgmt := false

		mode := InterfaceMode(config.GetStringValue("switchport mode", string(InterfaceModeAccess)))

		var untaggedVLAN *int32
		var taggedVLANs []int32

		enable := !config.Exists("shutdown", true)

		interfaceMode := models.InterfaceModeAccess

		defaultVLAN := int32(1)
		if strings.HasPrefix(nic, "VLAN") {
			var err error
			defaultVLANInt, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(nic, "VLAN")))
			// check if err is nil and defaultVLANInt can be converted to int32
			if err == nil && int32(defaultVLANInt) > 0 {
				defaultVLAN = int32(defaultVLANInt)
			}
		}

		switch mode {
		case InterfaceModeAccess:
			accessVlanID := config.GetInt32Value("switchport access vlan", defaultVLAN)
			untaggedVLAN = &accessVlanID
		case InterfaceModeTrunk:
			interfaceMode = models.InterfaceModeTrunk
			accessVlanID := defaultVLAN
			if config.Exists("switchport trunk vlan-untagged", false) {
				accessVlanID = config.GetInt32Value("switchport trunk vlan-untagged", defaultVLAN)
			} else if config.Exists("switchport trunk native vlan", false) {
				accessVlanID = config.GetInt32Value("switchport trunk native vlan", defaultVLAN)
			}

			untaggedVLAN = &accessVlanID
			taggedVLANs = config.GetInt32Values("switchport trunk allowed vlan add", vlanIDs)

			if config.Exists("switchport trunk allowed vlan add", false) {
				taggedVLANs = utils.ConvertVlans(config.GetStringValue("switchport trunk allowed vlan add", ""), ",")
			}

			// remove untagged vlan from tagged VLANs
			taggedVLANs = utils.DeleteVlanFromIDs(taggedVLANs, *untaggedVLAN)
		case InterfaceModeHybrid:
			interfaceMode = models.InterfaceModeTrunk
			//TODO implement hybrid mode
		default:
			return nil, fmt.Errorf("unknown interface mode %q", mode)
		}

		if config.Exists("ip address", false) {
			mgmt = true
		}

		var vlanMappingName *string
		if mapping := config.GetStringValue("switchport trunk vlan-translation mapping table", ""); mapping != "" {
			vlanMappingName = &mapping
		}

		var portIsolationGroup *string
		if group := config.GetStringValue("port-isolate group", ""); group != "" {
			portIsolationGroup = &group
		}

		interfaces = append(interfaces, &models.Interface{
			Name:               nic,
			Description:        config.GetStringValue("description", ""),
			Enabled:            enable,
			Mode:               interfaceMode,
			UntaggedVLAN:       untaggedVLAN,
			TaggedVLANs:        taggedVLANs,
			Management:         mgmt,
			VlanMappingName:    vlanMappingName,
			PortIsolationGroup: portIsolationGroup,
		})
	}

	return interfaces, nil
}

// GetConfiguration returns the configuration of a FSComS3 switch.
func (fs *FSComN5) GetConfiguration() (*Configuration, error) {
	output, err := fs.SendCommands("show running-config")
	if err != nil {
		return nil, err
	}

	cfg := ParseConfiguration(output)
	config := Configuration(cfg)

	return &config, nil
}

func (cfg Configuration) GetInterface(name string) (*models.Interface, error) {
	nics, err := cfg.ListInterfaces()
	if err != nil {
		return nil, err
	}

	for _, nic := range nics {
		if nic.Name == name {
			return nic, nil
		}
	}

	return nil, fmt.Errorf("interface %q not found", name)
}

// ParseConfiguration parses the configuration of a FSComS3 switch.
func ParseConfiguration(cfg string) iosconfig.Config {
	return iosconfig.Parse(cfg)
}
