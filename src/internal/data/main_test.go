package internal

import (
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testDbInstance *pgxpool.Pool

func TestMain(m *testing.M) {
	testDB := SetUpTestDB()
	testDbInstance = testDB.DBInstance
	defer testDB.TearDown()
	os.Exit(m.Run())
}

func TestTodoOperations(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T, *Postgres)
	}{
		{"InsertTodo", testInsertTodo},
		{"GetTodo", testGetTodo},
		{"DeleteTodo", testDeleteTodo},
		{"MarkTodoAsDone", testMarkTodoAsDone},
		{"EditTodo", testEditTodo},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := SetupTest(t)
			tt.test(t, pg)
		})
	}
}

func SetupTest(t *testing.T) *Postgres {
	pg := NewTodoDS(testDbInstance)
	return pg
}

func testInsertTodo(t *testing.T, pg *Postgres) {
	args := "Test Todo"
	err := pg.InsertTodo(args)
	if err != nil {
		t.Fatalf("InsertTodo failed: %v", err)
	}

	lastTodo, err := pg.GetLastInsertedTodo()
	if err != nil {
		t.Fatalf("Error getting last inserted todo: %v", err)
	}

	assertEqual(t, lastTodo.Task_name, args, "Task_name mismatch")
	assertEqual(t, lastTodo.Status, false, "Status mismatch")
}

func testGetTodo(t *testing.T, pg *Postgres) {
	args := "Test Todo"
	err := pg.InsertTodo(args)
	if err != nil {
		t.Fatalf("InsertTodo failed: %v", err)
	}

	todos, err := pg.GetTodo()
	if err != nil {
		t.Fatalf("Get todo failed: %v", err)
	}

	assertGreaterThanOrEqual(t, len(todos), 1, "Get todo length should be at least 1")
}

func testDeleteTodo(t *testing.T, pg *Postgres) {
	args := "TestTodo1"
	err := pg.InsertTodo(args)
	if err != nil {
		t.Fatalf("InsertTodo failed: %v", err)
	}

	err = pg.RemoveTodo(1)
	if err != nil {
		t.Fatalf("RemoveTodo failed: %v", err)
	}

	_, err = pg.SelectTodo(1)
	assertNotNil(t, err, "Expected error as todo should have been removed")
}

func testMarkTodoAsDone(t *testing.T, pg *Postgres) {
	args := "TestTodo1"
	err := pg.InsertTodo(args)
	if err != nil {
		t.Fatalf("InsertTodo failed: %v", err)
	}

	lastTodo, err := pg.GetLastInsertedTodo()
	if err != nil {
		t.Fatalf("Error getting last inserted todo: %v", err)
	}

	err = pg.MarkTodoAsDone(lastTodo.Id)
	if err != nil {
		t.Fatalf("MarkTodo as done failed: %v", err)
	}

	todo, err := pg.SelectTodo(lastTodo.Id)
	assertEqual(t, todo.Status, true, "Expected todo status to be true")
}

func testEditTodo(t *testing.T, pg *Postgres) {
	args := "Testing todoapp"
	err := pg.InsertTodo(args)
	if err != nil {
		t.Fatalf("InsertTodo failed: %v", err)
	}

	lastTodo, err := pg.GetLastInsertedTodo()
	if err != nil {
		t.Fatalf("Error getting last inserted todo: %v", err)
	}

	editedArg := "Edited todo"
	pg.EditTodo(lastTodo.Id, editedArg)

	editedTodo, err := pg.SelectTodo(lastTodo.Id)
	if err != nil {
		t.Fatalf("Error getting last inserted todo: %v", err)
	}

	assertEqual(t, editedTodo.Task_name, editedArg, "Editing todo failed")
}

// // Helper functions for assertion
func assertEqual(t *testing.T, actual, expected interface{}, message string) {
	if actual != expected {
		t.Errorf("%s: expected %v, got %v", message, expected, actual)
	}
}

func assertGreaterThanOrEqual(t *testing.T, value, min int, message string) {
	if value < min {
		t.Errorf(
			"%s: expected value to be greater than or equal to %d, got %d",
			message,
			min,
			value,
		)
	}
}

func assertNotNil(t *testing.T, obj interface{}, message string) {
	if obj == nil {
		t.Errorf("%s: expected non-nil value", message)
	}
}
