package juniper_els

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Juniper/go-netconf/netconf"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func (j *JuniperELS) GetRunningConfig() (*junosConfiguration, error) {
	reply, err := j.session.Exec(netconf.RawMethod("<get-config><source><running/></source></get-config>"))
	if err != nil {
		return nil, err
	}

	var cfg junosConfiguration
	if err := xml.Unmarshal([]byte(reply.RawReply), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetVlanMap returns the int32 representation of a specified vlan name
func (j *junosConfiguration) GetVlanIDByName(name string) (int32, error) {
	if name == "" {
		return 0, errors.New("no vlan defined")
	}
	for _, vlan := range j.VLANs {
		if vlan.Name == name {
			return vlan.ID, nil
		}
	}
	intval, err := strconv.Atoi(strings.TrimSpace(name))
	if err == nil {
		return int32(intval), nil
	}
	return 0, errors.New("GetVlanIDByName: no vlan found")
}

// GetInterfaceMode returns the mode of a specified interface
func (j *junosConfiguration) GetInterfaceMode(name string) (models.InterfaceMode, error) {
	for _, nic := range j.Interfaces {
		if nic.Name != name {
			continue
		}
		switch strings.TrimSpace(nic.PortMode) {
		case "trunk":
			return models.InterfaceModeTrunk, nil
		case "access":
			return models.InterfaceModeAccess, nil
		default: // "access" is default
			return models.InterfaceModeAccess, nil
		}
	}
	return "", fmt.Errorf("GetInterfaceMode: no interface found: %s", name)
}

// GetVlansByInterface returns the VLANs of a specified interface
func (j *junosConfiguration) GetVlansByInterface(name string) (*junosVlanMapEntry, error) {
	rgSingleVLAN := regexp.MustCompile("^[0-9]+$")
	rgVLANRange := regexp.MustCompile(`^[0-9]+-[0-9]+$`)
	for _, nic := range j.Interfaces {
		if nic.Name != name {
			continue
		}
		junosInterface := junosVlanMapEntry{}
		if rgSingleVLAN.MatchString(nic.UntaggedVLAN) {
			intId, err := strconv.Atoi(nic.UntaggedVLAN)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			junosInterface.UntaggedVLAN = int32(intId)
		} else {
			vlanID, err := j.GetVlanIDByName(nic.UntaggedVLAN)
			if err == nil {
				junosInterface.UntaggedVLAN = vlanID
			}
		}

		for _, vlan := range nic.TaggedVLANs {
			vlan = strings.TrimSpace(vlan)
			// Numeric -> to int32
			if rgSingleVLAN.MatchString(vlan) {
				intId, err := strconv.Atoi(vlan)
				if err != nil {
					log.Println(err.Error())
					continue
				}

				if int32(intId) == junosInterface.UntaggedVLAN {
					continue
				}

				junosInterface.TaggedVLANs = append(junosInterface.TaggedVLANs, int32(intId))
			} else if rgVLANRange.MatchString(vlan) {
				convertedVlans := utils.ConvertVlans(vlan, ",")
				convertedVlans = utils.DeleteVlanFromIDs(convertedVlans, junosInterface.UntaggedVLAN)
				junosInterface.TaggedVLANs = append(junosInterface.TaggedVLANs, convertedVlans...)
			} else {
				vlanID, err := j.GetVlanIDByName(nic.UntaggedVLAN)
				if err == nil {
					junosInterface.TaggedVLANs = append(junosInterface.TaggedVLANs, vlanID)
				}
			}
		}
		return &junosInterface, nil
	}
	return nil, fmt.Errorf("GetVlansByInterface: no interface found: %s", name)
}
