package httpserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

type HttpServer struct {
	Iface    string
	Port     uint16
	Server   *http.Server
	Database *gorm.DB
}

func getInterfaceIP(name string) (string, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return "", err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			if v.IP.To4() != nil {
				return v.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no IPv4 address found for interface %s", name)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var method string
		if strings.Contains(r.URL.Path, "/api/new/") {
			method = "POST"
		} else {
			method = r.Method
		}
		log.Printf("%s %s from %s",
			method,
			r.URL.Path,
			r.RemoteAddr,
		)
		next.ServeHTTP(w, r)
	})
}

func (h *HttpServer) Start() error {
	if h.Iface == "" {
		return fmt.Errorf("interface name not set")
	}
	if h.Port <= 0 {
		return fmt.Errorf("port not set")
	}

	ip, err := getInterfaceIP(h.Iface)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", ip, h.Port)

	log.Printf("Starting http on iface %s listening on %s", h.Iface, addr)

	router := httprouter.New()

	router.GET("/api/boot/:mac", h.BootScript)
	router.GET("/api/new/host/:mac/:hostname", h.NewHost)
	router.POST("/api/new/host/:mac/:hostname", h.NewHost)

	h.Server = &http.Server{
		Addr:    addr,
		Handler: loggingMiddleware(router),
	}

	errChan := make(chan error, 1)
	go func() {
		if err := h.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(100 * time.Millisecond):
		log.Printf("Started http")
		return nil
	}
}

func (h *HttpServer) Stop() error {
	if h.Server == nil {
		return fmt.Errorf("server not running")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return h.Server.Shutdown(ctx)
}
