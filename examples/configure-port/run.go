package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/g-portal/switchmgr-go/pkg/config"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors"
	"os"
	"strconv"
	"strings"
)

func main() {
	driver, err := vendors.New(vendors.Vendor(os.Getenv("SWITCH_VENDOR")))
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

	driver.Logger().SetLevel(log.DebugLevel)

	name := os.Getenv("PORT_NAME")
	if name == "" {
		log.Errorf("PORT_NAME is not set")
	}

	description := os.Getenv("PORT_DESCRIPTION")
	t := true

	update := &models.UpdateInterface{
		Name:         name,
		Enabled:      &t,
		UntaggedVLAN: nil,
		TaggedVLANs:  nil,
	}

	if description != "" {
		update.Description = &description
	}

	untaggedVLANString := os.Getenv("PORT_UNTAGGED_VLAN_ID")
	if untaggedVLANString != "" {
		vlanID, err := strconv.Atoi(untaggedVLANString)
		if err != nil {
			log.Errorf(fmt.Sprintf("Failed to parse untagged VLAN ID: %v", err))
		}

		untaggedVLAN := int32(vlanID)
		update.UntaggedVLAN = &untaggedVLAN
	}

	taggedVLANIds := os.Getenv("PORT_TAGGED_VLAN_IDS")
	if taggedVLANIds != "" {
		var taggedVLANs []int32
		for _, vlanIDString := range strings.Split(taggedVLANIds, ",") {
			vlanID, err := strconv.Atoi(vlanIDString)
			if err != nil {
				log.Errorf(fmt.Sprintf("Failed to parse tagged VLAN ID: %v", err))
			}

			taggedVLANs = append(taggedVLANs, int32(vlanID))
		}

		update.TaggedVLANs = taggedVLANs
	}

	changed, err := driver.ConfigureInterface(update)
	if err != nil {
		log.Errorf(fmt.Sprintf("Failed to configure interface: %v", err))
	}

	if changed {
		log.Infof("Interface %s has been configured", name)
	} else {
		log.Infof("Interface %s is already configured", name)
	}
}
