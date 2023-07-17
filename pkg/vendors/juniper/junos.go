package juniper

import (
	"fmt"
	"github.com/Juniper/go-netconf/netconf"
	"github.com/g-portal/switchmgr-go/pkg/config"
	"github.com/g-portal/switchmgr-go/pkg/vendors/unimplemented"
	"github.com/neverlee/keymutex"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

type Juniper struct {
	unimplemented.Unimplemented
	vlanMapping    map[string]int32
	interfaceVlans map[string]junosVlanMapEntry
	session        JuniperDriver
	identifier     string
}

var configMutex = keymutex.New(128)

func (j *Juniper) Connect(cfg config.Connection) error {
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
	j.identifier = fmt.Sprint(cfg.Host, ":", cfg.Port)

	session, err := netconf.DialSSHTimeout(fmt.Sprintf("%v:%d", cfg.Host, cfg.Port), sshConfig, time.Second*30)
	if err != nil {
		return fmt.Errorf("failed to connect to switch on ip %s: %s", cfg.Host, err.Error())
	}

	j.session = session
	return nil
}
func (j *Juniper) Disconnect() error {
	return j.session.Close()
}
