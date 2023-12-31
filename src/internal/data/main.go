package internal

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type TodoModel interface {
	InsertTodo(args string) error
	GetTodo() ([]Todo, error)
	RemoveTodo(id int) error
	MarkTodoAsDone(id int) error
	EditTodo(id int, task_name string) error
	SelectTodo(id int) (*Todo, error)
	GetLastInsertedTodo() (*Todo, error)
	Ping(ctx context.Context) error
}
type Todo struct {
	Id        int
	Task_name string
	Status    bool
}

type Postgres struct {
	DB *pgxpool.Pool
}

func NewTodoDS(DB *pgxpool.Pool) *Postgres {
	return &Postgres{DB: DB}
}

func NewPool(ctx context.Context, dsl string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, dsl)
	if err != nil {
		log.Fatal("unable to create db pool", err)
		return nil, err
	}

	return db, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.DB.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.DB.Close()
}

func (pg *Postgres) InsertTodo(args string) error {
	query := `INSERT INTO tasks ("task_name", "status") VALUES ($1, FALSE)`
	_, err := pg.DB.Exec(context.Background(), query, args)
	if err != nil {
		log.Fatal("unable to insert todo", err)
		return err
	}
	return nil
}

func (pg *Postgres) GetTodo() ([]Todo, error) {
	t := []Todo{}
	rows, err := pg.DB.Query(context.Background(), "select * from tasks order by id desc")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return nil, err
	}
	for rows.Next() {
		var id int
		var task string
		var done bool
		err := rows.Scan(&id, &task, &done)
		if err != nil {
			log.Fatal("unable to get tasks")
			return nil, err
		}
		todo := Todo{
			Id:        id,
			Task_name: task,
			Status:    done,
		}
		t = append(t, todo)
	}
	return t, nil
}

func (pg *Postgres) RemoveTodo(id int) error {
	query := `DELETE FROM tasks WHERE ID = $1`

	_, err := pg.DB.Exec(context.Background(), query, id)
	if err != nil {
		log.Fatal("unable to remove todo")
		return err
	}

	return nil
}

func (pg *Postgres) MarkTodoAsDone(id int) error {
	query := `UPDATE tasks SET status = TRUE WHERE id = $1`

	_, err := pg.DB.Exec(context.Background(), query, id)
	if err != nil {
		log.Fatal("unable to mark todo as done")
		return err
	}

	return nil
}

func (pg *Postgres) EditTodo(id int, task_name string) error {
	query := `UPDATE tasks SET task_name = $2 WHERE id = $1`

	_, err := pg.DB.Exec(context.Background(), query, id, task_name)
	if err != nil {
		log.Fatal("unable to edit todo")
		return err
	}

	return nil
}

func (pg *Postgres) GetLastInsertedTodo() (*Todo, error) {
	query := "SELECT id, task_name, status FROM tasks ORDER BY id DESC LIMIT 1"

	var todo Todo
	err := pg.DB.QueryRow(context.Background(), query).Scan(&todo.Id, &todo.Task_name, &todo.Status)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (pg *Postgres) SelectTodo(id int) (*Todo, error) {
	query := `SELECT task_name, id, status FROM tasks WHERE id = $1`
	var todo Todo
	err := pg.DB.QueryRow(context.Background(), query, id).
		Scan(&todo.Task_name, &todo.Id, &todo.Status)
	if err != nil {
		log.Println("Error executing query:", err)
		return &Todo{}, err
	}

	return &todo, nil
}
