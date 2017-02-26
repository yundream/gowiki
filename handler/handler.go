package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/yundream/gowiki/wiki"
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
	Wiki              *wiki.Wiki
}

func New() (*Handler, error) {
	w, err := wiki.New("localhost")
	if err != nil {
		return nil, err
	}
	h := &Handler{Wiki: w}
	h.Router = mux.NewRouter()
	h.Router.HandleFunc("/ping", h.Ping).Methods("GET")
	h.Router.HandleFunc("/w/{page:.+}", h.Viewer).Methods("GET")
	http.Handle("/", h.Middleware(h.Router))
	http.Handle("/theme/", http.FileServer(http.Dir("./")))
	return h, nil
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

func (h Handler) RenderPage(w http.ResponseWriter, pageName string) error {
	err := h.Template.ExecuteTemplate(w, "head", nil)
	if err != nil {
		return err
	}
	page, err := h.Wiki.ReadPage(pageName)
	if err != nil {
		return err
	}
	w.Write([]byte(page.Contents))
	err = h.Template.ExecuteTemplate(w, "tail", nil)
	if err != nil {
		return err
	}
	return nil
}
