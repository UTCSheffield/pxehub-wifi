package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"pxehub/internal/db"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

func (h *HttpServer) BootScript(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain")
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
	mac := ps.ByName("mac")
	if !macRegex.MatchString(mac) {
		fmt.Fprint(w, "Error: Invalid MAC Address")
		return
	}

	script, err := db.GetScriptByMAC(mac, h.Database, true)
	if err != nil {
		fmt.Fprint(w, "Error")
		log.Print("Error in http request", err)
		return
	}

	fmt.Fprint(w, script)
}
