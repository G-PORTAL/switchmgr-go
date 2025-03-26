package arista_eos

import (
	"encoding/json"
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/config"
	"github.com/g-portal/switchmgr-go/pkg/utils"
	"github.com/g-portal/switchmgr-go/pkg/vendors/registry"
	"github.com/g-portal/switchmgr-go/pkg/vendors/unimplemented"
	"golang.org/x/crypto/ssh"
	"io"
	"strings"
	"time"
)

const Vendor registry.Vendor = "arista_eos"

type AristaEOS struct {
	unimplemented.Unimplemented

	LoginCommands []string

	conn    *ssh.Client
	session *ssh.Session

	writer    io.WriteCloser
	reader    io.Reader
	errReader io.Reader
}

func (arista *AristaEOS) Vendor() registry.Vendor {
	return Vendor
}

// Connect Connecting to a FiberStore switch using SSH
func (arista *AristaEOS) Connect(cfg config.Connection) error {
	arista.Logger().SetPrefix(fmt.Sprintf("[arista/%s]", cfg.Host))

	var err error
	sshConfig := &ssh.ClientConfig{
		User:    cfg.Username,
		Timeout: 30 * time.Second,
		Auth: []ssh.AuthMethod{
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
				if len(questions) == 0 {
					return []string{}, nil
				}

				return []string{cfg.Password}, err
			}),
		},
		BannerCallback: func(message string) error {
			return nil
		},
		HostKeyAlgorithms: []string{ssh.KeyAlgoDSA, ssh.KeyAlgoRSA, ssh.KeyAlgoECDSA256, ssh.KeyAlgoED25519},
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
	}

	if cfg.Password != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.RetryableAuthMethod(ssh.Password(strings.TrimSpace(cfg.Password)), 3))
	}

	arista.conn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%v", cfg.Host, cfg.Port), sshConfig)
	if err != nil {
		return err
	}

	if arista.session, err = arista.conn.NewSession(); err != nil {
		return err
	}

	// Set up terminal modes
	arista.Logger().Debugf("Requesting pseudo terminal.")
	if err = arista.session.RequestPty("XTERM", 5000, 5000, ssh.TerminalModes{}); err != nil {
		return fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	// create pipe to the stdin of the SSH process
	arista.writer, err = arista.session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %s", err)
	}

	// create pipe to the stdout of the SSH process
	arista.reader, err = arista.session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %s", err)
	}

	// create pipe to the stderr of the SSH process
	arista.errReader, err = arista.session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %s", err)
	}

	// start the shell
	arista.Logger().Debugf("Starting the session.")
	if err = arista.session.Shell(); err != nil {
		return fmt.Errorf("failed to start the shell: %s", err)
	}

	// entering the shell
	if _, err = arista.SendCommands(arista.LoginCommands...); err != nil {
		return err
	}

	return nil
}

func (arista *AristaEOS) Disconnect() error {
	if arista.session != nil {
		if err := arista.session.Close(); err != nil && err != io.EOF {
			return err
		}
	}
	arista.Logger().Debug("Closed the session.")

	if arista.conn != nil {
		if err := arista.conn.Close(); err != nil && err != io.EOF {
			return err
		}
	}

	arista.Logger().Debug("Closed the connection.")

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
	return utils.SendCommands(arista.Logger(), arista.writer, arista.reader, arista.errReader, commands...)
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
