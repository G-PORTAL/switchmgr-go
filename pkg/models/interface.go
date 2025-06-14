package models

import (
	"bytes"
	"golang.org/x/exp/slices"
	textTemplate "text/template"
)

type InterfaceMode string

const (
	InterfaceModeAccess InterfaceMode = "access"
	InterfaceModeTrunk  InterfaceMode = "trunk"
)

// InterfaceType possible values are virtual, lag, 100base-tx, 1000base-t, 2.5gbase-t, 5gbase-t, 10gbase-t
// 10gbase-cx4, 1000base-x-gbic, 1000base-x-sfp, 10gbase-x-sfpp, 10gbase-x-xfp, 10gbase-x-xenpak, 10gbase-x-x2,
// 25gbase-x-sfp28, 50gbase-x-sfp56, 40gbase-x-qsfpp, 50gbase-x-sfp28, 100gbase-x-cfp, 100gbase-x-cfp2, 200gbase-x-cfp2,
// 100gbase-x-cfp4, 100gbase-x-cpak, 100gbase-x-qsfp28, 200gbase-x-qsfp56, 400gbase-x-qsfpdd, 400gbase-x-osfp,
// ieee802.11a, ieee802.11g, ieee802.11n, ieee802.11ac, ieee802.11ad, ieee802.11ax, gsm, cdma, lte, sonet-oc3,
// sonet-oc12, sonet-oc48, sonet-oc192, sonet-oc768, sonet-oc1920, sonet-oc3840, 1gfc-sfp, 2gfc-sfp, 4gfc-sfp,
// 8gfc-sfpp, 16gfc-sfpp, 32gfc-sfp28, 64gfc-qsfpp, 128gfc-sfp28, infiniband-sdr, infiniband-ddr, infiniband-qdr,
// infiniband-fdr10, infiniband-fdr, infiniband-edr, infiniband-hdr, infiniband-ndr, infiniband-xdr, t1, e1, t3, e3,
// cisco-stackwise, cisco-stackwise-plus, cisco-flexstack, cisco-flexstack-plus, juniper-vcp, extreme-summitstack,
// extreme-summitstack-128, extreme-summitstack-256, extreme-summitstack-512, other
type InterfaceType string

const (
	InterfaceTypeVirtual     InterfaceType = "virtual"
	InterfaceTypeLAG         InterfaceType = "lag"
	InterfaceType1GE         InterfaceType = "1000base-t"
	InterfaceType10GE        InterfaceType = "10gbase-t"
	InterfaceType10GSFPPlus  InterfaceType = "10gbase-x-sfpp"
	InterfaceType40QGSFPPlus InterfaceType = "40gbase-x-qsfpp"
)

// Valid returns true if the interface type is valid. This is used for validation
// of the interface type in the API. The SFP+ interface types are not valid.
func (t InterfaceType) Valid() bool {
	return t == InterfaceTypeVirtual ||
		t == InterfaceTypeLAG ||
		t == InterfaceType1GE ||
		t == InterfaceType10GE
}

type Interface struct {
	Name        string
	Description string
	Type        InterfaceType
	Enabled     bool

	MTU   uint32 // in bytes
	Speed uint32 // in Kbit/s

	Mode       InterfaceMode
	MacAddress MacAddress

	UntaggedVLAN *int32
	TaggedVLANs  []int32

	LagInterfaces []string
	IPAddresses   []string

	// The key is the original VLAN ID, the value is the new VLAN ID
	VlanMapping map[int32]int32

	Management    bool
	PortIsolation bool
}

type UpdateInterface struct {
	Name string

	Description *string
	Enabled     *bool

	MTU *uint32

	UntaggedVLAN *int32
	TaggedVLANs  []int32

	// The key is the original VLAN ID, the value is the new VLAN ID
	VlanMapping map[int32]int32

	// Change port isolation status
	PortIsolation *bool
}

// Disabled returns true if the interface should be disabled. This is the case if
// the Enabled field is set to false. If the Enabled field is nil, the interface
// should be .
func (u *UpdateInterface) Disabled() bool {
	return u.Enabled != nil && !*u.Enabled
}

// Template returns the a templated buffer for the interface, given by the template
// string. The template string should be a valid Go template. The template will
// be executed with the UpdateInterface as the data.
func (u *UpdateInterface) Template(template string) (bytes.Buffer, error) {
	var tpl bytes.Buffer

	tmpl, err := textTemplate.New("").Parse(template)
	if err != nil {
		return tpl, err
	}

	if err = tmpl.Execute(&tpl, u); err != nil {
		return tpl, err
	}

	return tpl, nil
}

func (u *UpdateInterface) Fill(i *Interface) {
	if u.Enabled == nil {
		u.Enabled = &i.Enabled
	}

	if u.MTU == nil {
		u.MTU = &i.MTU
	}
}

// Differs returns true if the interface differs from the update interface. Every
// field will be checked for differences, except for the name.
func (i *Interface) Differs(u *UpdateInterface) bool {
	if u.Description != nil && i.Description != *u.Description {
		return true
	}

	if u.Enabled != nil && i.Enabled != *u.Enabled {
		return true
	}

	if u.MTU != nil && i.MTU != *u.MTU {
		return true
	}

	if u.UntaggedVLAN != nil && (i.UntaggedVLAN == nil || *i.UntaggedVLAN != *u.UntaggedVLAN) {
		return true
	}

	if len(u.TaggedVLANs) != 0 && len(i.TaggedVLANs) != len(u.TaggedVLANs) {
		slices.Sort(i.TaggedVLANs)
		slices.Sort(u.TaggedVLANs)

		if !slices.Equal(i.TaggedVLANs, u.TaggedVLANs) {
			return true
		}
	}

	if len(u.VlanMapping) != 0 && len(i.VlanMapping) != len(u.VlanMapping) {
		return true
	}

	if len(u.VlanMapping) != 0 {
		for originalVlan, newVlan := range u.VlanMapping {
			if existingNewVlan, ok := i.VlanMapping[originalVlan]; ok && existingNewVlan != newVlan {
				return true
			}
		}
	}

	if len(u.VlanMapping) == 0 && (i.VlanMapping != nil && len(i.VlanMapping) > 0) {
		return true
	}

	if u.PortIsolation != nil && i.PortIsolation != *u.PortIsolation {
		return true
	}

	return false
}
