package models

import "net"

type MacAddress string

func (m MacAddress) String() string {
	parsed, err := net.ParseMAC(string(m))
	if err != nil {
		return string(m)
	}

	return parsed.String()
}

func (m MacAddress) Valid() bool {
	_, err := net.ParseMAC(m.String())
	if err != nil {
		return false
	}

	return true
}
