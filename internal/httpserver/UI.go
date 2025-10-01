package httpserver

import (
	"html/template"
	"net/http"
	"os"
	"pxehub/internal/db"
	"pxehub/ui"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (h *HttpServer) UI(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")

	path := strings.Trim(r.URL.Path, "/")

	switch path {
	case "":
		files := []string{
			"base.html",
			"index.html",
		}

		tmpl, err := template.ParseFS(ui.Content, files...)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		graphData, graphData1, graphData2, graphDates, err := db.GetRequestGraphData(time.Now(), h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		totalRequests, totalUnregisteredRequests, totalRegisteredRequests, err := db.GetTotalRequestCounts(h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		caser := cases.Title(language.English)
		data := map[string]any{
			"Title":                 caser.String("index"),
			"Username":              "index",
			"Path":                  r.URL.Path,
			"GraphData":             graphData,
			"UnregisteredGraphData": graphData1,
			"RegisteredGraphData":   graphData2,
			"GraphDates":            graphDates,
			"Month":                 time.Now().Month(),
			"UnregisteredRequests":  totalUnregisteredRequests,
			"RegisteredRequests":    totalRegisteredRequests,
			"TotalRequests":         totalRequests,
		}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "hosts":
		files := []string{
			"base.html",
			"hosts.html",
		}

		tmpl, err := template.ParseFS(ui.Content, files...)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		caser := cases.Title(language.English)
		data := map[string]any{
			"Title":    caser.String("hosts"),
			"Username": "hosts",
			"Path":     r.URL.Path,
		}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "tasks":
		files := []string{
			"base.html",
			"tasks.html",
		}

		tmpl, err := template.ParseFS(ui.Content, files...)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		caser := cases.Title(language.English)
		data := map[string]any{
			"Title":    caser.String("tasks"),
			"Username": "tasks",
			"Path":     r.URL.Path,
		}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	default:
		http.NotFound(w, r)
	}
}
