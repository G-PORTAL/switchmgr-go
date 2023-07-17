package juniper

import (
	"regexp"
)

var interfaceNameRegex = regexp.MustCompile(`^[A-Za-z]+-[0-9]+\/[0-9]+\/[0-9]+$`)

func isValidInterface(iface string) bool {
	return interfaceNameRegex.MatchString(iface)
}
