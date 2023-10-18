package fsos_n5_test

import (
	"github.com/g-portal/switchmgr-go/pkg/models"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_n5"
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_n5/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
)

var expectedEntries = map[string][]models.MacAddress{
	"AggregatePort 1": {
		"00:c0:1d:c0:ff:ee",
		"98:5d:82:47:d0:93",
		"98:5d:82:47:d6:d9",
	},
	"AggregatePort 7": {
		"48:df:37:79:be:20",
	},
	"AggregatePort 8": {
		"48:df:37:79:be:90",
	},
}

func TestParseArpTable(t *testing.T) {
	entries, err := fsos_n5.ParseArpTable(utils.ReadTestData("show mac address table", nil))
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
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
