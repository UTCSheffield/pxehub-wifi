package httpserver

import (
	"fmt"
	"net/http"
	"pxehub/internal/db"
	"regexp"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (h *HttpServer) GetWifiKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain")
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
	mac := ps.ByName("mac")
	if !macRegex.MatchString(mac) {
		http.Error(w, "Invalid Mac Address", http.StatusInternalServerError)
		return
	}

	host, err := db.GetHostByMAC(mac, h.Database)
	if err != nil {
		http.Error(w, "Fetching Host Failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	key, err := db.GetOrAssignWifiKeyToHost(host.ID, h.Database)
	if err != nil {
		return
	}

	fmt.Fprint(w, key.Key)
}

func (h *HttpServer) NewWifiKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := r.FormValue("wifiKey")
	redirect := r.FormValue("redirect") == "true"

	if key == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	db.CreateWifiKey(key, h.Database)

	if redirect {
		http.Redirect(w, r, "/wifikeys", http.StatusSeeOther)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}
}

func (h *HttpServer) EditWifiKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	var idPtr uint
	if id != "" {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid taskID", http.StatusBadRequest)
			return
		}
		idPtr = uint(idInt)
	}

	key := r.FormValue("wifiKey")
	redirect := r.FormValue("redirect") == "true"

	if err := db.EditWifiKey(key, idPtr, h.Database); err != nil {
		http.Error(w, "Update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if redirect {
		http.Redirect(w, r, "/wifikeys", http.StatusSeeOther)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}
}

func (h *HttpServer) DeleteWifiKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	var idPtr uint
	if id != "" {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid wifiKeyID", http.StatusBadRequest)
			return
		}
		idPtr = uint(idInt)
	}

	redirect := r.FormValue("redirect") == "true"

	if err := db.DeleteWifiKey(idPtr, h.Database); err != nil {
		http.Error(w, "Update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if redirect {
		http.Redirect(w, r, "/wifikeys", http.StatusSeeOther)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}
}
