package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type Options struct {
	Name string
	Age  int
}
type Handler struct {
	Theme             string
	Router            *mux.Router
	DocumentDirectory string
	Template          *template.Template
}

func New() *Handler {
	h := &Handler{}
	h.Router = mux.NewRouter()
	h.Router.HandleFunc("/ping", h.Ping).Methods("GET")
	h.Router.HandleFunc("/w/{page:.+}", h.Viewer).Methods("GET")
	http.Handle("/", h.Middleware(h.Router))
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

func (h *Handler) LoadTemplate(theme string) error {
	var err error
	h.Template, err = template.ParseFiles(
		"theme/"+theme+"/head.html",
		"theme/"+theme+"/tail.html")
	if err != nil {
		return err
	}
	return nil
}

func (h Handler) Render(w http.ResponseWriter, v *interface{}) error {
	err := h.Template.ExecuteTemplate(w, "head", v)
	if err != nil {
		return err
	}
	err = h.Template.ExecuteTemplate(w, "tail", v)
	if err != nil {
		return err
	}
	return nil
}
