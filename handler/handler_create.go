package handler

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	//"github.com/yundream/gowiki/wiki"
	"html/template"
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
	err = h.LoadEditor(pageName, w)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
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
	a := struct {
		PageName string
	}{
		PageName: pageName,
	}
	var doc bytes.Buffer
	err = t.Execute(&doc, a)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Fprint(w, doc.String())
	return nil
}
