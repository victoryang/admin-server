package main

import (
	"net/http"
	"net/http/pprof"
	"time"

	"fmt"

	"github.com/gorilla/mux"
)

const (
	host = "0.0.0.0"
	port = "9090"
)

type handlerFunc func(http.Handler) http.Handler
var handlerFns = []handlerFunc{
//	SetJwtMiddlewareHandler,
}

func RegisterHandlers (r *mux.Router, handlerFns ...handlerFunc) http.Handler {
	var f http.Handler
	f =r
	for _, hFn := range handlerFns {
		f = hFn(f)
	}
	return f
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle login")
	vars := mux.Vars(r)
	auth := vars["Authorization"]
	w.WriteHeader(http.StatusOK)
	if vars == "" {
		fmt.Fprintf(w, "auth fails")
		return
	}
	fmt.Fprintf(w, "auth good")
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	
}

func configureAdminHandler() http.Handler {
	r := mux.NewRouter()
	apiRouter := r.NewRoute().PathPrefix("/").Subrouter()
	//apiRouter.HandleFunc("/", handleHome).Methods("GET")
	apiRouter.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/")))).Methods("Get")
	apiRouter.HandleFunc("/", handleLogin).Methods("GET").Queries("username", "pass")
	
	/*TODO: get some status of controller back to admin*/
	/*admin := apiRouter.PathPrefix("/admin").Subrouter()
	admin.Methods("GET").Path("/usage").HandlerFunc(SetJwtMiddlewareFunc(getUsage))*/

	apiRouter.Path("/debug/cmdline").HandlerFunc(pprof.Cmdline)
	apiRouter.Path("/debug/profile").HandlerFunc(pprof.Profile)
	apiRouter.Path("/debug/symbol").HandlerFunc(pprof.Symbol)
	apiRouter.Path("/debug/trace").HandlerFunc(pprof.Trace)
	apiRouter.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)

	return RegisterHandlers(r, handlerFns...)
}

func startAdminServer() {
	adminServer := &http.Server{
		Addr: 			":9000",
		// Adding timeout of 10 minutes for unresponsive client connections.
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		Handler:        configureAdminHandler(),
		MaxHeaderBytes: 1 << 20,
	}

	quit := make(chan bool)
		// Configure TLS if certs are available.
	go func() {
		adminServer.ListenAndServe()
		quit<-
	}()
	fmt.Println("Admin server running...")
	<-quit
	fmt.Println("Admin server stop...")
}

func stopAdminServer() {
	//TODO
}

func main() {
	startAdminServer()
}