package fsos_s3

import (
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

const Vendor registry.Vendor = "fsos_s3"

type FSComS3 struct {
	unimplemented.Unimplemented

	LoginCommands []string

	conn    *ssh.Client
	session *ssh.Session

	writer    io.WriteCloser
	reader    io.Reader
	errReader io.Reader
}

func (fs *FSComS3) Vendor() registry.Vendor {
	return Vendor
}

// Connect Connecting to a FiberStore switch using SSH
func (fs *FSComS3) Connect(cfg config.Connection) error {
	fs.Logger().SetPrefix(fmt.Sprintf("[fscom/%s]", cfg.Host))

	var err error
	sshConfig := &ssh.ClientConfig{
		User:    cfg.Username,
		Timeout: 30 * time.Second,
		Auth: []ssh.AuthMethod{
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
				return []string{}, err
			}),
		},
		HostKeyAlgorithms: []string{ssh.KeyAlgoDSA, ssh.KeyAlgoRSA, ssh.KeyAlgoECDSA256, ssh.KeyAlgoED25519},
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
	}

	if cfg.Password != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.RetryableAuthMethod(ssh.Password(strings.TrimSpace(cfg.Password)), 3))
	}

	fs.conn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%v", cfg.Host, cfg.Port), sshConfig)
	if err != nil {
		return err
	}

	if fs.session, err = fs.conn.NewSession(); err != nil {
		return err
	}

	// Set up terminal modes
	fs.Logger().Debugf("Requesting pseudo terminal.")
	if err = fs.session.RequestPty("xterm", 5000, 5000, ssh.TerminalModes{}); err != nil {
		return fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	// create pipe to the stdin of the SSH process
	fs.writer, err = fs.session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %s", err)
	}

	// create pipe to the stdout of the SSH process
	fs.reader, err = fs.session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %s", err)
	}

	// create pipe to the stderr of the SSH process
	fs.errReader, err = fs.session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %s", err)
	}

	// start the shell
	fs.Logger().Debugf("Starting the session.")
	if err = fs.session.Shell(); err != nil {
		return fmt.Errorf("failed to start the shell: %s", err)
	}

	// entering the shell
	if _, err = fs.SendCommands(fs.LoginCommands...); err != nil {
		return err
	}

	return nil
}

func (fs *FSComS3) Disconnect() error {
	if fs.session != nil {
		if err := fs.session.Close(); err != nil && err != io.EOF {
			return err
		}
	}

	fs.Logger().Debug("Closed the session.")

	if fs.conn != nil {
		if err := fs.conn.Close(); err != nil && err != io.EOF {
			return err
		}
	}

	fs.Logger().Debug("Closed the connection.")

	return nil
}

// Save  saves the configuration to startup config
func (fs *FSComS3) Save() error {
	output, err := fs.SendCommands("write")
	if err != nil {
		return err
	}

	if !strings.Contains(strings.TrimSpace(output), "OK") {
		return fmt.Errorf("failed to save the configuration: %s", output)
	}

	return err
}

// SendCommands sends a command to the switch and returns the output
func (fs *FSComS3) SendCommands(commands ...string) (string, error) {
	outputs, err := utils.SendCommands(fs.Logger(), fs.writer, fs.reader, fs.errReader, commands...)
	if err != nil {
		return "", err
	}

	return strings.Join(outputs, "\n"), nil
}

func init() {
	registry.RegisterVendor(Vendor, &FSComS3{
		LoginCommands: []string{
			"enter", "terminal length 0",
		},
	})
}
