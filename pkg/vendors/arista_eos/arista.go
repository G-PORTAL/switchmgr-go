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

type AristaEOS struct {
	unimplemented.Unimplemented

	LoginCommands []string

	connection goeapi.Node
	mu         sync.Mutex
}

func (arista *AristaEOS) Vendor() registry.Vendor {
	return Vendor
}

// Connect Connecting to a FiberStore switch using SSH
func (arista *AristaEOS) Connect(cfg config.Connection) error {
	arista.Logger().SetPrefix(fmt.Sprintf("[arista/%s]", cfg.Host))

	var err error
	connection, err := goeapi.Connect("https", cfg.Host, cfg.Username, cfg.Password, 443)
	if err != nil {
		return fmt.Errorf("failed to connect to eapi: %w", err)
	}

	arista.connection = *connection

	return nil
}

func (arista *AristaEOS) Disconnect() error {
	return nil
}

// Save  saves the configuration to startup config
func (arista *AristaEOS) Save() error {
	output, err := arista.SendCommands("write memory")
	if err != nil {
		return err
	}

	if !strings.Contains(strings.TrimSpace(output[0]), "successfully") {
		return fmt.Errorf("failed to save the configuration: %s", output)
	}

	return err
}

// SendCommands sends a list of commands to the switch and returns the output
func (arista *AristaEOS) SendCommands(commands ...string) ([]string, error) {
	arista.mu.Lock()
	defer arista.mu.Unlock()

	response, err := arista.connection.RunCommands(commands, "text")
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
	registry.RegisterVendor(Vendor, &AristaEOS{
		LoginCommands: []string{
			"terminal length 0",
		},
	})
}
