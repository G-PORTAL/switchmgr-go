package fsos_s5

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fsos_s3"
	"github.com/g-portal/switchmgr-go/pkg/vendors/registry"
)

const Vendor registry.Vendor = "fsos_s5"

type FSComS5 struct {
	fsos_s3.FSComS3
}

func (fs *FSComS5) Vendor() registry.Vendor {
	return Vendor
}

func init() {
	registry.RegisterVendor(Vendor, &FSComS5{
		FSComS3: fsos_s3.FSComS3{
			LoginCommands: []string{
				"terminal length 0",
			},
		},
	})
}
