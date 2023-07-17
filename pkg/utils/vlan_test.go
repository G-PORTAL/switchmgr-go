package utils_test

import (
	"github.com/g-portal/switchmgr-go/pkg/utils"
	"testing"
)

func TestConvertVlanRange(t *testing.T) {
	if len(utils.ConvertVlans("1-10,10", ",")) != 10 {
		t.Error("expected 10 vlans")
	}

	if len(utils.ConvertVlans("", ",")) != 0 {
		t.Error("expected 0 vlans")
	}

	if len(utils.ConvertVlans("abcdef", ",")) != 0 {
		t.Error("expected 0 vlans")
	}

	if len(utils.ConvertVlans("1,5-6,10", ",")) != 4 {
		t.Error("expected 4 vlans")
	}

	if len(utils.ConvertVlans("1,5-6,10,11-12", ",")) != 6 {
		t.Error("expected 6 vlans")
	}

	vlans := utils.ConvertVlans("1,5-6,10,11-12,2000,11", ",")
	if len(vlans) != 7 {
		t.Error("expected 7 vlans")
	}

	if vlans[0] != 1 {
		t.Error("expected 1 as first vlan")
	}

	if vlans[1] != 5 {
		t.Error("expected 5 as second vlan")
	}

	if vlans[2] != 6 {
		t.Error("expected 6 as third vlan")
	}

	if vlans[3] != 10 {
		t.Error("expected 10 as fourth vlan")
	}

	if vlans[4] != 11 {
		t.Error("expected 11 as fifth vlan")
	}

	if vlans[5] != 12 {
		t.Error("expected 12 as sixth vlan")
	}

	if vlans[6] != 2000 {
		t.Error("expected 2000 as seventh vlan")
	}

	// Test with different delimiter
	if len(utils.ConvertVlans("1;5-6;10;11-12;2000;11", ";")) != 7 {
		t.Error("expected 7 vlans")
	}

	// Test with duplicated vlans
	if len(utils.ConvertVlans("1;5-6;10;11-12;2000;11;1", ";")) != 7 {
		t.Error("expected 7 vlans")
	}
}
