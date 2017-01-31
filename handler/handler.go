package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	Theme  string
	Router *mux.Router
}

func New() *Handler {
	h := &Handler{}
	h.Router = mux.NewRouter()
	h.Router.HandleFunc("/ping", h.Ping).Methods("GET")
	return h
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World\n")
}
