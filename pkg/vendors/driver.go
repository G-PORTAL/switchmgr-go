package vendors

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/g-portal/switchmgr-go/pkg/config"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_n5"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s3"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s5"
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper"
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper-els"
	"github.com/g-portal/switchmgr-go/pkg/vendors/unimplemented"
)

// Vendor is the type for the vendor name. It is used to identify the vendor
// and to load the correct driver.
type Vendor string

const (
	// VendorFSOSS3 FSComS3 (FiberStore), series S3XX
	VendorFSOSS3 Vendor = "fsos_s3"
	// VendorFSOSS5 FSComS3 (FiberStore), series S5XX
	VendorFSOSS5 Vendor = "fsos_s5"
	// VendorFSOSN5 FSComS3 (FiberStore), series N5XX
	VendorFSOSN5 Vendor = "fsos_n5"
	// VendorJuniper Juniper, up to version 15 (legacy)
	VendorJuniper Vendor = "juniper"
	// VendorJuniperELS Juniper, version 15.1 and higher with advanced
	// ELS (enhanced layer 2 software) features.
	VendorJuniperELS Vendor = "juniper_els"
)

// Valid checks if this lib supports the given vendor.
func (v Vendor) Valid() bool {
	switch v {
	case VendorFSOSS3, VendorFSOSS5, VendorFSOSN5, VendorJuniper, VendorJuniperELS:
		return true
	default:
		return false
	}
}

// Driver is the interface, which all vendor drivers have to implement. The
// driver is responsible for connecting to the switch, getting information
// about the switch and configuring the switch. The driver can also overwrite
// the logger. Some vendors do not support all features, so some functions
// may return the "not implemented" error.
type Driver interface {
	// Connect connects to the switch via SSH. The connection information is
	// provided via the first parameter. Keep in mind that the password is
	// not encrypted.
	Connect(cfg config.Connection) error
	// Disconnect disconnects from the switch.
	Disconnect() error

	// GetHardwareInfo returns the hardware information of the switch. The
	// HardwareInfo struct may contain not all information, depending on the
	// vendor.
	GetHardwareInfo() (*models.HardwareInfo, error)
	// ListArpTable returns the ARP table of the switch. The ARP table contains
	// all MAC addresses and the switch port they are connected to.
	ListArpTable() ([]models.ArpEntry, error)

	// ListInterfaces returns a list of all interfaces, build in the switch. The
	// interface list may contain not all information, depending on the vendor.
	ListInterfaces() ([]*models.Interface, error)
	// ConfigureInterface configures the interface with the given configuration.
	// Returns true if the interface was changed and false if not.
	ConfigureInterface(update *models.UpdateInterface) (bool, error)
	// GetInterface returns the interface with the given name. The interface
	// may contain not all information, depending on the vendor.
	GetInterface(name string) (*models.Interface, error)
	// ListLLDPNeighbors returns a list of neighbors, which are discovered
	// through LLDP. This information can be used to get insights about the
	// network topology or for debugging purposes.
	ListLLDPNeighbors() ([]models.LLDPNeighbor, error)

	// Logger returns the logger of the driver. This is a generic logger, which
	// can be overwritten by the vendor driver.
	Logger() *log.Logger

	ListVlans() ([]models.Vlan, error)
	ConfigureVlan(vlan *models.Vlan) (bool, error)
	DeleteVlan(id int32) error

	ListVlanMappings() ([]models.VlanMapping, error)
	ConfigureVlanMapping(mapping *models.VlanMapping) (bool, error)
	DeleteVLANMapping(name string)
}

// New returns a new driver for the given vendor. If the vendor is not
// supported, an error is returned.
func New(vendor Vendor) (Driver, error) {
	if !vendor.Valid() {
		return nil, fmt.Errorf("unsupported vendor %s", vendor)
	}

	switch vendor {
	case VendorJuniper:
		return &juniper.Juniper{}, nil
	case VendorJuniperELS:
		return &juniper_els.JuniperELS{}, nil
	case VendorFSOSS3:
		return &fsos_s3.FSComS3{
			LoginCommands: []string{
				"enter", "terminal length 0",
			},
		}, nil
	case VendorFSOSS5:
		return &fsos_s5.FSComS5{
			FSComS3: fsos_s3.FSComS3{
				LoginCommands: []string{
					"terminal length 0",
				},
			},
		}, nil
	case VendorFSOSN5:
		return &fsos_n5.FSComN5{
			FSComS3: fsos_s3.FSComS3{
				LoginCommands: []string{
					"terminal length 0",
				},
			},
		}, nil
	default:
		// Should never be reached because of the Valid() check above.
		return &unimplemented.Unimplemented{}, nil
	}
}
