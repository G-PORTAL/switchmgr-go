package fscom_test

import (
	"github.com/g-portal/switchmgr-go/pkg/vendors/fscom"
	"testing"
)

var plainConfig = `!version 2.2.0E build 99013
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
 shutdown
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
 switchport mode access
!
interface GigaEthernet0/23
 switchport mode access
 switchport access vlan 6
!
interface GigaEthernet0/24
 switchport trunk vlan-allowed 2000,4,2001-2005
 switchport trunk vlan-untagged 2001
 switchport pvid 2001
 switchport mode trunk
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

func TestListInterfaces(t *testing.T) {
	iosConfig, err := fscom.ParseConfiguration(plainConfig)
	if err != nil {
		t.Fatal(err)
	}

	cfg := fscom.Configuration(iosConfig)

	nics, err := cfg.ListInterfaces()
	if err != nil {
		t.Fatal(err)
	}

	if len(nics) != 30 {
		t.Fatal("Expected 30 interfaces, got", len(nics))
	}

	_, err = cfg.GetInterface("GigaEthernet0/124")
	if err == nil {
		t.Fatal("Expected error for non-existing interface")
	}

	nic, err := cfg.GetInterface("GigaEthernet0/19")
	if err != nil {
		t.Fatal(err)
	}

	if !nic.Enabled {
		t.Fatal("Expected interface to be enabled")
	}

	if nic.Name != "GigaEthernet0/19" {
		t.Fatal("Expected interface name to be GigaEthernet0/19, got", nic.Name)
	}

	if nic.UntaggedVLAN == nil {
		t.Fatal("Expected untagged VLAN to be set")
	}

	if *nic.UntaggedVLAN != 6 {
		t.Fatal("Expected untagged VLAN to be 6, got", *nic.UntaggedVLAN)
	}

	nic, err = cfg.GetInterface("GigaEthernet0/17")
	if err != nil {
		t.Fatal(err)
	}

	if nic.Enabled {
		t.Fatal("Expected interface to be disabled")
	}

	nic, err = cfg.GetInterface("GigaEthernet0/22")
	if err != nil {
		t.Fatal(err)
	}

	if nic.UntaggedVLAN == nil || *nic.UntaggedVLAN != 1 {
		t.Fatal("Expected untagged VLAN to be 1")
	}

	nic, err = cfg.GetInterface("GigaEthernet0/23")
	if err != nil {
		t.Fatal(err)
	}

	if nic.UntaggedVLAN == nil || *nic.UntaggedVLAN != 6 {
		t.Fatal("Expected untagged VLAN to be 6")
	}

	nic, err = cfg.GetInterface("GigaEthernet0/18")
	if err != nil {
		t.Fatal(err)
	}

	if nic.UntaggedVLAN == nil || *nic.UntaggedVLAN != 6 {
		t.Fatal("Expected untagged VLAN to be 6")
	}

	if len(nic.TaggedVLANs) != 2 {
		t.Fatal("Expected 2 tagged VLANs, got", len(nic.TaggedVLANs))
	}

	nic, err = cfg.GetInterface("GigaEthernet0/24")
	if err != nil {
		t.Fatal(err)
	}

	if nic.UntaggedVLAN == nil || *nic.UntaggedVLAN != 2001 {
		t.Fatal("Expected untagged VLAN to be 2001")
	}

	if len(nic.TaggedVLANs) != 6 {
		t.Fatalf("Expected 2 tagged VLANs, got %+v", nic.TaggedVLANs)
	}

	if nic.TaggedVLANs[0] != 4 {
		t.Fatalf("Expected tagged VLAN 0 to be 4, got %d", nic.TaggedVLANs[0])
	}

	if nic.TaggedVLANs[1] != 2000 {
		t.Fatalf("Expected tagged VLAN 1 to be 2000, got %d", nic.TaggedVLANs[1])
	}

	if nic.TaggedVLANs[2] != 2002 {
		t.Fatalf("Expected tagged VLAN 2 to be 2002, got %d", nic.TaggedVLANs[2])
	}

}

var interfaceResponse = `GigaEthernet0/1 is up, line protocol is up
  protocolstatus upTimes 286, downTimes 285, last transition 2000-5-18 23:40:55
  Ifindex is 165, unique port number is 1
  Hardware is Giga-TX, address is 649d.99c5.72df (bia 649d.99c5.72df)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 98224 bits/sec, 104 packets/sec
  5 minutes output rate 88937 bits/sec, 95 packets/sec
  Real time input rate 0%, 91600 bits/sec, 113 packets/sec
  Real time output rate 0%, 114466 bits/sec, 101 packets/sec
     Received 685756905 packets, 299009461365 bytes
     171480 broadcasts, 191420 multicasts, 685394005 ucasts
     0 discard, 2757 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 2757 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 785965482 packets, 365125500746 bytes
     128489455 broadcasts, 22127354 multicasts, 635348673 ucasts
     1 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 590591 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/2 is down, line protocol is down
  protocolstatus upTimes 0, downTimes 0, alloc at 2000-1-1 0:0:18
  Ifindex is 166, unique port number is 2
  Hardware is Giga-TX, address is 649d.99c5.72e0 (bia 649d.99c5.72e0)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-duplex,  Auto-speed,  Flow-Control Off
  5 minutes input rate 0 bits/sec, 0 packets/sec
  5 minutes output rate 0 bits/sec, 0 packets/sec
  Real time input rate 0%, 0 bits/sec, 0 packets/sec
  Real time output rate 0%, 0 bits/sec, 0 packets/sec
     Received 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 0 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 0 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/3 is up, line protocol is up
  protocolstatus upTimes 84, downTimes 83, last transition 2000-5-11 6:24:9
  Ifindex is 167, unique port number is 3
  Hardware is Giga-TX, address is 649d.99c5.72e1 (bia 649d.99c5.72e1)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 2844 bits/sec, 1 packets/sec
  5 minutes output rate 11509 bits/sec, 15 packets/sec
  Real time input rate 0%, 1752 bits/sec, 3 packets/sec
  Real time output rate 0%, 34646 bits/sec, 15 packets/sec
     Received 12594595 packets, 1946130072 bytes
     2783 broadcasts, 182986 multicasts, 12408826 ucasts
     0 discard, 22 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 22 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 174743875 packets, 27805467295 bytes
     128678151 broadcasts, 22136305 multicasts, 23929419 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 412 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/4 is up, line protocol is up
  protocolstatus upTimes 13, downTimes 12, last transition 2000-3-22 9:20:32
  Ifindex is 168, unique port number is 4
  Hardware is Giga-TX, address is 649d.99c5.72e2 (bia 649d.99c5.72e2)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 53 bits/sec, 0 packets/sec
  5 minutes output rate 8214 bits/sec, 12 packets/sec
  Real time input rate 0%, 0 bits/sec, 0 packets/sec
  Real time output rate 0%, 6327 bits/sec, 11 packets/sec
     Received 813138 packets, 191295238 bytes
     209105 broadcasts, 392152 multicasts, 211881 ucasts
     0 discard, 3448 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 3448 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 153463714 packets, 13403391580 bytes
     128478709 broadcasts, 21927672 multicasts, 3057333 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 1238 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/5 is up, line protocol is up
  protocolstatus upTimes 120, downTimes 119, last transition 2000-5-8 0:45:41
  Ifindex is 169, unique port number is 5
  Hardware is Giga-TX, address is 649d.99c5.72e3 (bia 649d.99c5.72e3)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 4244831 bits/sec, 1982 packets/sec
  5 minutes output rate 7089541 bits/sec, 6103 packets/sec
  Real time input rate 0%, 3648365 bits/sec, 1984 packets/sec
  Real time output rate 0%, 6965613 bits/sec, 6037 packets/sec
     Received 3656001854 packets, 990215066209 bytes
     213502 broadcasts, 412425 multicasts, 3655375927 ucasts
     0 discard, 2913464 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 2913464 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 8462708431 packets, 1308734484409 bytes
     128464737 broadcasts, 21906711 multicasts, 8312336983 ucasts
     3 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 679863 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/6 is up, line protocol is up
  protocolstatus upTimes 47, downTimes 46, last transition 2000-4-6 7:23:21
  Ifindex is 170, unique port number is 6
  Hardware is Giga-TX, address is 649d.99c5.72e4 (bia 649d.99c5.72e4)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 2301351 bits/sec, 1145 packets/sec
  5 minutes output rate 2671445 bits/sec, 2430 packets/sec
  Real time input rate 0%, 2328317 bits/sec, 1190 packets/sec
  Real time output rate 0%, 2716601 bits/sec, 2424 packets/sec
     Received 1078454261 packets, 327089862293 bytes
     211642 broadcasts, 391486 multicasts, 1077851133 ucasts
     0 discard, 21659 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 21659 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 2606630937 packets, 415585557728 bytes
     128472785 broadcasts, 21927963 multicasts, 2456230189 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 412 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/7 is up, line protocol is up
  protocolstatus upTimes 85, downTimes 84, last transition 2000-5-7 5:22:55
  Ifindex is 171, unique port number is 7
  Hardware is Giga-TX, address is 649d.99c5.72e5 (bia 649d.99c5.72e5)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 3705423 bits/sec, 2164 packets/sec
  5 minutes output rate 6207621 bits/sec, 5382 packets/sec
  Real time input rate 0%, 3390970 bits/sec, 2153 packets/sec
  Real time output rate 0%, 6350778 bits/sec, 5442 packets/sec
     Received 2691378772 packets, 782882770855 bytes
     212289 broadcasts, 392102 multicasts, 2690774381 ucasts
     0 discard, 500071 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 500071 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 5959344575 packets, 936201277770 bytes
     128469410 broadcasts, 21927348 multicasts, 5808947817 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 20693 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/8 is up, line protocol is up
  protocolstatus upTimes 82, downTimes 81, last transition 2000-5-11 6:26:47
  Ifindex is 172, unique port number is 8
  Hardware is Giga-TX, address is 649d.99c5.72e6 (bia 649d.99c5.72e6)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 2907 bits/sec, 1 packets/sec
  5 minutes output rate 11500 bits/sec, 15 packets/sec
  Real time input rate 0%, 2352 bits/sec, 4 packets/sec
  Real time output rate 0%, 36543 bits/sec, 16 packets/sec
     Received 12902963 packets, 1969460744 bytes
     2811 broadcasts, 182867 multicasts, 12717285 ucasts
     0 discard, 22 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 22 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 175097812 packets, 28274439894 bytes
     128678628 broadcasts, 22136451 multicasts, 24282733 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 405 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/9 is up, line protocol is up
  protocolstatus upTimes 84, downTimes 83, last transition 2000-5-11 6:24:28
  Ifindex is 173, unique port number is 9
  Hardware is Giga-TX, address is 649d.99c5.72e7 (bia 649d.99c5.72e7)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 2857 bits/sec, 1 packets/sec
  5 minutes output rate 11709 bits/sec, 15 packets/sec
  Real time input rate 0%, 2793 bits/sec, 4 packets/sec
  Real time output rate 0%, 38176 bits/sec, 16 packets/sec
     Received 13058985 packets, 2028414656 bytes
     2686 broadcasts, 182926 multicasts, 12873373 ucasts
     0 discard, 22 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 22 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 175338071 packets, 28093541326 bytes
     128678618 broadcasts, 22136313 multicasts, 24523140 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 406 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/10 is up, line protocol is up
  protocolstatus upTimes 84, downTimes 83, last transition 2000-5-11 6:23:46
  Ifindex is 174, unique port number is 10
  Hardware is Giga-TX, address is 649d.99c5.72e8 (bia 649d.99c5.72e8)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 3286 bits/sec, 2 packets/sec
  5 minutes output rate 11953 bits/sec, 15 packets/sec
  Real time input rate 0%, 4674 bits/sec, 3 packets/sec
  Real time output rate 0%, 35915 bits/sec, 15 packets/sec
     Received 13187209 packets, 2008349271 bytes
     2730 broadcasts, 182948 multicasts, 13001531 ucasts
     0 discard, 22 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 22 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 175296958 packets, 28394783261 bytes
     128678529 broadcasts, 22136236 multicasts, 24482193 ucasts
     1 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 406 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/11 is up, line protocol is up
  protocolstatus upTimes 84, downTimes 83, last transition 2000-5-11 6:26:44
  Ifindex is 175, unique port number is 11
  Hardware is Giga-TX, address is 649d.99c5.72e9 (bia 649d.99c5.72e9)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 3050 bits/sec, 2 packets/sec
  5 minutes output rate 12097 bits/sec, 15 packets/sec
  Real time input rate 0%, 5064 bits/sec, 3 packets/sec
  Real time output rate 0%, 37902 bits/sec, 16 packets/sec
     Received 11658921 packets, 1811312829 bytes
     2763 broadcasts, 183023 multicasts, 11473135 ucasts
     0 discard, 22 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 22 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 173724704 packets, 27340224533 bytes
     128678757 broadcasts, 22136225 multicasts, 22909722 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 406 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/12 is up, line protocol is up
  protocolstatus upTimes 84, downTimes 83, last transition 2000-5-11 6:24:5
  Ifindex is 176, unique port number is 12
  Hardware is Giga-TX, address is 649d.99c5.72ea (bia 649d.99c5.72ea)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 2822 bits/sec, 2 packets/sec
  5 minutes output rate 11600 bits/sec, 15 packets/sec
  Real time input rate 0%, 1752 bits/sec, 3 packets/sec
  Real time output rate 0%, 37782 bits/sec, 16 packets/sec
     Received 13206378 packets, 1998325062 bytes
     2775 broadcasts, 182896 multicasts, 13020707 ucasts
     0 discard, 22 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 22 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 175324274 packets, 28202356540 bytes
     128678555 broadcasts, 22136262 multicasts, 24509457 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 406 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/13 is up, line protocol is up
  protocolstatus upTimes 95, downTimes 94, last transition 2000-5-11 6:23:47
  Ifindex is 177, unique port number is 13
  Hardware is Giga-TX, address is 649d.99c5.72eb (bia 649d.99c5.72eb)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 2734 bits/sec, 1 packets/sec
  5 minutes output rate 11407 bits/sec, 15 packets/sec
  Real time input rate 0%, 3761 bits/sec, 3 packets/sec
  Real time output rate 0%, 39105 bits/sec, 16 packets/sec
     Received 13028582 packets, 2002940143 bytes
     49661 broadcasts, 148212 multicasts, 12830709 ucasts
     0 discard, 22 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 22 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 152442020 packets, 27005363550 bytes
     110591177 broadcasts, 17381765 multicasts, 24469078 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 403 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/14 is up, line protocol is up
  protocolstatus upTimes 92, downTimes 91, last transition 2000-5-11 6:23:43
  Ifindex is 178, unique port number is 14
  Hardware is Giga-TX, address is 649d.99c5.72ec (bia 649d.99c5.72ec)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 3714 bits/sec, 2 packets/sec
  5 minutes output rate 12517 bits/sec, 16 packets/sec
  Real time input rate 0%, 1241 bits/sec, 2 packets/sec
  Real time output rate 0%, 36987 bits/sec, 15 packets/sec
     Received 12941557 packets, 1979026809 bytes
     5308 broadcasts, 183071 multicasts, 12753178 ucasts
     0 discard, 282 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 282 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 175830344 packets, 28598615606 bytes
     128675372 broadcasts, 22136156 multicasts, 25018816 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 406 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/15 is up, line protocol is up
  protocolstatus upTimes 84, downTimes 83, last transition 2000-5-11 6:24:55
  Ifindex is 179, unique port number is 15
  Hardware is Giga-TX, address is 649d.99c5.72ed (bia 649d.99c5.72ed)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 3045 bits/sec, 1 packets/sec
  5 minutes output rate 11771 bits/sec, 15 packets/sec
  Real time input rate 0%, 5089 bits/sec, 4 packets/sec
  Real time output rate 0%, 38282 bits/sec, 15 packets/sec
     Received 12901870 packets, 1988942523 bytes
     2003 broadcasts, 182909 multicasts, 12716958 ucasts
     0 discard, 22 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 22 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 175049582 packets, 28101460376 bytes
     128679420 broadcasts, 22136308 multicasts, 24233854 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 406 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/16 is up, line protocol is up
  protocolstatus upTimes 84, downTimes 83, last transition 2000-5-11 6:24:23
  Ifindex is 180, unique port number is 16
  Hardware is Giga-TX, address is 649d.99c5.72ee (bia 649d.99c5.72ee)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-Duplex(Full),  Auto-Speed(1000Mb/s),  Flow-Control Off
  5 minutes input rate 3396 bits/sec, 2 packets/sec
  5 minutes output rate 12751 bits/sec, 15 packets/sec
  Real time input rate 0%, 1688 bits/sec, 3 packets/sec
  Real time output rate 0%, 37233 bits/sec, 15 packets/sec
     Received 12840889 packets, 1971500938 bytes
     2772 broadcasts, 182977 multicasts, 12655140 ucasts
     0 discard, 22 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 22 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 175168191 packets, 28256754693 bytes
     128678655 broadcasts, 22136324 multicasts, 24353212 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 406 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/17 is down, line protocol is down
  protocolstatus upTimes 0, downTimes 0, alloc at 2000-1-1 0:0:18
  Ifindex is 181, unique port number is 17
  Hardware is Giga-TX, address is 649d.99c5.72ef (bia 649d.99c5.72ef)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-duplex,  Auto-speed,  Flow-Control Off
  5 minutes input rate 0 bits/sec, 0 packets/sec
  5 minutes output rate 0 bits/sec, 0 packets/sec
  Real time input rate 0%, 0 bits/sec, 0 packets/sec
  Real time output rate 0%, 0 bits/sec, 0 packets/sec
     Received 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 0 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 0 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/18 is down, line protocol is down
  protocolstatus upTimes 0, downTimes 0, alloc at 2000-1-1 0:0:18
  Ifindex is 182, unique port number is 18
  Hardware is Giga-TX, address is 649d.99c5.72f0 (bia 649d.99c5.72f0)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-duplex,  Auto-speed,  Flow-Control Off
  5 minutes input rate 0 bits/sec, 0 packets/sec
  5 minutes output rate 0 bits/sec, 0 packets/sec
  Real time input rate 0%, 0 bits/sec, 0 packets/sec
  Real time output rate 0%, 0 bits/sec, 0 packets/sec
     Received 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 0 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 0 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/19 is down, line protocol is down
  protocolstatus upTimes 0, downTimes 0, alloc at 2000-1-1 0:0:18
  Ifindex is 183, unique port number is 19
  Hardware is Giga-TX, address is 649d.99c5.72f1 (bia 649d.99c5.72f1)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-duplex,  Auto-speed,  Flow-Control Off
  5 minutes input rate 0 bits/sec, 0 packets/sec
  5 minutes output rate 0 bits/sec, 0 packets/sec
  Real time input rate 0%, 0 bits/sec, 0 packets/sec
  Real time output rate 0%, 0 bits/sec, 0 packets/sec
     Received 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 0 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 0 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/20 is down, line protocol is down
  protocolstatus upTimes 0, downTimes 0, alloc at 2000-1-1 0:0:18
  Ifindex is 184, unique port number is 20
  Hardware is Giga-TX, address is 649d.99c5.72f2 (bia 649d.99c5.72f2)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-duplex,  Auto-speed,  Flow-Control Off
  5 minutes input rate 0 bits/sec, 0 packets/sec
  5 minutes output rate 0 bits/sec, 0 packets/sec
  Real time input rate 0%, 0 bits/sec, 0 packets/sec
  Real time output rate 0%, 0 bits/sec, 0 packets/sec
     Received 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 0 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 0 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/21 is down, line protocol is down
  protocolstatus upTimes 0, downTimes 0, alloc at 2000-1-1 0:0:18
  Ifindex is 185, unique port number is 21
  Hardware is Giga-TX, address is 649d.99c5.72f3 (bia 649d.99c5.72f3)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-duplex,  Auto-speed,  Flow-Control Off
  5 minutes input rate 0 bits/sec, 0 packets/sec
  5 minutes output rate 0 bits/sec, 0 packets/sec
  Real time input rate 0%, 0 bits/sec, 0 packets/sec
  Real time output rate 0%, 0 bits/sec, 0 packets/sec
     Received 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 0 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 0 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/22 is down, line protocol is down
  protocolstatus upTimes 0, downTimes 0, alloc at 2000-1-1 0:0:19
  Ifindex is 186, unique port number is 22
  Hardware is Giga-TX, address is 649d.99c5.72f4 (bia 649d.99c5.72f4)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-duplex,  Auto-speed,  Flow-Control Off
  5 minutes input rate 0 bits/sec, 0 packets/sec
  5 minutes output rate 0 bits/sec, 0 packets/sec
  Real time input rate 0%, 0 bits/sec, 0 packets/sec
  Real time output rate 0%, 0 bits/sec, 0 packets/sec
     Received 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 0 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 0 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/23 is down, line protocol is down
  protocolstatus upTimes 0, downTimes 0, alloc at 2000-1-1 0:0:19
  Ifindex is 187, unique port number is 23
  Hardware is Giga-TX, address is 649d.99c5.72f5 (bia 649d.99c5.72f5)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-duplex,  Auto-speed,  Flow-Control Off
  5 minutes input rate 0 bits/sec, 0 packets/sec
  5 minutes output rate 0 bits/sec, 0 packets/sec
  Real time input rate 0%, 0 bits/sec, 0 packets/sec
  Real time output rate 0%, 0 bits/sec, 0 packets/sec
     Received 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 0 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 0 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
