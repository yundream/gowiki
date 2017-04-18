package handler

import (
	"github.com/gorilla/mux"
	"github.com/yundream/gowiki/wiki"
	"net/http"
	"strings"
)

func (h *Handler) SavePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageName := vars["page"]
	tag := strings.Fields(r.PostFormValue("tag"))

	contents := wiki.Page{
		Name:     pageName,
		Title:    r.PostFormValue("title"),
		Contents: r.PostFormValue("wikidata"),
		Publish:  true,
		Tag:      tag,
	}
	err := h.Wiki.SavePage(contents)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte("OK"))
	}
}
