package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) CreatePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageName := vars["page"]
	ok, err := h.Wiki.IsPage(pageName)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	if ok {
		fmt.Fprint(w, "페이지가 이미 존재합니다.")
		return
	}
	err = h.Editor(pageName, w)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
}