GigaEthernet0/24 is down, line protocol is down
  protocolstatus upTimes 0, downTimes 0, alloc at 2000-1-1 0:0:19
  Ifindex is 188, unique port number is 24
  Hardware is Giga-TX, address is 649d.99c5.72f6 (bia 649d.99c5.72f6)
  MTU 1500 bytes, BW 1000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Auto-duplex,  Auto-speed,  Flow-Control Off
  5 minutes input rate 0 bits/sec, 0 packets/sec
  5 minutes output rate 0 bits/sec, 0 packets/sec
  Real time input rate 0%, 0 bits/sec, 0 packets/sec
  Real time output rate 0%, 0 bits/sec, 0 packets/sec
     Received 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 0 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 0 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
TGigaEthernet0/25 is up, line protocol is up
  protocolstatus upTimes 10, downTimes 9, last transition 2000-3-4 2:41:7
  Ifindex is 189, unique port number is 49
  Hardware is 10Giga-FX-SFP, address is 649d.99c5.72f7 (bia 649d.99c5.72f7)
  MTU 1500 bytes, BW 10000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Full-duplex,  10000Mb/s,  Flow-Control Off
Transceiver Info:
    SFP,LC,1310nm,10000BASE-FX-LR,LOS:no
    SM 10KM
    DDM:YES,Vend:OEM,PN:SFP-10G-LR
    SerialNum:CLS22101702144,Date:2022-10-17
