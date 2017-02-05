package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	Theme             string
	Router            *mux.Router
	DocumentDirectory string
}

func New() *Handler {
	h := &Handler{}
	h.Router = mux.NewRouter()
	h.Router.HandleFunc("/ping", h.Ping).Methods("GET")
	http.Handle("/", h.Middleware(h.Router))
	http.Handle("/w/{page:.+}")
	http.Handle("/theme/", http.FileServer(http.Dir("./")))
	return h
}

func (h Handler) Run(port string) error {
	err := http.ListenAndServe(port, nil)
	return err
}

func (h Handler) Middleware(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("URL ", r.URL.Path)
		handle.ServeHTTP(w, r)
	})
}
func (h Handler) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World\n")
}
