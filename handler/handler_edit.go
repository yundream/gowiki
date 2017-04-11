package handler

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/yundream/gowiki/wiki"
	"html/template"
	"net/http"
	"strings"
)

func (h *Handler) EditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageName := vars["page"]
	err := h.LoadEditor(pageName, w)
	fmt.Println(err)
}

func (h Handler) LoadEditor(pageName string, w http.ResponseWriter) error {
	err := h.Template.ExecuteTemplate(w, "head", nil)
	if err != nil {
		return err
	}
	defer func() {
		err = h.Template.ExecuteTemplate(w, "tail", nil)
	}()
	t, err := template.ParseFiles("plugin/editor/editor.tmpl")
	if err != nil {
		return err
	}
	page, err := h.Wiki.ReadRawPage(pageName)
	tagStr := strings.Join(page.Tag, " ")
	a := struct {
		wiki.Page
		TagStr string
	}{page, tagStr}
	var doc bytes.Buffer
	err = t.Execute(&doc, a)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Fprint(w, doc.String())
	return nil
}
