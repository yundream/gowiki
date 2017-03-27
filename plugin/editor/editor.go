package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/yundream/gowiki/wiki"
	"html/template"
	"net/http"
)

func Function_editor(data string, parameter string, w http.ResponseWriter, r *http.Request) string {
	vars := mux.Vars(r)
	pageName := vars["page"]
	t, err := template.ParseFiles("plugin/editor/editor.tmpl")
	if err != nil {
		return ""
	}
	fmt.Println("editor:", pageName)

	wiki.ReadRawPage(pageName)

	var doc bytes.Buffer
	a := struct{}{}
	err = t.Execute(&doc, a)
	if err != nil {
		return err.Error()
	}
	return doc.String()
}