DDM info:
    TX power:-2.88 dBm, RX power:-2.88 dBm
    SFP temperature:54.00 C,supply voltage :3.31V,Bias Current.:35.55mA

  DDM Thresholds:        Low-Alarm    Low-Warning   High-Warning  High-Alarm
  TX power(dBm):            -8.50        -7.50         0.00         0.50
  RX power(dBm):           -18.51       -16.00         0.50         1.50
  SFP temperature(C):          -5            0           70           75
  Supply voltage(v):         3.00         3.10         3.60         3.70
  Bias Current(mA):         15.00        20.00        75.00        85.00
  5 minutes input rate 16075561 bits/sec, 14001 packets/sec
  5 minutes output rate 10385421 bits/sec, 5420 packets/sec
  Real time input rate 0%, 16443261 bits/sec, 14020 packets/sec
  Real time output rate 0%, 9497043 bits/sec, 5479 packets/sec
     Received 17765364835 packets, 3187059797066 bytes
     190960026 broadcasts, 20454690 multicasts, 17553950119 ucasts
     0 discard, 111687773 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 111687773 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 8227189951 packets, 2398260137250 bytes
     547882 broadcasts, 5911452 multicasts, 8220730617 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 3443572 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
TGigaEthernet0/26 is down, line protocol is down
  protocolstatus upTimes 0, downTimes 0, alloc at 2000-1-1 0:0:19
  Ifindex is 190, unique port number is 50
  Hardware is 10Giga-FX, address is 649d.99c5.72f8 (bia 649d.99c5.72f8)
  MTU 1500 bytes, BW 10000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Full-duplex,  10000Mb/s,  Flow-Control Off
  5 minutes input rate 0 bits/sec, 0 packets/sec
  5 minutes output rate 0 bits/sec, 0 packets/sec
  Real time input rate 0%, 0 bits/sec, 0 packets/sec
  Real time output rate 0%, 0 bits/sec, 0 packets/sec
     Received 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 0 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 0 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
