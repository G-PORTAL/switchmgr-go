package models

type Interface struct {
	Name        string
	Description string
	Enabled     bool

	MacAddress MacAddress

	UntaggedVLAN *int32
	TaggedVLANs  []int32
}
