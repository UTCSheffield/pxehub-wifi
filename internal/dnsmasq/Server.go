package dnsmasq

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type DnsmasqServer struct {
	Iface       string
	RangeStart  string
	RangeEnd    string
	Mask        string
	Router      string
	Nameservers []string
	TFTPDir     string
	ConfigPath  string
	cmd         *exec.Cmd
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func (d *DnsmasqServer) prepareTFTP() error {
	dir, err := os.MkdirTemp("", "dnsmasq-tftp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp TFTP dir: %w", err)
	}
	d.TFTPDir = dir

	binaries := map[string]string{
		"ipxe.pxe": "https://boot.ipxe.org/ipxe.pxe",
		"ipxe.efi": "https://boot.ipxe.org/ipxe.efi",
	}

	for name, url := range binaries {
		dest := filepath.Join(dir, name)
		log.Printf("Downloading %s...", name)
		if err := downloadFile(url, dest); err != nil {
			return fmt.Errorf("failed to download %s: %w", name, err)
		}
		if err := os.Chmod(dest, 0755); err != nil {
			return fmt.Errorf("failed to chmod %s: %w", name, err)
		}
	}

	script := `#!ipxe
dhcp
chain --autofree http://${next-server}/api/boot/${net0/mac}
	`

	scriptPath := filepath.Join(d.TFTPDir, "autoexec.ipxe")

	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return err
	}

	return nil
}

func (d *DnsmasqServer) generateConfig() (string, error) {
	var opts strings.Builder

	if d.Router != "" {
		opts.WriteString(fmt.Sprintf("dhcp-option=3,%s\n", d.Router))
	}

	if len(d.Nameservers) > 0 {
		opts.WriteString(fmt.Sprintf("dhcp-option=6,%s\n", strings.Join(d.Nameservers, ",")))
		opts.WriteString("no-resolv\n")
		for _, ns := range d.Nameservers {
			opts.WriteString(fmt.Sprintf("server=%s\n", ns))
		}
	}

	confTemplate := `
interface=%s
bind-interfaces
port=0
dhcp-range=%s,%s,%s,12h
enable-tftp
tftp-root=%s

# PXE client architecture matching
dhcp-match=set:archx86, option:client-arch, 0
dhcp-match=set:archx64, option:client-arch, 6
dhcp-match=set:archx32, option:client-arch, 7
dhcp-match=set:archia32, option:client-arch, 1

# Assign boot files based on architecture
dhcp-boot=tag:archx86,ipxe.pxe
dhcp-boot=tag:archx64,ipxe.efi
dhcp-boot=tag:archx32,ipxe.efi
dhcp-boot=tag:archia32,ipxe.pxe

# iPXE override
dhcp-userclass=set:ipxe,iPXE
dhcp-boot=tag:ipxe,autoexec.ipxe

%s

log-dhcp
`

	conf := fmt.Sprintf(confTemplate,
		d.Iface, d.RangeStart, d.RangeEnd,
		d.Mask, d.TFTPDir, opts.String(),
	)

	tmpDir := os.TempDir()
	confPath := filepath.Join(tmpDir, "dnsmasq.conf")

	if err := os.WriteFile(confPath, []byte(conf), 0644); err != nil {
		return "", err
	}
	d.ConfigPath = confPath
	return confPath, nil
}

func (d *DnsmasqServer) Start() error {
	if err := d.prepareTFTP(); err != nil {
		return err
	}

	path, err := exec.LookPath("dnsmasq")
	if err != nil {
		return fmt.Errorf("dnsmasq not found in PATH: %w", err)
	}

	confPath, err := d.generateConfig()
	if err != nil {
		return err
	}

	d.cmd = exec.Command(path, "--no-daemon", "--conf-file="+confPath)

	stdoutPipe, err := d.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	stderrPipe, err := d.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			log.Print(scanner.Text())
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			log.Printf("error reading dnsmasq stdout: %v", err)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			log.Print(scanner.Text())
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			log.Printf("error reading dnsmasq stderr: %v", err)
		}
	}()

	go func() {
		if err := d.cmd.Start(); err != nil {
			log.Printf("Failed to start dnsmasq: %v", err)
			return
		}

		log.Printf(
			"Starting dnsmasq [pid %d] on iface %s serving %s-%s (router %s, nameservers %v, TFTPDir %s)",
			d.cmd.Process.Pid, d.Iface, d.RangeStart, d.RangeEnd, d.Router, d.Nameservers, d.TFTPDir,
		)

		if err := d.cmd.Wait(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 5 {
			} else {
				log.Printf("dnsmasq exited with error: %v", err)
			}
		} else {
			log.Printf("dnsmasq exited cleanly")
		}
	}()

	return nil
}

func (d *DnsmasqServer) Stop() error {
	if d.cmd != nil && d.cmd.Process != nil {
		if err := d.cmd.Process.Signal(os.Interrupt); err != nil {
			log.Printf("Failed to send interrupt: %v", err)
		}

		done := make(chan error, 1)
		go func() {
			done <- d.cmd.Wait()
		}()

		select {
		case <-done:

		case <-time.After(3 * time.Second):
			log.Printf("dnsmasq did not exit, killing...")
			if err := d.cmd.Process.Kill(); err != nil {
				log.Printf("Failed to kill dnsmasq: %v", err)
			}
		}
	}

	if d.ConfigPath != "" {
		if err := os.Remove(d.ConfigPath); err == nil {
			log.Printf("Deleted config: %s", d.ConfigPath)
		}
	}

	if d.TFTPDir != "" {
		if err := os.RemoveAll(d.TFTPDir); err == nil {
			log.Printf("Deleted TFTP dir: %s", d.TFTPDir)
		}
		d.TFTPDir = ""
	}

	return nil
}
