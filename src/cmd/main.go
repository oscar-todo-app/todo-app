package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	internal "github.com/oscarsjlh/todo/internal/data"
	mg "github.com/oscarsjlh/todo/migrations"
	"github.com/oscarsjlh/todo/static"
)

type application struct {
	todos  internal.TodoModel
	logger *slog.Logger
}

func main() {
	ctx := context.Context(context.Background())
	dsn := getdgburl()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	err := mg.MigrateDb(dsn)
	if err != nil {
		log.Fatal("Db is not set up properly chekc the env vars")
	}
	if err != nil {
		log.Fatal("Failed to migrate DB")
	}
	db, err := internal.NewPool(ctx, dsn)
	if err != nil {
		log.Fatal("Failed to set up DB")
	}
	app := &application{
		todos:  &internal.Postgres{DB: db},
		logger: logger,
	}
	port := ":3000"
	logger.Info("Starting server", "addr", port)
	err = http.ListenAndServe(port, app.serverRoutes())
	if err != nil {
		log.Fatal("Unable to start http server")
	}
}

func (app *application) serverRoutes() http.Handler {
	// use embed for the static files
	mux := http.NewServeMux()
	assets, _ := static.Assets()
	fs := http.FileServer(http.FS(assets))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", app.GetTodosHandler)
	mux.HandleFunc("/health", app.Health)
	mux.HandleFunc("/new-todo", app.InsertTodoHandler)
	mux.HandleFunc("/delete/", app.RemoveTodoHandler)
	mux.HandleFunc("/update/", app.MarkTodoDoneHandler)
	mux.HandleFunc("/modify/", app.EditHandlerForm)
	mux.HandleFunc("/edit/", app.EditTodoHandler)
	return app.logRequest(commonHeaders(mux))
}

func getdgburl() string {
	dbPass := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	production := os.Getenv("PRODUCTION")

	if production != "" {
		dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbName)
		return dsn
	}
	return os.Getenv("TODO_DB_DSN")
}
