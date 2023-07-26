package iosconfig

import (
	"fmt"
	"strings"
)

// Interfaces return a map of interfaces and their configuration. If there are
// no interfaces, an empty map is returned. An interface starts with the
// "interface" prefix.
func (cfg Config) Interfaces() map[string]ConfigValues {
	interfaces := make(map[string]ConfigValues)
	for key, values := range cfg {
		if strings.HasPrefix(key, "interface ") {
			nicName := strings.TrimSpace(strings.TrimPrefix(key, "interface "))
			interfaces[nicName] = values
		}
	}

	return interfaces
}

// Interface returns the configuration of a specific interface. If the interface
// is not found, an error is returned.
func (cfg Config) Interface(name string) (ConfigValues, error) {
	for nic, config := range cfg.Interfaces() {
		if nic == name {
			return config, nil
		}
	}

	return nil, fmt.Errorf("interface %q not found", name)
}
