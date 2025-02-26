package main

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors"
	"github.com/g-portal/switchmgr-go/pkg/vendors/registry"
	"github.com/g-portal/switchmgr-go/pkg/vendors/unimplemented"
)

const Vendor registry.Vendor = "example"

type ExampleDriver struct {
	unimplemented.Unimplemented
}

func (d *ExampleDriver) Vendor() registry.Vendor {
	return Vendor
}

// Driver Exported symbol
var Driver vendors.Driver = &ExampleDriver{}
