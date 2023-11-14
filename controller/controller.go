package controller

import (
	"embed"
	"github.com/FACorreiaa/go-ollama/core/account"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
)

//go:embed static
var staticFS embed.FS

//go:embed html
var htmlFS embed.FS

type core struct {
	accounts *account.Accounts
}

type Handlers struct {
	pgool       *pgxpool.Pool
	formDecoder *form.Decoder
	validator   *validator.Validate
	translator  ut.Translator
	sessions    *sessions.CookieStore
	core        *core
}

func Router(pool *pgxpool.Pool, sessionSecret []byte) http.Handler {
	validate := validator.New()
	translator, _ := ut.New(en.New(), en.New()).GetTranslator("en")
	if err := en_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		slog.Error("Error registering translations", "error", err)

	}
	formDecoder := form.NewDecoder()

	r := mux.NewRouter()
	h := Handlers{
		pgool:       pool,
		formDecoder: formDecoder,
		validator:   validate,
		translator:  translator,
		sessions:    sessions.NewCookieStore(sessionSecret),
		core: &core{
			accounts: account.NewAccounts(pool, validate),
		},
	}

	// Static files
	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(staticFS)))
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		file, _ := staticFS.ReadFile("static/favicon.ico")
		w.Header().Set("Content-Type", "image/x-icon")
		w.Write(file)
	})

	// Public routes, authentication is optional
	optauth := r.NewRoute().Subrouter()
	optauth.Use(h.authMiddleware)
	optauth.HandleFunc("/", handler(h.homePage)).Methods(http.MethodGet)

	// Routes that shouldn't be available to authenticated users
	noauth := r.NewRoute().Subrouter()
	noauth.Use(h.authMiddleware)
	noauth.Use(h.redirectIfAuth)

	noauth.HandleFunc("/login", handler(h.loginPage)).Methods(http.MethodGet)
	noauth.HandleFunc("/login", handler(h.loginPost)).Methods(http.MethodPost)
	noauth.HandleFunc("/register", handler(h.registerPage)).Methods(http.MethodGet)
	noauth.HandleFunc("/register", handler(h.registerPost)).Methods(http.MethodPost)

	// Authenticated routes
	auth := r.NewRoute().Subrouter()
	auth.Use(h.authMiddleware)
	auth.Use(h.requireAuth)

	auth.HandleFunc("/logout", handler(h.logout)).Methods(http.MethodPost)
	auth.HandleFunc("/settings", handler(h.settingsPage)).Methods(http.MethodGet)

	return r
}

func handler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			slog.Error("Error handling request", "error", err)
		}
	}
}

func (h *Handlers) formErrors(err error) []string {
	decodeErrors, isDecodeError := err.(form.DecodeErrors)
	if isDecodeError {
		errors := []string{}
		for _, decodeError := range decodeErrors {
			errors = append(errors, decodeError.Error())
		}

		return errors
	}

	validateErrors, isValidateError := err.(validator.ValidationErrors)
	if isValidateError {
		errors := []string{}
		for _, validateError := range validateErrors {
			errors = append(errors, validateError.Translate(h.translator))
		}
		return errors
	}

	return []string{err.Error()}
}
