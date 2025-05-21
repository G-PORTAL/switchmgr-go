package fsos_n5

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s3"
	"github.com/g-portal/switchmgr-go/pkg/vendors/registry"
)

const Vendor registry.Vendor = "fsos_n5"

type FSComN5 struct {
	fsos_s3.FSComS3
}

func (fs *FSComN5) Vendor() registry.Vendor {
	return Vendor
}

func init() {
	registry.RegisterVendorFactory(Vendor, func() interface{} {
		return &FSComN5{
			FSComS3: fsos_s3.FSComS3{
				LoginCommands: []string{
					"terminal length 0",
				},
			},
		}
	})
}
