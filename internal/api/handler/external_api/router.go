package external_api

import (
	"flag"
	"fmt"
	"github.com/FACorreiaa/go-ollama/internal/api/handler/external_api/health"
	"github.com/FACorreiaa/go-ollama/internal/api/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

var templates *template.Template

func init() {
	templateFiles, err := filepath.Glob("internal/api/handler/external_api/static/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range templateFiles {
		fmt.Println("Template file:", file)
	}

	templates = template.Must(template.ParseGlob("internal/api/handler/external_api/static/templates/*.html"))
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

func initRouter(s *service.Service) *mux.Router {
	var dir string

	flag.StringVar(&dir, "dir", "./internal/api/handler/external_api/static/templates", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()
	handlerHealth := new(health.HandlerHealth)

	r.PathPrefix("/v1").Handler(http.StripPrefix("/v1", http.FileServer(http.Dir(dir))))

	r.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.ServeFile(w, r, "index.html")
	})

	//check server
	r.HandleFunc("/health", handlerHealth.HealthCheck).Methods("GET")
	r.HandleFunc("/", indexHandler)

	return r

}
