package fsos_test

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
)

var expectedEntries = map[string][]models.MacAddress{
	"eth-0-15": {
		"5c:ba:2c:3a:1c:17",
	},
	"eth-0-14": {
		"5c:ba:2c:39:17:2b",
	},
	"eth-0-17": {
		"9c:6b:00:1d:74:d4",
	},
	"eth-0-16": {
		"50:7c:6f:16:f9:26",
		"5c:ed:8c:31:0a:e7",
	},
	"agg1": {
		"64:9d:99:0a:7a:cb",
		"64:9d:99:06:11:45",
		"00:de:ad:be:ef:02",
		"00:de:ad:be:ef:01",
		"00:c0:1d:c0:ff:ee",
	},
}

func TestParseArpTable(t *testing.T) {
	entries, err := fsos.ParseArpTable(utils.ReadTestData("show mac address table", nil))
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(entries) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(entries))
	}

	for _, entry := range entries {
		expectedMacs, ok := expectedEntries[entry.SwitchPort]
		if !ok {
			t.Fatalf("unexpected port %s", entry.SwitchPort)
		}

		if len(entry.MacAddresses) != len(expectedMacs) {
			t.Fatalf("invalid amount of macs for port %q, got %v, expected %v", entry.SwitchPort,
				len(entry.MacAddresses), len(expectedMacs))
		}

		if diff := cmp.Diff(entry.MacAddresses, expectedMacs, cmpopts.SortSlices(func(a, b string) bool { return a < b })); diff != "" {
			t.Fatalf("mac addresses for port %q don't match. diff: %s", entry.SwitchPort, diff)
		}

	}

}
