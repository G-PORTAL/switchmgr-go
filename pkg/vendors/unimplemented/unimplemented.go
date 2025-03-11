package unimplemented

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/g-portal/switchmgr-go/pkg/config"
	"github.com/g-portal/switchmgr-go/pkg/models"
	"net"
	"os"
	"time"
)

type Unimplemented struct {
	logger *log.Logger
}

func (j *Unimplemented) Connect(cfg config.Connection) error {
	return errors.New("not implemented")
}

func (j *Unimplemented) Disconnect() error {
	return errors.New("not implemented")
}

func (j *Unimplemented) GetHardwareInfo() (*models.HardwareInfo, error) {
	return nil, errors.New("not implemented")
}

func (j *Unimplemented) ListArpTable() ([]models.ArpEntry, error) {
	return nil, errors.New("not implemented")
}

func (j *Unimplemented) ListInterfaces() ([]*models.Interface, error) {
	return nil, errors.New("not implemented")
}
func (j *Unimplemented) AddVRFRoute(vrfName, interfaceName string, network *net.IPNet) error {
	return errors.New("not implemented")
}
func (j *Unimplemented) RemoveVRFRoute(vrfName, interfaceName string, network *net.IPNet) error {
	return errors.New("not implemented")
}
func (j *Unimplemented) ConfigureInterface(update *models.UpdateInterface) (bool, error) {
	return false, errors.New("not implemented")
}
func (j *Unimplemented) GetInterface(name string) (*models.Interface, error) {
	return nil, errors.New("not implemented")
}

func (j *Unimplemented) ListLLDPNeighbors() ([]models.LLDPNeighbor, error) {
	return nil, nil
}

func (j *Unimplemented) Logger() *log.Logger {
	if j.logger != nil {
		return j.logger
	}

	j.logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		Level:           log.DebugLevel,
		Formatter:       log.TextFormatter,
		TimeFormat:      time.RFC3339,
		Prefix:          "[switchmgr]",
	})

	return j.logger
}
