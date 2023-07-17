package juniper_els

import (
	"encoding/xml"
	"github.com/Juniper/go-netconf/netconf"
	"github.com/g-portal/switchmgr-go/pkg/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func (j *JuniperELS) updateVlanMap() error {
	cfg, err := j.GetRunningConfig()
	if err != nil {
		return err
	}

	j.vlanMapping = make(map[string]int32)
	j.interfaceVlans = make(map[string]junosVlanMapEntry)
	for _, vlan := range cfg.VLANs {
		j.vlanMapping[vlan.Name] = vlan.ID
	}

	rgSingleVLAN := regexp.MustCompile("^[0-9]+$")
	rgVLANRange := regexp.MustCompile(`^[0-9]+-[0-9]+$`)
	for _, nic := range cfg.Interfaces {
		junosInterface := junosVlanMapEntry{}
		if rgSingleVLAN.MatchString(nic.UntaggedVLAN) {
			intId, err := strconv.Atoi(nic.UntaggedVLAN)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			junosInterface.UntaggedVLAN = int32(intId)
		} else {
			if id, ok := j.vlanMapping[nic.UntaggedVLAN]; ok {
				junosInterface.UntaggedVLAN = id
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

				// If the tagged vlan is the same as the untagged, we skip it.
				if int32(intId) == junosInterface.UntaggedVLAN {
					continue
				}

				junosInterface.TaggedVLANs = append(junosInterface.TaggedVLANs, int32(intId))
			} else if rgVLANRange.MatchString(vlan) {
				convertedVlans := utils.ConvertVlans(vlan, ",")
				convertedVlans = utils.DeleteVlanFromIDs(convertedVlans, junosInterface.UntaggedVLAN)
				junosInterface.TaggedVLANs = append(junosInterface.TaggedVLANs, convertedVlans...)
			} else {
				if id, ok := j.vlanMapping[vlan]; ok {
					if id == junosInterface.UntaggedVLAN {
						continue
					}

					junosInterface.TaggedVLANs = append(junosInterface.TaggedVLANs, id)
				}
			}
		}

		j.interfaceVlans[nic.Name] = junosInterface
	}
	return nil
}
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
