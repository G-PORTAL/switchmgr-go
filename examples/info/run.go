package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/g-portal/switchmgr-go/pkg/config"
	"github.com/g-portal/switchmgr-go/pkg/vendors"
	"github.com/g-portal/switchmgr-go/pkg/vendors/registry"
	"os"
)

func main() {
	driver, err := vendors.New(registry.Vendor(os.Getenv("SWITCH_VENDOR")))
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}

	if err := driver.Connect(config.Connection{
		Host:     os.Getenv("SWITCH_HOST"),
		Port:     22,
		Username: os.Getenv("SWITCH_USERNAME"),
		Password: os.Getenv("SWITCH_PASSWORD"),
	}); err != nil {
		log.Errorf(fmt.Sprintf("Failed to connect to switch: %v", err))
	}
	defer func() {
		if err := driver.Disconnect(); err != nil {
			log.Errorf(fmt.Sprintf("Failed to disconnect from switch: %v", err))
		}
	}()

	driver.Logger().SetLevel(log.InfoLevel)
	driver.Logger().Info("+++++++++++++++++++")
	driver.Logger().Info("System Information:")
	driver.Logger().Info("+++++++++++++++++++")
	if info, err := driver.GetHardwareInfo(); err == nil {
		driver.Logger().Infof("Hostname: %s", info.Hostname)
		driver.Logger().Infof("Serial: %s", info.Serial)
		driver.Logger().Infof("Model: %s", info.Model)
		driver.Logger().Infof("Firmware Version: %s", info.FirmwareVersion)
	}

	driver.Logger().Info("+++++++++++++++++++")
	driver.Logger().Info("ARP Table:")
	driver.Logger().Info("+++++++++++++++++++")
	if arp, err := driver.ListArpTable(); err == nil {
		for _, entry := range arp {

			var macs []string
			for _, mac := range entry.MacAddresses {
				macs = append(macs, mac.String())
			}

			driver.Logger().Infof("Interface %q: %q", entry.SwitchPort, macs)
		}
	}

	driver.Logger().Info("+++++++++++++++++++")
	driver.Logger().Info("LLDP Table:")
	driver.Logger().Info("+++++++++++++++++++")
	if lldp, err := driver.ListLLDPNeighbors(); err == nil {
		for _, neighbor := range lldp {
			driver.Logger().Infof("Interface: %q: %s", neighbor.LocalInterface, neighbor.RemoteHostname)
		}
	}

	driver.Logger().Info("+++++++++++++++++++")
	driver.Logger().Info("Interfaces:")
	driver.Logger().Info("+++++++++++++++++++")
	nics, err := driver.ListInterfaces()
	if err != nil {
		driver.Logger().Errorf("Failed to get interfaces: %v", err)
	} else {
		for _, nic := range nics {
			untagged := "-"
			if nic.UntaggedVLAN != nil {
				untagged = fmt.Sprintf("%d", *nic.UntaggedVLAN)
			}
			driver.Logger().Infof("Interface: %q (%s): %s (untagged: %s, tagged: %+v, enabled: %v)",
				nic.Name, nic.Description, nic.MacAddress, untagged, nic.TaggedVLANs, nic.Enabled)
		}
	}

}
