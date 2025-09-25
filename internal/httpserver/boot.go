package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"pxehub/internal/database"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func (h *HttpServer) BootScript(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mac := ps.ByName("mac")

	ctx := context.Background()

	host, err := gorm.G[database.Host](h.Database).Where("mac = ?", mac).First(ctx)
	if err != nil {
		fmt.Fprintln(w, "Error")
		return
	}

	script, err := gorm.G[database.Script](h.Database).Where("ID = ?", host.ScriptID).First(ctx)
	if err != nil {
		fmt.Fprintln(w, "Error")
		return
	}

	fmt.Fprint(w, script.Content)
}
