package external_api

import (
	"flag"
	"github.com/FACorreiaa/go-ollama/internal/api/handler/external_api/health"
	"github.com/FACorreiaa/go-ollama/internal/api/service"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
	"text/template"
)

func initRouter(s *service.Service) *mux.Router {
	var dir string

	flag.StringVar(&dir, "dir", "./internal/api/handler/external_api/static/templates", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()
	handlerHealth := new(health.HandlerHealth)
	//check server
	r.HandleFunc("/health", handlerHealth.HealthCheck).Methods("GET")

	templates := template.Must(template.ParseFiles(
		filepath.Join(dir, "header.tmpl"),
		filepath.Join(dir, "index.html"),
	))

	// Render your template
	//err := templates.ExecuteTemplate(w, "base.html", nil)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//}

	r.PathPrefix("/v1").Handler(http.StripPrefix("/v1", http.FileServer(http.Dir(dir))))

	r.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		//http.ServeFile(w, r, "index.html")
	})
	//http.Handle("/", r)

	return r

}
