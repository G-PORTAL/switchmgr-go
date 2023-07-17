# Switch Manager Go

This is a simple switch manager written in Go

# fscom
mode trunk -> heißt alle vlans tagged außer wenn "vlan-untagged XXX" gesetzt ist
    switchport pvid 6 -> oder auch untagged (default 1)
    switchport trunk allowed vlan 6,7,8-10 -> tagged

mode access -> heißt nur ein vlan untagged (default)
    switchport access vlan 6 -> untagged

interface GigaEthernet0/4
    switchport trunk vlan-allowed 1,4,6
    switchport trunk vlan-untagged 1
    switchport pvid 1
    switchport mode trunk
