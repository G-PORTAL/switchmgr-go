package registry

import (
	"fmt"
	"sync"
)

// Vendor is the type for the vendor name. It is used to identify the vendor
// and to load the correct driver.
type Vendor string

var registeredVendors map[Vendor]interface{}
var registeredVendorsMtx sync.RWMutex

// Valid checks if this lib supports the given vendor.
func (v Vendor) Valid() bool {
	registeredVendorsMtx.RLock()
	defer registeredVendorsMtx.RUnlock()

	// we have no implementation
	if _, ok := registeredVendors[v]; !ok {
		return false
	}

	return true
}

// GetVendor returns the vendor implementation
func GetVendor(v Vendor) (interface{}, error) {
	registeredVendorsMtx.RLock()
	defer registeredVendorsMtx.RUnlock()

	// we have no implementation
	if implementation, ok := registeredVendors[v]; ok {
		return implementation, nil
	}

	return nil, fmt.Errorf("failed to get implementation for vendor %s", v)
}

// RegisterVendor registers a vendor plugin
func RegisterVendor(v Vendor, implementation interface{}) {
	registeredVendorsMtx.Lock()
	defer registeredVendorsMtx.Unlock()

	if registeredVendors == nil {
		registeredVendors = make(map[Vendor]interface{})
	}

	registeredVendors[v] = implementation
}
