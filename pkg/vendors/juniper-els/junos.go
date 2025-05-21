package juniper_els

import (
	"context"
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/config"
	"github.com/g-portal/switchmgr-go/pkg/vendors/registry"
	"github.com/g-portal/switchmgr-go/pkg/vendors/unimplemented"
	"github.com/neverlee/keymutex"
	"github.com/openshift-telco/go-netconf-client/netconf"
	"github.com/openshift-telco/go-netconf-client/netconf/message"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

// JuniperELS is a vendor implementation for Juniper EX switches running
// JunOS 15.1 or newer. It uses the ELS (Enhanced Layer 2 Software) configuration
// style and is not compatible with older JunOS versions. The driver is not
// fully implemented yet.
type JuniperELS struct {
	unimplemented.Unimplemented
	session    JuniperDriver
	identifier string
}

const Vendor registry.Vendor = "juniper_els"

// Timeout is the default timeout for netconf operations in seconds
const Timeout = 300

var configMutex = keymutex.New(128)

func (j *JuniperELS) Vendor() registry.Vendor {
	return Vendor
}

func (j *JuniperELS) Connect(cfg config.Connection) error {
	sshConfig := &ssh.ClientConfig{
		Config: ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr",
				"aes128-gcm@openssh.com",
				"arcfour256", "arcfour128", "arcfour",
				"aes128-cbc",
			},
		},
		Timeout: time.Second * 30,
		User:    cfg.Username,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Auth: []ssh.AuthMethod{ssh.Password(cfg.Password)},
	}

	session, err := netconf.NewSessionFromSSHConfigTimeout(context.Background(), fmt.Sprintf("%v:%d", cfg.Host, cfg.Port), sshConfig, time.Second*30)
	if err != nil {
		return fmt.Errorf("failed to connect to switch on ip %s: %s", cfg.Host, err.Error())
	}

	capabilities := netconf.DefaultCapabilities
	err = session.SendHello(&message.Hello{Capabilities: capabilities})
	if err != nil {
		return fmt.Errorf("failed to send hello to switch on ip %s: %s", cfg.Host, err.Error())
	}

	j.session = session
	return nil
}
func (j *JuniperELS) Disconnect() error {
	return j.session.Close()
}

func init() {
	registry.RegisterVendorFactory(Vendor, func() interface{} {
		return &JuniperELS{}
	})
}
