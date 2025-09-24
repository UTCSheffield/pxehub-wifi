package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"pxehub/internal/dnsmasq"
	"pxehub/internal/httpserver"
)

func main() {
	//TODO: config file / command line args
	dhcpTftpServer := dnsmasq.DnsmasqServer{
		Iface:       "virbr1",
		RangeStart:  "192.168.100.10",
		RangeEnd:    "192.168.100.254",
		Mask:        "255.255.255.0",
		Router:      "192.168.100.1",
		Nameservers: []string{"1.1.1.1"},
	}

	httpServer := httpserver.HttpServer{
		Iface: "virbr1",
		Port:  80,
	}

	if err := dhcpTftpServer.Start(); err != nil {
		fmt.Printf("dnsmasq failed: %v", err)
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}

	if err := httpServer.Start(); err != nil {
		fmt.Printf("http failed: %v", err)
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("Shutting down dnsmasq...")
	if err := dhcpTftpServer.Stop(); err != nil {
		log.Printf("failed to stop dnsmasq: %v", err)
	}
	log.Println("Shutting down http...")
	if err := httpServer.Stop(); err != nil {
		log.Printf("failed to stop http: %v", err)
	}
}
