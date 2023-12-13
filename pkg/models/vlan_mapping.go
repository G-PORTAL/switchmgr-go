package models

type VlanMapping struct {
	GroupName string

	Mapping map[int32]int32
}
