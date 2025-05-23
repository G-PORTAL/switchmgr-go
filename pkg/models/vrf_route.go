package models

import "net"

type VRFRoute struct {
	Network       net.IPNet
	InterfaceName string
}
