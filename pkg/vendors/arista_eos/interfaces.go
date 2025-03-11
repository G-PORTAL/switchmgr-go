package arista_eos

import (
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/iosconfig"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"net"
	"sort"
	"strconv"
	"strings"
)

type AristaInterfacesResponse struct {
	Interfaces map[string]struct {
		Name               string `json:"name"`
		ForwardingModel    string `json:"forwardingModel"`
		LineProtocolStatus string `json:"lineProtocolStatus"`
		InterfaceStatus    string `json:"interfaceStatus"`
		Hardware           string `json:"hardware"`
		InterfaceAddress   []struct {
			PrimaryIP struct {
				Address string `json:"address"`
				MaskLen int    `json:"maskLen"`
			} `json:"primaryIp"`
			SecondaryIps struct {
			} `json:"secondaryIps"`
			SecondaryIpsOrderedList []interface{} `json:"secondaryIpsOrderedList"`
			VirtualIP               struct {
				Address string `json:"address"`
				MaskLen int    `json:"maskLen"`
			} `json:"virtualIp"`
			VirtualSecondaryIps struct {
			} `json:"virtualSecondaryIps"`
			VirtualSecondaryIpsOrderedList []interface{} `json:"virtualSecondaryIpsOrderedList"`
			BroadcastAddress               string        `json:"broadcastAddress"`
			Dhcp                           bool          `json:"dhcp"`
		} `json:"interfaceAddress"`
		PhysicalAddress           string  `json:"physicalAddress"`
		BurnedInAddress           string  `json:"burnedInAddress"`
		Description               string  `json:"description"`
		Bandwidth                 int     `json:"bandwidth"`
		Mtu                       int     `json:"mtu"`
		L3MtuConfigured           bool    `json:"l3MtuConfigured"`
		L2Mru                     int     `json:"l2Mru"`
		LastStatusChangeTimestamp float64 `json:"lastStatusChangeTimestamp"`
		InterfaceStatistics       struct {
			UpdateInterval float64 `json:"updateInterval"`
			InBitsRate     float64 `json:"inBitsRate"`
			InPktsRate     float64 `json:"inPktsRate"`
			OutBitsRate    float64 `json:"outBitsRate"`
			OutPktsRate    float64 `json:"outPktsRate"`
		} `json:"interfaceStatistics"`
		InterfaceCounters struct {
			InOctets          int64   `json:"inOctets"`
			InUcastPkts       int     `json:"inUcastPkts"`
			InMulticastPkts   int     `json:"inMulticastPkts"`
			InBroadcastPkts   int     `json:"inBroadcastPkts"`
			InDiscards        int     `json:"inDiscards"`
			InTotalPkts       int     `json:"inTotalPkts"`
			OutOctets         int64   `json:"outOctets"`
			OutUcastPkts      int     `json:"outUcastPkts"`
			OutMulticastPkts  int     `json:"outMulticastPkts"`
			OutBroadcastPkts  int     `json:"outBroadcastPkts"`
			OutDiscards       int     `json:"outDiscards"`
			OutTotalPkts      int     `json:"outTotalPkts"`
			LinkStatusChanges int     `json:"linkStatusChanges"`
			LastClear         float64 `json:"lastClear"`
			TotalInErrors     int     `json:"totalInErrors"`
			InputErrorsDetail struct {
				RuntFrames      int `json:"runtFrames"`
				GiantFrames     int `json:"giantFrames"`
				FcsErrors       int `json:"fcsErrors"`
				AlignmentErrors int `json:"alignmentErrors"`
				SymbolErrors    int `json:"symbolErrors"`
				RxPause         int `json:"rxPause"`
			} `json:"inputErrorsDetail"`
			TotalOutErrors     int `json:"totalOutErrors"`
			OutputErrorsDetail struct {
				Collisions            int `json:"collisions"`
				LateCollisions        int `json:"lateCollisions"`
				DeferredTransmissions int `json:"deferredTransmissions"`
				TxPause               int `json:"txPause"`
			} `json:"outputErrorsDetail"`
			CounterRefreshTime float64 `json:"counterRefreshTime"`
		} `json:"interfaceCounters"`
		Duplex        string `json:"duplex"`
		AutoNegotiate string `json:"autoNegotiate"`
		LoopbackMode  string `json:"loopbackMode"`
		Lanes         int    `json:"lanes"`
	} `json:"interfaces"`
}

