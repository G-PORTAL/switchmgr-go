package utils

import (
	"golang.org/x/exp/slices"
	"regexp"
	"strconv"
	"strings"
)

// ConvertVlans converts a string like "1-10" to a list of VLAN IDs.

var vlanRangeRgx = regexp.MustCompile("^[0-9]+\\-[0-9]+$")

func ConvertVlans(input string, seperator string) []int32 {
	vlanIDs := make([]int32, 0)
	if input == "" {
		return vlanIDs
	}

	for _, s := range strings.Split(input, seperator) {
		if !vlanRangeRgx.MatchString(s) {
			vlanID, err := strconv.Atoi(s)
			if err == nil {
				vlanIDs = append(vlanIDs, int32(vlanID))
			}

			continue
		}

		parts := strings.Split(s, "-")
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		end, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		n := start
		for n <= end {
			vlanIDs = append(vlanIDs, int32(n))
			n += 1
		}
	}

	// remove duplicates
	vlanIDs = uniqueVlanIDs(vlanIDs)

	// sort
	slices.Sort(vlanIDs)

	return vlanIDs
}

func DeleteVlanFromIDs(vlanIDs []int32, vlanID int32) []int32 {
	for i, id := range vlanIDs {
		if id == vlanID {
			return append(vlanIDs[:i], vlanIDs[i+1:]...)
		}
	}

	// sort
	slices.Sort(vlanIDs)

	return vlanIDs
}

// uniqueVlanIDs
func uniqueVlanIDs(intSlice []int32) []int32 {
	keys := make(map[int32]bool)
	list := []int32{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
