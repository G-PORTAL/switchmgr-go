package sshconn

import (
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

// Dial is a copy/thin wrapper around the real ssh.Dial() that
// additionally sets a deadline on the underlying connection.
//
// It's necessary because a ssh.Dial() can happens right after the
// reboot command is issued. The net.Dial() itself is successful but
// then during the ssh session set up the TCP connection ends
// because of the reboot. The golang "ssh" package has no concept of
// "ssh -o ServerAliveInterval=10" or similar, so the code will
// just hang in a read forever. This was observed running the
// spread "cerberus" tests on ubuntu 23.04.
//
// Note that half of the function is just a copy of
// golang.org/x/crypto/ssh/client.go:func Dial()
// Source: https://github.com/canonical/spread/pull/160/files
var Dial = func(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	conn, err := net.DialTimeout(network, addr, config.Timeout)
	if err != nil {
		return nil, err
	}

	// See e.g. https://github.com/golang/go/issues/51926
	if config.Timeout > 0 {
		if err := conn.SetDeadline(time.Now().Add(config.Timeout)); err != nil {
			return nil, err
		}
		defer func() {
			_ = conn.SetDeadline(time.Time{})
		}()
	}
	// end of the new code
	c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(c, chans, reqs), nil
}
