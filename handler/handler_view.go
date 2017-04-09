package handler

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/yundream/gowiki/wiki"
	"html/template"
	"net/http"
	"strings"
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

func (h Handler) RenderPage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	pageName := vars["page"]

	err := h.Template.ExecuteTemplate(w, "head", nil)
	if err != nil {
		return err
	}

	dirent := strings.Split(r.URL.Path, "/")
	if dirent[1] == "c" {
		pageName = "editor"
	}

	page, err := h.Wiki.ReadPage(pageName, w, r)
	switch err {
	case wiki.StatusPageNotFound:
		var doc bytes.Buffer
		t, err := template.ParseFiles("plugin/viewer/viewer.tmpl")
		if err != nil {
			w.Write([]byte(err.Error()))
			return err
		}
		a := struct{ PageName string }{pageName}
		err = t.Execute(&doc, a)
		if err != nil {
			w.Write([]byte(err.Error()))
			return err
		}
		w.Write(doc.Bytes())
	case nil:
		w.Write([]byte(page.Contents))
	}
	err = h.Template.ExecuteTemplate(w, "tail", nil)
	if err != nil {
		return err
	}
	return nil
}
