package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/yundream/gowiki/plugin"
	"github.com/yundream/gowiki/sessions"
	"github.com/yundream/gowiki/wiki"
	"html/template"
	"net/http"
)

type Options struct {
	Name string
	Age  int
}

type DocInfo struct {
	PageName string
	Session  sessions.SessionData
}

type Handler struct {
	Theme             string
	Router            *mux.Router
	DocumentDirectory string
	Template          *template.Template
	Wiki              *wiki.Wiki
	P                 *plugin.PlugIns
}

func New() (*Handler, error) {
	p, err := plugin.Load()
	if err != nil {
		return nil, err
	}
	w, err := wiki.New("localhost", p)
	if err != nil {
		return nil, err
	}
	h := &Handler{Wiki: w, P: p}
	h.Router = mux.NewRouter()
	h.Router.HandleFunc("/ping", h.Ping).Methods("GET")
	h.Router.HandleFunc("/w/{page:.+}", h.Viewer).Methods("GET")
	h.Router.HandleFunc("/c/{page:.+}", h.CreatePage).Methods("GET")
	h.Router.HandleFunc("/e/{page:.+}", h.EditPage).Methods("GET")
	h.Router.HandleFunc("/s/{page:.+}", h.SavePage).Methods("POST")
	h.Router.HandleFunc("/api/{name}", h.CallAPI).Methods("POST", "GET", "DELETE")
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

func GetJwt(r *http.Request) sessions.SessionData {
	cookie, _ := r.Cookie("session-jwt")
	info := sessions.SessionData{}
	if cookie == nil {
		return info
	}
	info, _ = sessions.Validation(cookie.Value)
	return info
}
