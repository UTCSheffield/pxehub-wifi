package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"pxehub/internal/db"
	"pxehub/internal/dnsmasq"
	httpserver "pxehub/internal/http"
)

func readConf(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	conf := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		conf[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return conf, scanner.Err()
}

func main() {
	dirs := []string{
		"/opt/pxehub",
		"/opt/pxehub/http",
		"/opt/pxehub/tftp",
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Printf("Failed to create %s: %v\n", dir, err)
				continue
			}
			log.Printf("Created %s\n", dir)
		} else if err != nil {
			log.Printf("Error checking %s: %v\n", dir, err)
		} else {
		}
	}

	conf, err := readConf("/opt/pxehub/pxehub.conf")
	if err != nil {
		log.Println("Error reading conf:", err)
		return
	}

	dnsList := []string{}
	if val, ok := conf["DNS_SERVERS"]; ok {
		for _, ip := range strings.Split(val, ",") {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				dnsList = append(dnsList, ip)
			}
		}
	}

	database := db.OpenDB("/opt/pxehub/pxehub.db")

	log.Println(db.GetUnassignedWifiKey(database))

	dhcpTftpServer := dnsmasq.DnsmasqServer{
		Iface:       conf["INTERFACE"],
		RangeStart:  conf["DHCP_RANGE_START"],
		RangeEnd:    conf["DHCP_RANGE_END"],
		Mask:        conf["DHCP_MASK"],
		Router:      conf["DHCP_ROUTER"],
		Nameservers: dnsList,
		TFTPDir:     "/opt/pxehub/tftp",
	}

	httpServer := httpserver.HttpServer{
		Address:   conf["HTTP_BIND"],
		Database:  database,
		ExtrasDir: "/opt/pxehub/http",
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
