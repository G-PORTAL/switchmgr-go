package arista_eos

import (
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/iosconfig"
)

func (arista *AristaEOS) getRunningConfig() (iosconfig.Config, error) {
	output, err := arista.SendCommands("show running-config")
	if err != nil {
		return nil, fmt.Errorf("failed to show running-config: %w", err)
	}

	return iosconfig.Parse(output[0]), nil
}
