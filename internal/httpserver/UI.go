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

		graphData1, graphData2, graphDates, err := db.GetRequestGraphData(time.Now(), h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		totalRequests, err := db.GetTotalRequestCount(h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		totalHosts, err := db.GetTotalHostCount(h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		activeTasks, err := db.GetActiveTaskCount(h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		caser := cases.Title(language.English)
		data := map[string]any{
			"Title":                 caser.String("Home"),
			"Name":                  "User",
			"Path":                  r.URL.Path,
			"RegisteredGraphData":   graphData1,
			"UnregisteredGraphData": graphData2,
			"GraphDates":            graphDates,
			"Month":                 time.Now().Month(),
			"TotalRequests":         totalRequests,
			"TotalHosts":            totalHosts,
			"ActiveTasks":           activeTasks,
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

		hostsHtml, err := db.GetHostsAsHTML(h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		caser := cases.Title(language.English)
		data := map[string]any{
			"Title": caser.String("hosts"),
			"Name":  "User",
			"Path":  r.URL.Path,
			"Hosts": template.HTML(hostsHtml),
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

		tasksHtml, err := db.GetTasksAsHTML(h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		caser := cases.Title(language.English)
		data := map[string]any{
			"Title":    caser.String("tasks"),
			"Username": "tasks",
			"Path":     r.URL.Path,
			"Tasks":    tasksHtml,
		}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	default:
		http.NotFound(w, r)
	}
}
