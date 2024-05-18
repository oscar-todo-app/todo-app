package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	internal "github.com/oscarsjlh/todo/internal/data"
	logs "github.com/oscarsjlh/todo/logging"
	mg "github.com/oscarsjlh/todo/migrations"
	"github.com/oscarsjlh/todo/static"
	"github.com/oscarsjlh/todo/traces"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type application struct {
	todos internal.TodoModel
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := traces.SetupOtelSKD(ctx)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()
	dsn := getdgburl()
	err = mg.MigrateDb(dsn)
	if err != nil {
		slog.Error("Db is not set up properly chekc the env vars")
	}
	if err != nil {
		slog.Error("Failed to migrate DB")
	}
	db, err := internal.NewPool(ctx, dsn)
	if err != nil {
		slog.Error("Failed to set up DB")
	}
	app := &application{
		todos: &internal.Postgres{DB: db},
	}
	srv := &http.Server{
		Addr:         ":3000",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      serverRoutes(app),
	}
	srvErr := make(chan error, 1)
	go func() {
		slog.Info("Server listening in port 3000")
		srvErr <- srv.ListenAndServe()
	}()
	select {
	case err = <-srvErr:
		return
	case <-ctx.Done():
		stop()
	}
	err = srv.Shutdown(context.Background())
	return
}

func serverRoutes(app *application) http.Handler {
	mux := http.NewServeMux()

	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, logs.LoggingMiddleware(handler))
	}

	// use embed for the static files
	assets, _ := static.Assets()
	fs := http.FileServer(http.FS(assets))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	handleFunc("/", app.GetTodosHandler)
	handleFunc("/health", app.Health)
	handleFunc("/new-todo", app.InsertTodoHandler)
	handleFunc("/delete/", app.RemoveTodoHandler)
	handleFunc("/update/", app.MarkTodoDoneHandler)
	handleFunc("/modify/", app.EditHandlerForm)
	handleFunc("/edit/", app.EditTodoHandler)
	handler := otelhttp.NewHandler(mux, "/")
	return handler
}

func getdgburl() string {
	dbPass := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	production := os.Getenv("PRODUCTION")

	if production != "" {
		dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=require", dbUser, dbPass, dbHost, dbName)
		return dsn
	}
	return os.Getenv("TODO_DB_DSN")
}
