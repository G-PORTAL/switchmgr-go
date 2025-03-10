package juniper_els

import (
	"fmt"
	"github.com/openshift-telco/go-netconf-client/netconf"
	"github.com/openshift-telco/go-netconf-client/netconf/message"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type JuniperDriver interface {
	SyncRPC(operation message.RPCMethod, timeout int32) (*message.RPCReply, error)
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
		session: &MockDriver{
			mockBasePath: directory + "/testdata/juniper-els/",
		},
	}
}

func (j *MockDriver) SyncRPC(operation message.RPCMethod, timeout int32) (*message.RPCReply, error) {
	mockFileRegex := regexp.MustCompile(`<([^>]+)>`)

	rpc, ok := operation.(*message.RPC)
	if !ok {
		return nil, fmt.Errorf("operation is not an RPC method, got %T", operation)
	}

	rpcPayload, ok := rpc.Data.(string)
	if !ok {
		return nil, fmt.Errorf("data is not a byte array, got %T", rpc)
	}

	match := mockFileRegex.FindStringSubmatch(strings.Replace(rpcPayload, "/", "", 1))
	if len(match) != 2 {
		return nil, fmt.Errorf("could not find mock file for RPC method")
	}
	reply := &message.RPCReply{}
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

func (j *LiveDriver) SyncRPC(operation message.RPCMethod, timeout int32) (*message.RPCReply, error) {
	return j.session.SyncRPC(operation, timeout)
}
func (j *LiveDriver) Close() error {
	return j.session.Close()
}
