package juniper

import (
	"encoding/xml"
	"github.com/Juniper/go-netconf/netconf"
)

func (j *Juniper) GetRunningConfig() (*junosConfiguration, error) {
	reply, err := j.session.Exec(netconf.RawMethod("<get-config><source><running/></source></get-config>"))
	if err != nil {
		return nil, err
	}

	var cfg junosConfiguration
	if err := xml.Unmarshal([]byte(reply.RawReply), &cfg); err != nil {
		return nil, err
	}
	return &cfg, err

}
