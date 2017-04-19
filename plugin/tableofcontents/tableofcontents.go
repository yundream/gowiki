package main

import (
	"bytes"
	"html/template"
	"net/http"
)

func Function_tableofcontents(data string, parameter string, w http.ResponseWriter, r *http.Request) string {
	t, err := template.ParseFiles("plugin/tableofcontents/table.tmpl")
	var doc bytes.Buffer
	err = t.Execute(&doc, nil)
	if err != nil {
		return "Table Create Error"
	}
	return doc.String()
}
