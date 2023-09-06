package fscom

import (
	"bufio"
	"fmt"
	"github.com/g-portal/switchmgr-go/pkg/config"
	"github.com/g-portal/switchmgr-go/pkg/vendors/unimplemented"
	"golang.org/x/crypto/ssh"
	"golang.org/x/exp/slices"
	"io"
	"regexp"
	"strings"
	"time"
)

type FSCom struct {
	unimplemented.Unimplemented

	LoginCommands []string

	conn    *ssh.Client
	session *ssh.Session

	writer io.WriteCloser
	reader io.Reader
}

// Connect Connecting to a FiberStore switch using SSH
func (fs *FSCom) Connect(cfg config.Connection) error {
	fs.Logger().SetPrefix(fmt.Sprintf("[fscom/%s]", cfg.Host))

	var err error
	sshConfig := &ssh.ClientConfig{
		User: cfg.Username,
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

func (fs *FSCom) Disconnect() error {
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

// / save saves the configuration to startup config
func (fs *FSCom) Save() error {
	output, err := fs.SendCommands("write")
	if err != nil {
		return err
	}

	if !slices.Contains([]string{"Saving current configuration...\r\n\nOK!", "Building configuration...\r\n\n[OK]"}, strings.TrimSpace(output)) {
		return fmt.Errorf("failed to save the configuration: %s", output)
	}

	return err
}

// SendCommands sends a command to the switch and returns the output
func (fs *FSCom) SendCommands(commands ...string) (string, error) {
	output := ""

	startTime := time.Now()
	fs.Logger().Debugf("SendCommands %q", commands)
	defer func() {
		fs.Logger().Debugf("SendCommands %q took %s", commands, time.Since(startTime).String())
	}()

	reader := bufio.NewReader(fs.reader)
	for _, s := range commands {
		// send the command to the switch
		if _, err := fmt.Fprintf(fs.writer, "%s\n\n", s); err != nil {
			return "", fmt.Errorf("failed to send command %q: %s", s, err)
		}

		// read the output of the command
		commandOutput, err := readUntil(s, reader, fs.writer)
		if err != nil {
			return "", err
		}

		fs.Logger().Debugf("cmd: %q output: %q", s, commandOutput)
		output += commandOutput
	}

	return output, nil
}

var moreRgx = regexp.MustCompile(`--More--[\s\\b]+`)

func readUntil(command string, reader *bufio.Reader, writer io.Writer) (string, error) {
	output := ""

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if strings.TrimSpace(line) == "" {
				break
			}

			continue
		}

		// sending enter when switch asks for more even if terminal length is set to 0
		if moreRgx.MatchString(line) {
			if _, err := fmt.Fprintf(writer, "\n"); err != nil {
				return "", fmt.Errorf("failed to send space: %s", err)
			}

			line = moreRgx.ReplaceAllString(line, "")
		}

		if commandFinished(line) {
			break
		}

		if !strings.HasSuffix(line, strings.TrimSpace(command)+"\r\n") && line != "" {
			output += fmt.Sprintf("%s\n", line)
		}
	}

	return output, nil
}

func commandFinished(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasSuffix(line, "#") || strings.HasSuffix(line, ">")
}
