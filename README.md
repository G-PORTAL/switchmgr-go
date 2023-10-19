# Switch Manager

This is a simple switch management library, written in Go. Currently, it supports only
Juniper, Juniper ELS and FScom switches, but can be extended easily. 

Extensive tests ensure that the library works as expected, but this is no guarantee. Software
updates for the switches can change the behavior of the switches and therefore break the library.

# Supported Vendors

- Juniper
- FibreStore

# Tested Switches

We tested a bunch of switches, but we can't guarantee that the library works with all switches
of the same type. If you have a switch that is not listed here, please test it and let us know
if it works or not.

- FibreStore S3900
- FibreStore S5800
- FibreStore N5860 (not all functions are supported yet)
- Juniper EX3200
- Juniper EX3300
- Juniper EX3400
- Juniper EX4200
- Juniper EX4550
- Juniper EX4600
- Juniper QFX510002
- Juniper QFX5100-48S-8C

# Usage

All vendors implementing the same interface. See vendors/drivers.go for the full
list of functions. The following example shows how to get the firmware version
as an example (the error handling is cut out of the example).

```go
package main

import (
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/config"
	"github.com/g-portal/switchmgr-go/pkg/vendors"
)

func main() {
	driver, err := vendors.New(vendors.VendorFSOSS5)
	driver.Connect(config.Connection{
		Host: "10.10.10.2",
		Port: 22,
		Username: "admin",
		Password: "admin",
	})
	defer driver.Disconnect()

	info, err := driver.GetHardwareInfo()
	if err != nil {
		panic(err)
	}

	fmt.Println(info.FirmwareVersion)
}

```

## Test

    go test ./... -count=1 -v

## Linting

    golangci-lint run
