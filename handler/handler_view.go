package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) Viewer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageName := vars["page"]
	h.RenderPage(w, pageName)
}
