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

func serveLogin() {

}

func configureAdminHandler() http.Handler {
	r := mux.NewRouter()
	apiRouter := r.NewRoute().PathPrefix("/").Subrouter()
	apiRouter.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
	//apiRouter.HandleFunc("/", serveLogin)
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

		// Configure TLS if certs are available.
	err := adminServer.ListenAndServe()
	if err!= nil {
		fmt.Println("Admin server error.")
	}
	fmt.Println("Admin server running...")
}

func stopAdminServer() {
	//TODO
}

func main() {
	startAdminServer()
}