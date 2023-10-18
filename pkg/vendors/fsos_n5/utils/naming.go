package utils

import "strings"

// ConvertInterface converts the interface name to a consistent format
func ConvertInterface(name string) string {
	name = strings.TrimSpace(name)

	if strings.HasPrefix(name, "Hu") {
		return strings.Replace(name, "Hu", "HundredGigabitEthernet ", 1)
	} else if strings.HasPrefix(name, "Te") {
		return strings.Replace(name, "Te", "TenGigabitEthernet ", 1)
	}

	return name
}