TGigaEthernet0/27 is down, line protocol is down
  protocolstatus upTimes 0, downTimes 0, alloc at 2000-1-1 0:0:20
  Ifindex is 191, unique port number is 51
  Hardware is 10Giga-FX, address is 649d.99c5.72f9 (bia 649d.99c5.72f9)
  MTU 1500 bytes, BW 10000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Full-duplex,  10000Mb/s,  Flow-Control Off
  5 minutes input rate 0 bits/sec, 0 packets/sec
  5 minutes output rate 0 bits/sec, 0 packets/sec
  Real time input rate 0%, 0 bits/sec, 0 packets/sec
  Real time output rate 0%, 0 bits/sec, 0 packets/sec
     Received 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 0 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 0 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
TGigaEthernet0/28 is down, line protocol is down
  protocolstatus upTimes 0, downTimes 0, alloc at 2000-1-1 0:0:20
  Ifindex is 192, unique port number is 52
  Hardware is 10Giga-FX, address is 649d.99c5.72fa (bia 649d.99c5.72fa)
  MTU 1500 bytes, BW 10000000 kbit, DLY 10 usec
  Encapsulation ARPA
  Full-duplex,  10000Mb/s,  Flow-Control Off
  5 minutes input rate 0 bits/sec, 0 packets/sec
  5 minutes output rate 0 bits/sec, 0 packets/sec
  Real time input rate 0%, 0 bits/sec, 0 packets/sec
  Real time output rate 0%, 0 bits/sec, 0 packets/sec
     Received 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 align, 0 FCS, 0 symbol
     0 jabber, 0 oversize, 0 undersize
     0 carriersense, 0 collision, 0 fragment
     0 L3 packets, 0 discards, 0 Header errors
     Transmitted 0 packets, 0 bytes
     0 broadcasts, 0 multicasts, 0 ucasts
     0 discard, 0 error, 0 PAUSE
     0 sqettest, 0 deferred, 0 oversize
     0 single, 0 multiple, 0 excessive, 0 late
     0 L3 forwards
