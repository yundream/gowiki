package handler

import (
	"net/http"
)

func (h Handler) Viewer(w http.ResponseWriter, r *http.Request) {
	h.Render(w, nil)
}
