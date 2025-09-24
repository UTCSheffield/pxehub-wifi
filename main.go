package main

import (
	"log"

	"pxehub/internal/dnsmasq"
)

func main() {
	//TODO: config file / command line args
	server := dnsmasq.DnsmasqServer{
		Iface:      "eth0",
		RangeStart: "192.168.100.10",
		RangeEnd:   "192.168.100.50",
		Router:     "192.168.100.1",
		Mask:       "255.255.255.0",
		TFTPDir:    "/srv/tftp",
		NextServer: "192.168.100.1",
		BootFile:   "ipxe.efi",
		IpxeScript: "autoexec.ipxe",
	}

	if err := server.Start(); err != nil {
		log.Fatalf("dnsmasq failed: %v", err)
	}
}
