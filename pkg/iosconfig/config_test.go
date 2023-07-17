package iosconfig_test

import (
	"github.com/g-portal/switchmgr-go/pkg/iosconfig"
	"testing"
)

var plainConfig = `
!version 2.2.0E build 99013
service timestamps log date
service timestamps debug date
service password-encryption
!
hostname switch06_tw
!
!
lldp run
!
!
!
!
!
!
spanning-tree mode rstp
!
!
!
!
!
!
!
!
!
!
!
!
!
!
!
!
aaa authentication login ssh local
aaa authentication login default local
aaa authentication enable default none
aaa authorization exec default local
!
username admin password 7 011c36405d337d241938322a461a27567f
username webinterface password 7 014b0e206e36442a093e702f481122412c1a15581a44125
42d221d252b57626933
!
!
!
!
!
interface Null0
!
interface GigaEthernet0/1
 switchport trunk vlan-untagged 1
 switchport mode trunk
!
interface GigaEthernet0/2
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/3
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/4
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/5
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/6
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/7
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/8
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/9
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/10
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/11
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/12
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/13
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/14
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/15
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/16
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/17
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/18
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/19
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/20
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/21
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/22
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/23
 switchport mode trunk
 switchport pvid 6
!
interface GigaEthernet0/24
 switchport mode trunk
 switchport pvid 6
!
interface TGigaEthernet0/25
 switchport mode trunk
!
interface TGigaEthernet0/26
!
interface TGigaEthernet0/27
!
interface TGigaEthernet0/28
!
interface VLAN4
 ip address 10.16.143.253 255.255.240.0
 no ip directed-broadcast
!
!
!
vlan 4
 name platform-mgmt
!
vlan 6
 name intern_daemon
!
vlan 1,4,6
!
!
!
!
!
!
!
!
!
!
!
!
!
!
!
ip route default 10.16.128.1 
ip exf
!
ipv6 exf
!
!
!
!
ip telnet attack-defense
no ip telnet enable
!
ip http language english
ip http server
!
!
!
!
line console 0
 exec-timeout 1500
!
line vty 0 31
 exec-timeout 1500
!
!
!
ip sshd auth-method ssh
ip sshd enable
!
!
!
!
!
!
!Pending configurations for absent linecards:
!
!No configurations pending global

`

func TestConfig_Interfaces(t *testing.T) {
	//cfg := iosconfig.
	cfg, err := iosconfig.Parse(plainConfig)
	if err != nil {
		t.Error(err)
	}

	if len(cfg) != 57 {
		t.Errorf("Config length is not 57, it is %d", len(cfg))
	}

	if len(cfg.Interfaces()) != 30 {
		t.Errorf("Interface count is not 30, it is %d", len(cfg.Interfaces()))
	}

	nicConfig, err := cfg.Interface("GigaEthernet0/12")
	if err != nil {
		t.Error(err)
	}

	nicConfigMode := nicConfig.GetStringValue("switchport mode", "")
	if nicConfigMode != "trunk" {
		t.Errorf("switchport mode is not trunk, it is %q", nicConfigMode)
	}

	nicConfigPVID := nicConfig.GetInt32Value("switchport pvid", 0)
	if nicConfigPVID != 6 {
		t.Errorf("switchport pvid is not 6, it is %d", nicConfigPVID)
	}
}
