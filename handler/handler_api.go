package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) CallAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fname := vars["name"]
	rtv, err := h.P.Exec(fname, "", w, r)
	if err != nil {
		w.Write([]byte("error"))
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(rtv))
}
