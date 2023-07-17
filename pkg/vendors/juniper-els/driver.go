package juniper_els

import (
	"fmt"
	"github.com/Juniper/go-netconf/netconf"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type JuniperDriver interface {
	Exec(methods ...netconf.RPCMethod) (*netconf.RPCReply, error)
	Close() error
}

type MockDriver struct {
	mockBasePath string
}

func NewMockDriver() *JuniperELS {
	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_, err = os.Stat(directory + "/testdata")
	for os.IsNotExist(err) {
		directory, _ = filepath.Abs(directory + "/..")
		_, err = os.Stat(directory + "/testdata")
	}
	return &JuniperELS{
		vlanMapping:    make(map[string]int32),
		interfaceVlans: make(map[string]junosVlanMapEntry),
		session: &MockDriver{
			mockBasePath: directory + "/testdata/juniper-els/",
		},
	}
}

func (j *MockDriver) Exec(methods ...netconf.RPCMethod) (*netconf.RPCReply, error) {
	if len(methods) != 1 {
		return nil, fmt.Errorf("mock driver does not support multiple RPC methods")
	}

	method, ok := methods[0].(netconf.RawMethod)
	if !ok {
		return nil, fmt.Errorf("mock driver only supports raw RPC methods")
	}

	mockFileRegex := regexp.MustCompile(`<([^>]+)>`)

	match := mockFileRegex.FindStringSubmatch(strings.Replace(string(method), "/", "", 1))
	if len(match) != 2 {
		return nil, fmt.Errorf("could not find mock file for RPC method")
	}
	reply := &netconf.RPCReply{}
	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s.xml", j.mockBasePath, match[1]))
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	reply.RawReply = string(data)
	return reply, nil
}

func (j *MockDriver) Close() error {
	return nil
}

type LiveDriver struct {
	session *netconf.Session
}

func (j *LiveDriver) Exec(methods ...netconf.RPCMethod) (*netconf.RPCReply, error) {
	return j.session.Exec(methods...)
}
func (j *LiveDriver) Close() error {
	return j.session.Close()
}
