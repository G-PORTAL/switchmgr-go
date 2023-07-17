package iosconfig

import (
	"fmt"
	"strings"
)

// Interfaces returns a map of interfaces and their configuration
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

// Interface returns the configuration of a specific interface
func (cfg Config) Interface(name string) (ConfigValues, error) {
	for nic, config := range cfg.Interfaces() {
		if nic == name {
			return config, nil
		}
	}

	return nil, fmt.Errorf("interface %q not found", name)
}
