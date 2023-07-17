package models

type HardwareInfo struct {
	Hostname string
	Vendor   string
	Model    string

	FirmwareVersion string
	Serial          string
}
