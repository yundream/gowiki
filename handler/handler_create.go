package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/yundream/gowiki/wiki"
	"net/http"
)

func (h *Handler) CreatePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageName := vars["page"]
	err := h.LoadEditor(pageName, w)
	switch err {
	case wiki.StatusPageNotFound:
		fmt.Println("Editor Page not Found", pageName)
	}
}
