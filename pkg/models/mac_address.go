package models

import "net"

// MacAddress is a custom type for MAC addresses. It implements the Stringer
// interface, so it can be used as a string. It also implements the Valid()
// method, which checks if the MAC address is valid.
type MacAddress string

// String returns the MAC address as a string. First, we try to parse the MAC
// address, and if that fails, we return the original string.
func (m MacAddress) String() string {
	parsed, err := net.ParseMAC(string(m))
	if err != nil {
		return string(m)
	}

	return parsed.String()
}

// Valid checks if the MAC address is a valid MAC address.
func (m MacAddress) Valid() bool {
	_, err := net.ParseMAC(m.String())
	return err == nil
}
