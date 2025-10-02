package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"pxehub/internal/db"
	"regexp"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var postRegisterScript = `#!ipxe
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
			script := strings.ReplaceAll(postRegisterScript, "#err ", "")
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
			script := strings.ReplaceAll(postRegisterScript, "#err ", "")
			fmt.Fprint(w, script)
		case "POST":
			fmt.Fprint(w, "Error")
		}
		log.Print(err)
		return
	}

	switch r.Method {
	case "GET":
		script := strings.ReplaceAll(postRegisterScript, "#suc ", "")
		fmt.Fprint(w, script)
	case "POST":
		fmt.Fprint(w, "Success")
	}
}

func (h *HttpServer) EditHost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var idPtr int
	if id != "" {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid taskID", http.StatusBadRequest)
			return
		}
		idPtr = idInt
	}

	name := r.FormValue("hostName")
	mac := r.FormValue("taskMac")
	taskID := r.FormValue("taskID")
	redirect := r.FormValue("redirect") == "true"

	var taskIDPtr int
	if taskID != "" {
		idInt, err := strconv.Atoi(taskID)
		if err != nil {
			http.Error(w, "Invalid taskID", http.StatusBadRequest)
			return
		}
		taskIDPtr = idInt
	}

	if err := db.EditHost(name, mac, taskIDPtr, idPtr, h.Database); err != nil {
		http.Error(w, "Update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if redirect {
		http.Redirect(w, r, "/hosts", http.StatusSeeOther)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}
}
