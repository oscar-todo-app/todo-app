package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	internal "github.com/oscarsjlh/todo/internal/data"
	ui "github.com/oscarsjlh/todo/templates"
)

type TodoData struct {
	Todos []internal.Todo
}

func renderTemplate(w http.ResponseWriter, tmplName string, data TodoData) error {
	tmpl, err := template.ParseFS(ui.Files, "*.html")
	if err != nil {
		return err
	}
	fmt.Printf("Data is %v", data)
	err = tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
