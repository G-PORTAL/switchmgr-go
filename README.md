# Switch Manager

This is a simple switch management library, written in Go. Currently, it supports only
Juniper, Juniper ELS and FScom switches, but can be extended easily. 

Extensive tests ensure that the library works as expected, but this is no guarantee. Software
updates for the switches can change the behavior of the switches and therefore break the library. 

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
	driver, err := vendors.New(vendors.VendorFiberStore)
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
