package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"pxehub/internal/db"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var script = `#!ipxe
#suc chain --autofree http://${next-server}/api/boot/${net0/mac}
#err echo "Host registering failed"
#err read end
`

func (h *HttpServer) NewHost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain")
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
	mac := ps.ByName("mac")
	if !macRegex.MatchString(mac) {
		switch r.Method {
		case "GET":
			script := strings.ReplaceAll(script, "#err ", "")
			fmt.Fprint(w, script)
		case "POST":
			fmt.Fprint(w, "Invalid MAC Address")
		}
		return
	}

	hostname := ps.ByName("hostname")

	err := db.CreateHost(mac, hostname, h.Database)
	if err != nil {
		switch r.Method {
		case "GET":
			script := strings.ReplaceAll(script, "#err ", "")
			fmt.Fprint(w, script)
		case "POST":
			fmt.Fprint(w, "Error")
		}
		log.Print(err)
		return
	}

	switch r.Method {
	case "GET":
		script := strings.ReplaceAll(script, "#suc ", "")
		fmt.Fprint(w, script)
	case "POST":
		fmt.Fprint(w, "Success")
	}
}
