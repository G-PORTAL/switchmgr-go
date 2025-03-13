package utils

import (
	"bufio"
	"fmt"
	"github.com/charmbracelet/log"
	"io"
	"regexp"
	"strings"
	"time"
)

func SendCommands(logger *log.Logger, stdin io.WriteCloser, stdout io.Reader, commands ...string) ([]string, error) {
	startTime := time.Now()
	reader := bufio.NewReaderSize(stdout, 10240)
	logger.Debugf("SendCommands %q", commands)
	defer func() {
		logger.Debugf("SendCommands %q took %s", commands, time.Since(startTime).String())
	}()

	output := make([]string, 0)
	for _, s := range commands {
		// send the command to the switch
		if _, err := fmt.Fprintf(stdin, "%s\n\n", s); err != nil {
			return nil, fmt.Errorf("failed to send command %q: %s", s, err)
		}

		// read the output of the command
		commandOutput, err := readUntil(s, reader, stdin, 30*time.Second)
		if err != nil {
			return nil, fmt.Errorf("failed to read command output: %s", err)
		}

		logger.Debugf("cmd: %q output: %q", s, commandOutput)
		output = append(output, strings.TrimSpace(commandOutput))
	}

	return output, nil
}

var moreRgx = regexp.MustCompile(`--More--[\s\\b]+`)

func readUntil(command string, reader *bufio.Reader, writer io.Writer, duration time.Duration) (string, error) {
	output := ""

	for {
		line, err := readWithTimeout(reader, duration)
		if err != nil {
			if strings.TrimSpace(line) == "" {
				break
			}
			return "", fmt.Errorf("failed to read line: %s", err)
		}

		// sending enter when switch asks for more even if terminal length is set to 0
		if moreRgx.MatchString(line) {
			if _, err := fmt.Fprintf(writer, "\n\n\n"); err != nil {
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

func readWithTimeout(reader *bufio.Reader, timeout time.Duration) (string, error) {
	lineChan := make(chan string)
	errChan := make(chan error)

	go func() {
		line, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err
		} else {
			lineChan <- line
		}
	}()

	select {
	case line := <-lineChan:
		return line, nil
	case err := <-errChan:
		return "", err
	case <-time.After(timeout):
		return "", fmt.Errorf("timeout reached while reading line")
	}
}
