package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) EditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageName := vars["page"]
	err := h.LoadEditor(pageName, w)
	fmt.Println(err)
}