func (arista *AristaEOS) ListInterfaces() ([]*models.Interface, error) {
	var response AristaInterfacesResponse
	err := arista.GetJsonResponse(&response, "show interfaces")
	if err != nil {
		return nil, err
	}

	config, err := arista.getRunningConfig()
	if err != nil {
		return nil, err
	}

	return ParseInterfaces(response, config)
}

func (arista *AristaEOS) GetInterface(name string) (*models.Interface, error) {
	nics, err := arista.ListInterfaces()
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

func ParseInterfaces(output AristaInterfacesResponse, config iosconfig.Config) ([]*models.Interface, error) {
	var interfaces []*models.Interface

	for _, nic := range output.Interfaces {
		macAddress := ""
		tempMacAddress, err := net.ParseMAC(nic.PhysicalAddress)
		if err == nil {
			macAddress = tempMacAddress.String()
		}

		mgmt := false
		ips := make([]string, 0)
		if len(nic.InterfaceAddress) > 0 {
			for _, address := range nic.InterfaceAddress {
				if address.PrimaryIP.Address == "" {
					continue
				}

				mgmt = true
				ips = append(ips, address.PrimaryIP.Address)
			}
		}

		vlanMapping := map[int32]int32{}

		// nic.Bandwidth in bytes
		interfaceType := models.InterfaceTypeVirtual
		switch nic.Bandwidth {
		case 1000000000: //
			interfaceType = models.InterfaceType1GE
		case 10000000000: //
			interfaceType = models.InterfaceType10GE
		case 40000000000: //
			interfaceType = models.InterfaceType40QGSFPPlus
		}

		mode := "access"
		var untaggedVlan *int32
		var taggedVlans []int32

		interfaceConfig, err := config.Interface(nic.Name)
		if err == nil {
			mode = interfaceConfig.GetStringValue("switchport mode", "access")

			if interfaceConfig.Exists("switchport access vlan", true) {
				untaggedVlanFromConfig := interfaceConfig.GetInt32Value("switchport access vlan", 1)
				untaggedVlan = &untaggedVlanFromConfig
			} else {
				untaggedVlanFromConfig := interfaceConfig.GetInt32Value("switchport trunk native vlan", 1)
				untaggedVlan = &untaggedVlanFromConfig
			}

			// get vlan mapping
			mappings := interfaceConfig.GetLines("switchport vlan translation", true)
			for _, mapping := range mappings {
				vlanList := strings.Split(mapping, " ")
				if len(vlanList) != 2 {
					continue
				}

				originalVlan, err := strconv.Atoi(vlanList[0])
				if err != nil {
					continue

				}

				newVlan, err := strconv.Atoi(vlanList[1])
				if err != nil {
					continue
				}

				vlanMapping[int32(originalVlan)] = int32(newVlan)
			}

			vlanIDs := interfaceConfig.GetInt32Values("switchport trunk allowed vlan", []int32{})
			// delete untaggedVlan from taggedVlans
			for _, vlanID := range vlanIDs {
				if vlanID == *untaggedVlan {
					continue
				}
				taggedVlans = append(taggedVlans, vlanID)
			}
		}

		interfaces = append(interfaces, &models.Interface{
			Name:         nic.Name,
			Description:  nic.Description,
			Type:         interfaceType,
			Enabled:      nic.InterfaceStatus == "connected",
			MTU:          uint32(nic.Mtu),
			Speed:        uint32(nic.Bandwidth / 1000),
			Mode:         models.InterfaceMode(mode),
			MacAddress:   models.MacAddress(macAddress),
			UntaggedVLAN: untaggedVlan,
			VlanMapping:  vlanMapping,
			TaggedVLANs:  taggedVlans,
			IPAddresses:  ips,
			Management:   mgmt,
		})
	}

	// sort by interface name
	sort.Slice(interfaces, func(i, j int) bool {
		return interfaces[i].Name < interfaces[j].Name
	})

	return interfaces, nil
}
