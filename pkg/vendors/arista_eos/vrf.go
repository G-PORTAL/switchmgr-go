package arista_eos

import (
	"fmt"
	"net"
)

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
