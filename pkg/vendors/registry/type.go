package registry

import (
	"fmt"
	"sync"
)

// Vendor is the type for the vendor name. It is used to identify the vendor
// and to load the correct driver.
type Vendor string

var vendorFactories = make(map[Vendor]func() interface{})
var registeredVendorsMtx sync.RWMutex

// Valid checks if this lib supports the given vendor.
func (v Vendor) Valid() bool {
	registeredVendorsMtx.RLock()
	defer registeredVendorsMtx.RUnlock()

	_, ok := vendorFactories[v]

	return ok
}

// GetVendor returns the vendor implementation
func GetVendor(vendor Vendor) (interface{}, error) {
	registeredVendorsMtx.RLock()
	defer registeredVendorsMtx.RUnlock()

	if factory, ok := vendorFactories[vendor]; ok {
		return factory(), nil
	}

	return nil, fmt.Errorf("vendor %s not found", vendor)
}

func RegisterVendorFactory(vendor Vendor, factory func() interface{}) {
	registeredVendorsMtx.Lock()
	defer registeredVendorsMtx.Unlock()

	vendorFactories[vendor] = factory
}
