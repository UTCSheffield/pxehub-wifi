package httpserver

import (
	"fmt"
	"net/http"
	"pxehub/internal/db"

	"github.com/julienschmidt/httprouter"
)

func (h *HttpServer) NewTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := r.FormValue("taskName")
	script := r.FormValue("taskScript")
	redirect := r.FormValue("redirect") == "true"

	if name == "" || script == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	db.CreateTask(name, script, h.Database)

	if redirect {
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
	} else {
		fmt.Fprintln(w, "Success")
	}
}

func (h *HttpServer) EditTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	name := r.FormValue("taskName")
	script := r.FormValue("taskScript")
	redirect := r.FormValue("redirect") == "true"

	if err := db.EditTask(name, script, id, h.Database); err != nil {
		http.Error(w, "Update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if redirect {
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}
}
