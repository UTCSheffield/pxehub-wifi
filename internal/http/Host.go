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

func (h *HttpServer) NewHostiPXE(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain")
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
	mac := ps.ByName("mac")
	if !macRegex.MatchString(mac) {
		script := strings.ReplaceAll(postRegisterScript, "#err ", "")
		fmt.Fprint(w, script)
		return
	}

	hostname := ps.ByName("hostname")

	err := db.CreateHost(mac, hostname, 0, false, h.Database)
	if err != nil {
		script := strings.ReplaceAll(postRegisterScript, "#err ", "")
		fmt.Fprint(w, script)
		log.Print(err)
		return
	}

	script := strings.ReplaceAll(postRegisterScript, "#suc ", "")
	fmt.Fprint(w, script)
}

func (h *HttpServer) NewHost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
	mac := r.FormValue("hostMac")
	if !macRegex.MatchString(mac) {
		http.Error(w, fmt.Sprintf("Invalid Mac Address: %s", mac), http.StatusInternalServerError)
		return
	}
	name := r.FormValue("hostName")
	taskID := r.FormValue("taskID")
	redirect := r.FormValue("redirect") == "true"
	taskPerm := r.FormValue("taskPerm") == "on"

	var taskIDPtr int
	if taskID != "" {
		idInt, err := strconv.Atoi(taskID)
		if err != nil {
			http.Error(w, "Invalid taskID", http.StatusBadRequest)
			return
		}
		taskIDPtr = idInt
	}

	if err := db.CreateHost(mac, name, taskIDPtr, taskPerm, h.Database); err != nil {
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

func (h *HttpServer) EditHost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	name := r.FormValue("hostName")
	mac := r.FormValue("hostMac")
	taskID := r.FormValue("taskID")
	redirect := r.FormValue("redirect") == "true"
	taskPerm := r.FormValue("taskPerm") == "on"

	var idPtr uint
	if id != "" {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid taskID", http.StatusBadRequest)
			return
		}
		idPtr = uint(idInt)
	}

	var taskIDPtr *int
	if taskID != "" {
		idInt, err := strconv.Atoi(taskID)
		if err != nil {
			http.Error(w, "Invalid taskID", http.StatusBadRequest)
			return
		}
		taskIDPtr = &idInt
	}

	if err := db.EditHost(name, mac, taskIDPtr, taskPerm, idPtr, h.Database); err != nil {
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

func (h *HttpServer) DeleteHost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	redirect := r.FormValue("redirect") == "true"

	if err := db.DeleteHost(id, h.Database); err != nil {
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
