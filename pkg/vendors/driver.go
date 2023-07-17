package vendors

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/g-portal/switchmgr-go/pkg/config"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fscom"
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper"
	"github.com/g-portal/switchmgr-go/pkg/vendors/juniper/els"
	"github.com/g-portal/switchmgr-go/pkg/vendors/unimplemented"
)

type Vendor string

const (
	VendorFiberStore Vendor = "fscom"
	VendorJuniper    Vendor = "juniper"
	VendorJuniperELS Vendor = "juniper_els"
)

func (v Vendor) Valid() bool {
	switch v {
	case VendorFiberStore, VendorJuniper, VendorJuniperELS:
		return true
	default:
		return false
	}
}

type Driver interface {
	Connect(cfg config.Connection) error
	Disconnect() error

	GetHardwareInfo() (*models.HardwareInfo, error)
	ListArpTable() ([]models.ArpEntry, error)

	ListInterfaces() ([]*models.Interface, error)
	ConfigureInterface(port *models.Interface) error
	GetInterface(name string) (*models.Interface, error)

	ListLLDPNeighbours() ([]models.LLDPNeighbour, error)

	Logger() *log.Logger
}

func New(vendor Vendor) (Driver, error) {
	if !vendor.Valid() {
		return nil, fmt.Errorf("unsupported vendor %s", vendor)
	}

	switch vendor {
	case VendorFiberStore:
		return &fscom.FSCom{}, nil
	case VendorJuniper:
		return &juniper.Juniper{}, nil
	case VendorJuniperELS:
		return &els.JuniperELS{}, nil
	default:
		return &unimplemented.Unimplemented{}, nil
	}
}
