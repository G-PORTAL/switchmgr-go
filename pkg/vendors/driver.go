package vendors

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/g-portal/switchmgr-go/pkg/config"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/registry"
	"os"
	"plugin"
	"strings"

	// load all native vendor drivers
	_ "github.com/g-portal/switchmgr-go/pkg/vendors/fsos_n5"
	_ "github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s3"
	_ "github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s5"
	_ "github.com/g-portal/switchmgr-go/pkg/vendors/juniper"
	_ "github.com/g-portal/switchmgr-go/pkg/vendors/juniper-els"
)

const pluginsVar = "SWITCHMGR_PLUGINS"

// Driver is the interface, which all vendor drivers have to implement. The
// driver is responsible for connecting to the switch, getting information
// about the switch and configuring the switch. The driver can also overwrite
// the logger. Some vendors do not support all features, so some functions
// may return the "not implemented" error.
type Driver interface {
	// Vendor returns the vendor and the driver instance. The vendor is used
	// to identify the driver and the driver instance is used to call the
	// functions of the driver.
	Vendor() registry.Vendor
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
}

// New returns a new driver for the given vendor. If the vendor is not
// supported, an error is returned.
func New(vendor registry.Vendor) (Driver, error) {
	// load plugins
	err := loadPlugins()
	if err != nil {
		return nil, fmt.Errorf("failed to load plugins: %v", err)
	}

	if !vendor.Valid() {
		return nil, fmt.Errorf("unsupported vendor %s", vendor)
	}

	// Return the driver instance for the given vendor.
	implementation, err := registry.GetVendor(vendor)
	if err != nil {
		return nil, err
	}

	if driver, ok := implementation.(Driver); ok {
		return driver, nil
	}

	return nil, fmt.Errorf("vendor %s does not implement the driver interface", vendor)
}

func loadPlugins() error {
	pluginPaths := os.Getenv(pluginsVar)
	// we have no plugins, so we allow the program to continue
	if pluginPaths == "" {
		return nil
	}

	for _, path := range strings.Split(pluginPaths, ",") {
		p, err := plugin.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open plugin %s: %v", path, err)
		}

		sym, err := p.Lookup("Driver")
		if err != nil {
			return fmt.Errorf("failed to lookup Driver symbol in plugin %s: %v", path, err)
		}

		// check if the symbol is implementing the driver interface
		if _, ok := sym.(*Driver); !ok {
			return fmt.Errorf("the plugin %s does not implement the driver interface", path)
		}

		driver := *sym.(*Driver)

		registry.RegisterVendor(driver.Vendor(), driver)
	}

	return nil
}
