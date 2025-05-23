package arista_eos

import (
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"net"
	"strings"
)

func (arista *AristaEOS) ListVRFRoutes(vrfName string) ([]models.VRFRoute, error) {
	config, err := arista.getRunningConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get running config: %v", err)
	}

	vrfConfig := fmt.Sprintf("ip route vrf %s", vrfName)
	routes := make([]models.VRFRoute, 0)
	for _, line := range config.Values().GetLines(vrfConfig, true) {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid route line: %s", line)
		}

		_, ipNet, err := net.ParseCIDR(parts[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse CIDR %s: %v", parts[1], err)
		}

		routes = append(routes, models.VRFRoute{
			Network:       *ipNet,
			InterfaceName: parts[1],
		})
	}

	return routes, nil
}

func (arista *AristaEOS) AddVRFRoute(vrfName, interfaceName string, network *net.IPNet) error {
	config, err := arista.getRunningConfig()
	if err != nil {
		return fmt.Errorf("failed to get running config: %v", err)
	}

	if network == nil {
		return fmt.Errorf("ip network is nil")
	}

	vrfConfig := fmt.Sprintf("ip route vrf %s %s %s", vrfName, network.String(), interfaceName)

	if config.Values().Exists(vrfConfig, true) {
		return nil
	}

	_, err = arista.SendCommands("enable", "configure", vrfConfig)
	return err
}

func (arista *AristaEOS) RemoveVRFRoute(vrfName, interfaceName string, network *net.IPNet) error {
	config, err := arista.getRunningConfig()
	if err != nil {
		return fmt.Errorf("failed to get running config: %v", err)
	}

	if network == nil {
		return fmt.Errorf("ip network is nil")
	}

	vrfConfig := fmt.Sprintf("ip route vrf %s %s %s", vrfName, network.String(), interfaceName)
	if !config.Values().Exists(vrfConfig, true) {
		return nil
	}

	_, err = arista.SendCommands("enable", "configure", "no "+vrfConfig)
	return err
}
