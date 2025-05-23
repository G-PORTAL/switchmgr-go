package arista_eos

import (
	"encoding/json"
	"fmt"
	"github.com/aristanetworks/goeapi"
	"github.com/g-portal/switchmgr-go/pkg/config"
	"github.com/g-portal/switchmgr-go/pkg/vendors/registry"
	"github.com/g-portal/switchmgr-go/pkg/vendors/unimplemented"
	"strings"
	"sync"
)

const Vendor registry.Vendor = "arista_eos"

var aristaMutex sync.Mutex

// AristaEOS Please note that the goeapi library is not thread-safe
// See also: https://github.com/aristanetworks/goeapi/issues/73
type AristaEOS struct {
	unimplemented.Unimplemented

	LoginCommands []string

	connection goeapi.Node
}

func (arista *AristaEOS) Vendor() registry.Vendor {
	return Vendor
}

// Connect Connecting to a FiberStore switch using SSH
func (arista *AristaEOS) Connect(cfg config.Connection) error {
	arista.Logger().SetPrefix(fmt.Sprintf("[arista/%s]", cfg.Host))

	aristaMutex.Lock()
	defer aristaMutex.Unlock()

	var err error
	connection, err := goeapi.Connect("https", cfg.Host, cfg.Username, cfg.Password, 443)
	if err != nil {
		return fmt.Errorf("failed to connect to eapi: %w", err)
	}

	arista.connection = *connection
	arista.connection.GetConnection().SetDisableKeepAlive(true)
	arista.connection.GetConnection().SetTimeout(60)

	return nil
}

func (arista *AristaEOS) Disconnect() error {
	return nil
}

// Save  saves the configuration to startup config
func (arista *AristaEOS) Save() error {
	output, err := arista.SendCommands("enable", "write memory")
	if err != nil {
		return err
	}

	if !strings.Contains(strings.TrimSpace(output[1]), "successfully") {
		return fmt.Errorf("failed to save the configuration: %s", output)
	}

	return err
}

// cmdsToInterface, copied from https://github.com/aristanetworks/goeapi/blob/7090068b8735dc15c22444cbda080db4052ae8af/client.go#L440
func (arista *AristaEOS) cmdsToInterface(commands []string) []interface{} {
	if commands == nil || len(commands) == 0 {
		return nil
	}
	var interfaceSlice []interface{}
	length := len(commands)

	interfaceSlice = make([]interface{}, length)

	for i := 0; i < length; i++ {
		interfaceSlice[i] = commands[i]
	}
	return interfaceSlice
}

// SendCommands sends a list of commands to the switch and returns the output
func (arista *AristaEOS) SendCommands(commands ...string) ([]string, error) {
	aristaMutex.Lock()
	defer aristaMutex.Unlock()

	// we use arista.connection.RunCommands but it always adds "enable" before all commands, even for
	// those where we don't need those privileges.
	response, err := arista.connection.GetConnection().Execute(arista.cmdsToInterface(commands), "text")
	if err != nil {
		return nil, err
	}

	output := make([]string, 0)
	for i, result := range response.Result {
		commandOutput, ok := result["output"]
		if !ok {
			arista.Logger().Debugf("Output of command %q is missing", commands[i])
			continue
		}

		stringOutput, ok := commandOutput.(string)
		if !ok {
			arista.Logger().Debugf("Output of command %q is from invalid type: %T", commands[i], commandOutput)
			continue
		}

		output = append(output, stringOutput)
	}

	return output, nil
}

// GetJsonResponse Runs a command and returns the json as bytes
func (arista *AristaEOS) GetJsonResponse(resp interface{}, command string) error {
	response, err := arista.SendCommands(fmt.Sprintf("%s | json", command))
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(response[0]), resp)
}

func init() {
	registry.RegisterVendorFactory(Vendor, func() interface{} {
		return &AristaEOS{
			LoginCommands: []string{
				"terminal length 0",
			},
		}
	})
}
