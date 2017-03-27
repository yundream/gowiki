package handler

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/yundream/gowiki/wiki"
	"html/template"
	"net/http"
)

func (h *Handler) Viewer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageName := vars["page"]
	err := h.RenderPage(w, r)
	switch err {
	case wiki.StatusPageNotFound:
		t, err := template.ParseFiles("plugin/viewer/viewer.tmpl")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		var doc bytes.Buffer
		a := struct{ PageName string }{pageName}
		err = t.Execute(&doc, a)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(doc.Bytes())
	}
}