VLAN4 is up, line protocol is up
  protocolstatus upTimes 1, downTimes 0, last transition 2000-2-23 3:45:32
  Ifindex is 679
  Hardware is EtherSVI, Address is 649d.99c5.72de(649d.99c5.72de)
  Interface address is 10.16.143.253/20
  MTU 1500 bytes, BW 1000000 kbit, DLY 2000 usec
  Encapsulation ARPA
  ARP type: ARPA, ARP timeout 04:00:00
  Peak input rate 0 pps, output 0 pps
    6952436 packets input, 1085266921 bytes
    Received 3623103 broadcasts, 2318265 multicasts
    0 mpls unicasts, 0 mpls multicasts, 0 mpls input discards
    0 input errors, 930304 discards, 0 protocol unknown
    68345 packets output, 16990458 bytes
    Transmited 2264 broadcasts, 0 multicasts
    0 mpls unicasts, 0 mpls multicasts, 0 mpls output discards
    0 output errors, 0 discards
Null0 is up, line protocol is up
  protocolstatus upTimes 1, downTimes 0, last transition 2000-1-1 0:0:21
  Ifindex is 677
  Hardware is Null
  MTU 1500 bytes, BW 10000000 kbit, DLY 10000 usec
  Encapsulation NULL
`

func TestParseInterfaces(t *testing.T) {
	nics, err := fscom.ParseInterfaces(interfaceResponse)
	if err != nil {
		t.Fatal(err)
	}

	if len(nics) != 29 {
		t.Fatalf("expected 29 nics, got %d", len(nics))
	}

	mac, ok := nics["TGigaEthernet0/25"]
	if !ok {
		t.Fatalf("interface TGigaEthernet0/25 not found")
	}

	if mac.String() != "64:9d:99:c5:72:f7" {
		t.Fatalf("expected mac 64:9d:99:c5:72:f7, got %s", mac)
	}

	_, ok = nics["TGigaEthernet0/123"]
	if ok {
		t.Fatalf("interface TGigaEthernet0/123 found, but should not exist")
	}

}
