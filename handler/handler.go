package handler

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/yundream/gowiki/wiki"
	"html/template"
	"net/http"
	"strings"
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
	h.Router.HandleFunc("/c/{page:.+}", h.CreatePage).Methods("GET")
	http.Handle("/", h.Middleware(h.Router))
	http.Handle("/theme/", http.FileServer(http.Dir("./")))
	return h, nil
}

func (h Handler) Run(port string) error {
	fmt.Println("Application running... ", port)
	err := http.ListenAndServe(port, nil)
	return err
}

func (h Handler) Middleware(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("URL ", r.URL.Path)
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

func (h Handler) LoadEditor(pageName string, w http.ResponseWriter) error {
	err := h.Template.ExecuteTemplate(w, "head", nil)
	if err != nil {
		return err
	}

	err = h.Template.ExecuteTemplate(w, "tail", nil)
	if err != nil {
		return err
	}

	return nil
}
func (h Handler) RenderPage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	pageName := vars["page"]

	err := h.Template.ExecuteTemplate(w, "head", nil)
	if err != nil {
		return err
	}

	dirent := strings.Split(r.URL.Path, "/")
	if dirent[1] == "c" {
		pageName = "editor"
	}

	page, err := h.Wiki.ReadPage(pageName, w, r)
	switch err {
	case wiki.StatusPageNotFound:
		var doc bytes.Buffer
		t, err := template.ParseFiles("plugin/viewer/viewer.tmpl")
		if err != nil {
			w.Write([]byte(err.Error()))
			return err
		}
		a := struct{ PageName string }{pageName}
		err = t.Execute(&doc, a)
		if err != nil {
			w.Write([]byte(err.Error()))
			return err
		}
		w.Write(doc.Bytes())
	case nil:
		w.Write([]byte(page.Contents))
	}
	err = h.Template.ExecuteTemplate(w, "tail", nil)
	if err != nil {
		return err
	}
	return nil
}
