package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/oscarsjlh/todo/internal/data"
	mg "github.com/oscarsjlh/todo/migrations"
	"github.com/oscarsjlh/todo/static"
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
	assets, _ := static.Assets()
	fs := http.FileServer(http.FS(assets))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", app.GetTodosHandler)
	http.HandleFunc("/health", app.Health)
	http.HandleFunc("/new-todo", app.InsertTodoHandler)
	http.HandleFunc("/delete/", app.RemoveTodoHandler)
	http.HandleFunc("/update/", app.MarkTodoDoneHandler)
	http.HandleFunc("/modify/", app.EditHandlerForm)
	http.HandleFunc("/edit/", app.EditTodoHandler)
}

func getdgburl() string {
	dbPass := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	production := os.Getenv("PRODUCTION")

	if production != "" {
		dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbPass, dbUser, dbHost, dbName)
		return dsn
	}
	return os.Getenv("TODO_DB_DSN")
}
