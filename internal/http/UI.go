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

func parseTemplates(files ...string) (*template.Template, error) {
	return template.New(files[0]).Funcs(template.FuncMap{
		"contains": strings.Contains,
	}).ParseFS(ui.Content, files...)
}

func (h *HttpServer) UI(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	path := strings.Trim(r.URL.Path, "/")
	caser := cases.Title(language.English)

	switch path {
	case "":
		files := []string{"base.html", "index.html"}
		tmpl, err := parseTemplates(files...)
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
			return
		}
		totalRequests, err := db.GetTotalRequestCount(h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		totalHosts, err := db.GetTotalHostCount(h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		activeTasks, err := db.GetActiveTaskCount(h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		availableWifiKeys, err := db.GetUnassignedWifiKeyCount(h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
			"AvailableWifiKeys":     availableWifiKeys,
		}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "hosts", "hosts/new":
		files := []string{"base.html", "hosts.html"}
		tmpl, err := parseTemplates(files...)
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
			return
		}

		tasks, err := db.GetTasks(h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"Title": caser.String("hosts"),
			"Name":  "User",
			"Path":  r.URL.Path,
			"Hosts": template.HTML(hostsHtml),
			"Tasks": tasks,
		}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "tasks", "tasks/new":
		files := []string{"base.html", "tasks.html"}
		tmpl, err := parseTemplates(files...)
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
			return
		}

		data := map[string]any{
			"Title": caser.String("tasks"),
			"Name":  "User",
			"Path":  r.URL.Path,
			"Tasks": tasksHtml,
		}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "wifikeys", "wifikeys/new":
		files := []string{"base.html", "wifikeys.html"}
		tmpl, err := parseTemplates(files...)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wifiHtml, err := db.GetWifiKeysAsHTML(h.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"Title":    caser.String("tasks"),
			"Name":     "User",
			"Path":     r.URL.Path,
			"WifiKeys": wifiHtml,
		}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	default:
		if strings.HasPrefix(path, "tasks/edit/") {
			id := ps.ByName("id")
			files := []string{"base.html", "tasks_edit.html"}
			tmpl, err := parseTemplates(files...)
			if err != nil {
				if os.IsNotExist(err) {
					http.NotFound(w, r)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			task, err := db.GetTaskByID(id, h.Database)
			if err != nil {
				http.Error(w, "Task not found", http.StatusNotFound)
				return
			}

			data := map[string]any{
				"Title": caser.String("edit task"),
				"Name":  "User",
				"Path":  r.URL.Path,
				"Task":  task,
			}

			if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		} else if strings.HasPrefix(path, "hosts/edit/") {
			id := ps.ByName("id")
			files := []string{"base.html", "hosts_edit.html"}
			tmpl, err := parseTemplates(files...)
			if err != nil {
				if os.IsNotExist(err) {
					http.NotFound(w, r)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			host, err := db.GetHostByID(id, h.Database)
			if err != nil {
				http.Error(w, "Host not found", http.StatusNotFound)
				return
			}

			tasks, err := db.GetTasks(h.Database)
			if err != nil {
				http.Error(w, "Tasks not found", http.StatusNotFound)
				return
			}

			data := map[string]any{
				"Title": caser.String("edit task"),
				"Name":  "User",
				"Path":  r.URL.Path,
				"Host":  host,
				"Tasks": tasks,
			}

			if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		} else if strings.HasPrefix(path, "wifikeys/edit/") {
			id := ps.ByName("id")
			files := []string{"base.html", "wifikeys_edit.html"}
			tmpl, err := parseTemplates(files...)
			if err != nil {
				if os.IsNotExist(err) {
					http.NotFound(w, r)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			key, err := db.GetWifiKeyByID(id, h.Database)
			if err != nil {
				http.Error(w, "Wifi Key not found", http.StatusNotFound)
				return
			}

			data := map[string]any{
				"Title": caser.String("edit wifi key"),
				"Name":  "User",
				"Path":  r.URL.Path,
				"Key":   key,
			}

			if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		http.NotFound(w, r)
	}
}
