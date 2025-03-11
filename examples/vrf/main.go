package main

import (
	"github.com/charmbracelet/log"
	"github.com/g-portal/switchmgr-go/pkg/config"
	"github.com/g-portal/switchmgr-go/pkg/vendors"
	"github.com/g-portal/switchmgr-go/pkg/vendors/registry"
	"net"
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
		log.Errorf("Failed to connect to switch: %v", err)
	}
	defer func() {
		if err := driver.Disconnect(); err != nil {
			log.Errorf("Failed to disconnect from switch: %v", err)
		}
	}()

	driver.Logger().SetLevel(log.DebugLevel)

	vrfName := os.Getenv("VRF_NAME")
	if vrfName == "" {
		log.Errorf("VRF_NAME is not set")
	}

	interfaceName := os.Getenv("INTERFACE_NAME")
	if interfaceName == "" {
		log.Errorf("INTERFACE_NAME is not set")
	}

	ipAddress := os.Getenv("IP_ADDRESS")
	if ipAddress == "" {
		log.Errorf("IP_ADDRESS is not set")
	}

	_, network, err := net.ParseCIDR(ipAddress)
	if err != nil {
		log.Errorf("Failed to parse IP address: %v", err)
	}

	mode := os.Getenv("MODE")
	switch mode {
	case "remove":
		err := driver.RemoveVRFRoute(vrfName, interfaceName, network)
		if err != nil {
			log.Errorf("Failed to remove VRF route: %v", err)
		}

	default:
		err := driver.AddVRFRoute(vrfName, interfaceName, network)
		if err != nil {
			log.Errorf("Failed to add VRF route: %v", err)
		}
	}
}
