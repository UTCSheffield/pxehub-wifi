# pxehub
### This fork includes wifi key assignment to hosts
Dynamic iPXE booting and host tracking.

* Allows registering of hosts
* Allows setting an iPXE script per host
* Identifies hosts by mac address
* Combines DHCP and TFTP together
* Web UI for managing hosts and scripts

TODO:
* Web Authentication and Users

## Config Example
```
HTTP_BIND=192.168.1.1:80
INTERFACE=eth0
DHCP_RANGE_START=192.168.1.10
DHCP_RANGE_END=192.168.1.254
DHCP_MASK=255.255.255.0
DHCP_ROUTER=192.168.1.1
```

## Fetching a Wifi Key
Send a GET request to this url:
`http://{server}/api/get/wifikey/{mac}`
{server} is the IP or hostname of your server
{mac} is the mac address of the device, lowercase, delimited by colons