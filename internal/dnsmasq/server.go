package dnsmasq

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type DnsmasqServer struct {
	Iface      string
	RangeStart string
	RangeEnd   string
	Router     string
	Mask       string
	TFTPDir    string
	NextServer string
	BootFile   string // default PXE loader (for non-iPXE clients)
	IpxeScript string // script served if client is iPXE
	ConfigPath string
	cmd        *exec.Cmd
}

func NewDnsmasqServer(iface, rangeStart, rangeEnd, router, mask, tftpDir, nextServer string) *DnsmasqServer {
	return &DnsmasqServer{
		Iface:      iface,
		RangeStart: rangeStart,
		RangeEnd:   rangeEnd,
		Router:     router,
		Mask:       mask,
		TFTPDir:    tftpDir,
		NextServer: nextServer,
	}
}

func (d *DnsmasqServer) GenerateConfig() (string, error) {
	conf := fmt.Sprintf(`
interface=%s
bind-interfaces
dhcp-range=%s,%s,%s,12h
dhcp-option=3,%s
dhcp-option=66,%s
enable-tftp
tftp-root=%s

# PXE chainloading
dhcp-boot=%s
dhcp-userclass=set:ipxe,iPXE
dhcp-boot=tag:ipxe,%s

log-dhcp
`, d.Iface, d.RangeStart, d.RangeEnd, d.Mask, d.Router, d.NextServer,
		d.TFTPDir, d.BootFile, d.IpxeScript)

	tmpDir := os.TempDir()
	confPath := filepath.Join(tmpDir, "dnsmasq.conf")

	if err := os.WriteFile(confPath, []byte(conf), 0644); err != nil {
		return "", err
	}
	d.ConfigPath = confPath
	return confPath, nil
}

func (d *DnsmasqServer) Start() error {
	path, err := exec.LookPath("dnsmasq")
	if err != nil {
		return fmt.Errorf("dnsmasq not found in PATH: %w", err)
	}

	confPath, err := d.GenerateConfig()
	if err != nil {
		return err
	}

	d.cmd = exec.Command(path, "--no-daemon", "--conf-file="+confPath)
	d.cmd.Stdout = os.Stdout
	d.cmd.Stderr = os.Stderr

	log.Printf("Starting dnsmasq on iface %s serving %s-%s (router %s, next-server %s, bootfile autoexec.ipxe)",
		d.Iface, d.RangeStart, d.RangeEnd, d.Router, d.NextServer)

	return d.cmd.Start()
}

func (d *DnsmasqServer) Stop() error {
	if d.cmd != nil && d.cmd.Process != nil {
		return d.cmd.Process.Kill()
	}
	return nil
}
