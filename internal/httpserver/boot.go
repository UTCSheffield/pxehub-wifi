package httpserver

import (
	"fmt"
	"net/http"
	"pxehub/internal/db"

	"github.com/julienschmidt/httprouter"
)

func (h *HttpServer) BootScript(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mac := ps.ByName("mac")

	script, err := db.GetScriptByMAC(mac, h.Database)
	if err != nil {
		fmt.Fprint(w, "Error")
		return
	}

	fmt.Fprint(w, script)
}
