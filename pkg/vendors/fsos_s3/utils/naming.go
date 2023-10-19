package utils

import "strings"

// ConvertInterface converts the interface name to a consistent format
func ConvertInterface(name string) string {
	name = strings.TrimSpace(name)

	if strings.HasPrefix(name, "tg") {
		return strings.Replace(name, "tg", "TGigaEthernet", 1)
	} else if strings.HasPrefix(name, "TGi") {
		return strings.Replace(name, "TGi", "TGigaEthernet", 1)
	} else if strings.HasPrefix(name, "g") {
		return strings.Replace(name, "g", "GigaEthernet", 1)
	}

	return name
}
