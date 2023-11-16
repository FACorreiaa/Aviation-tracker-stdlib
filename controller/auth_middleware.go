package controller

import (
	"context"
	"fmt"
	"github.com/FACorreiaa/go-ollama/core/account"
	"net/http"
)

type ctxKey int

const (
	ctxKeyAuthUser ctxKey = iota
)

// Middleware to set the current logged in user in the context.
// See `Handlers.requireAuth` or `Handlers.redirectIfAuth` middleware
func (h *Handlers) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.sessions.Get(r, "auth")
		fmt.Println("session", session)

		token := session.Values["token"]
		fmt.Println("token", token)

		if token, ok := token.(string); ok {
			user, err := h.core.accounts.UserFromSessionToken(r.Context(), account.Token(token))
			fmt.Println("user", user)
			fmt.Println("token", token)
			fmt.Println("account.Token(token)", account.Token(token))
			fmt.Println("string(token)", string(token))

			if err == nil {
				ctx := context.WithValue(r.Context(), ctxKeyAuthUser, user)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyAuthUser)
		if user == nil {
			http.Redirect(w, r, "/login?return_to="+r.URL.Path, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) redirectIfAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyAuthUser)
		if user != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
