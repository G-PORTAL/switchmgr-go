package fsos_s5_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_n5/utils"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s5"
	"testing"
)

func TestListVlanMappings(t *testing.T) {
	iosConfig := fsos_s5.ParseConfiguration(utils.ReadTestData("show running-config", nil))
	mappings := iosConfig.VlanMappings()

	expectedMappings := map[string]map[int32]int32{
		"custom_mapping": {
			10: 100,
			15: 150,
		},
	}
	for mapName := range expectedMappings {
		if mapping, ok := mappings[mapName]; !ok {
			t.Fatalf("Expected %s mapping to exist", mapName)
		} else {
			for vlan, expected := range expectedMappings[mapName] {
				if _, ok2 := mapping[vlan]; !ok2 || mapping[vlan] != expected {
					t.Fatalf("Expected mapping %d => %d to exist", vlan, expected)
				}
			}
		}
	}
}
