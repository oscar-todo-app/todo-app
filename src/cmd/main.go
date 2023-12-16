package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/oscarsjlh/todo/internal/data"
	mg "github.com/oscarsjlh/todo/migrations"
)

type application struct {
	todos internal.TodoModel
}

func main() {
	ctx := context.Context(context.Background())
	dsn := os.Getenv("TODO_DB_DSN")
	err := mg.MigrateDb(dsn)
	if err != nil {
		log.Fatal("Failed to migrate DB")
	}
	db, err := internal.NewPool(ctx, dsn)
	if err != nil {
		return
	}
	app := &application{
		todos: &internal.Postgres{DB: db},
	}
	serverRoutes(app)
	http.ListenAndServe(":3000", nil)
}

func serverRoutes(app *application) {
	// use embed for the static files

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", app.GetTodosHandler)
	http.HandleFunc("/new-todo", app.InsertTodoHandler)
	http.HandleFunc("/delete/", app.RemoveTodoHandler)
	http.HandleFunc("/update/", app.MarkTodoDoneHandler)
	http.HandleFunc("/modify/", app.EditHandlerForm)
	http.HandleFunc("/edit/", app.EditTodoHandler)
}
