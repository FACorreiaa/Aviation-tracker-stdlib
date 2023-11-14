package pprof

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/pprof"
)

func InitPprof(router *mux.Router) {
	prefixRouter := router.PathPrefix("/debug/pprof").Subrouter()
	//prefixRouter.Use(middleware.NoCache)
	prefixRouter.HandleFunc("/", gorillaWrapper(pprof.Index)).Methods("GET")
	prefixRouter.HandleFunc("/cmdline", gorillaWrapper(pprof.Cmdline)).Methods("GET")
	prefixRouter.HandleFunc("/profile", gorillaWrapper(pprof.Profile)).Methods("GET")
	prefixRouter.HandleFunc("/symbol", gorillaWrapper(pprof.Symbol)).Methods("GET")
	prefixRouter.HandleFunc("/symbol", gorillaWrapper(pprof.Symbol)).Methods("POST")
	prefixRouter.HandleFunc("/trace", gorillaWrapper(pprof.Trace)).Methods("GET")
	prefixRouter.HandleFunc("/allocs", gorillaWrapper(handlerFunc(pprof.Handler("allocs")))).Methods("GET")
	prefixRouter.HandleFunc("/block", gorillaWrapper(handlerFunc(pprof.Handler("allocs")))).Methods("GET")
	prefixRouter.HandleFunc("/goroutine", gorillaWrapper(handlerFunc(pprof.Handler("allocs")))).Methods("GET")
	prefixRouter.HandleFunc("/heap", gorillaWrapper(handlerFunc(pprof.Handler("allocs")))).Methods("GET")
	prefixRouter.HandleFunc("/mutex", gorillaWrapper(handlerFunc(pprof.Handler("allocs")))).Methods("GET")
	prefixRouter.HandleFunc("/threadcreate", gorillaWrapper(handlerFunc(pprof.Handler("allocs")))).Methods("GET")
}

func handlerFunc(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

func gorillaWrapper(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fn(w, r.WithContext(ctx))
	}
}
