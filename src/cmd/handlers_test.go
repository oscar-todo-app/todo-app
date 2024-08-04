package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	internal "github.com/oscarsjlh/todo/internal/data"
)

type TestSuite struct {
	app    *application
	server *httptest.Server
}

func (suite *TestSuite) SetUpTest() {
	suite.server = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			suite.app.serverRoutes()
		}),
	)
}

func (suite *TestSuite) TearDownTest() {
	suite.server.Close()
	internal.SetUpTestDB().TearDown()
}

func (suite *TestSuite) InsertTodo() {
	args := "test-todo"
	err := suite.app.todos.InsertTodo(args)
	if err != nil {
		log.Fatal("failed to insert todo")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestHandlers(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T, *TestSuite)
	}{
		{"GetTodos", testGetTodoHandler},
		{"updateHome", testUpdateHome},
		{"InsertTodo", testInsertTodoHandler},
		{"RemoveTodo", testDeleteTodo},
		{"MarkasDone", testMarkAsDone},
		{"UpdateTodo", testUpdateTodo},
		{"Edit", testEditTodo},
	}
	suite := &TestSuite{}
	suite.app = &application{
		todos: &internal.Postgres{DB: internal.SetUpTestDB().DBInstance},
	}
	suite.SetUpTest()
	defer suite.TearDownTest()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t, suite)
		})
	}
}

func getBody(req *bytes.Buffer) string {
	body, err := io.ReadAll(req)
	if err != nil {
		fmt.Print("Unable to convert to string")
	}
	println(body)
	return string(body)
}

func testGetTodoHandler(t *testing.T, suite *TestSuite) {
	suite.InsertTodo()

	req, err := http.NewRequest("GET", suite.server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	suite.app.GetTodosHandler(w, req)
	body := getBody(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, err)

	assert.Contains(t, body, "test-todo")
}

func testUpdateHome(t *testing.T, suite *TestSuite) {
	suite.InsertTodo()
	w := httptest.NewRecorder()
	suite.app.UpdateHomeHandler(w)
	body := getBody(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, body, "test-todo")
}

func testInsertTodoHandler(t *testing.T, suite *TestSuite) {
	queryParam := "insert-todo"
	req, err := http.NewRequest("GET", suite.server.URL+"/new-todo?todo="+queryParam, nil)
	if err != nil {
		println("request faild")
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	suite.app.InsertTodoHandler(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	want, err := suite.app.todos.GetLastInsertedTodo()
	if err != nil {
		return
	}
	assert.Contains(t, queryParam, want.Task_name)
}

func testDeleteTodo(t *testing.T, suite *TestSuite) {
	suite.InsertTodo()
	id, err := suite.app.todos.GetLastInsertedTodo()
	if err != nil {
		log.Fatal("failed to get inserted todo")
	}
	sid := strconv.Itoa(id.Id)

	req, err := http.NewRequest("DELETE", suite.server.URL+"/delete/"+sid, nil)
	if err != nil {
		println("request faild")
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	suite.app.RemoveTodoHandler(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	_, err = suite.app.todos.SelectTodo(id.Id)
	assert.NotNil(t, err, "Expected error as todo should have been removed")
}

func testMarkAsDone(t *testing.T, suite *TestSuite) {
	suite.InsertTodo()
	id, err := suite.app.todos.GetLastInsertedTodo()
	if err != nil {
		log.Fatal("failed to get inserted todo")
	}
	sid := strconv.Itoa(id.Id)

	req, err := http.NewRequest("PUT", suite.server.URL+"/update/"+sid, nil)
	if err != nil {
		println("request faild")
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	suite.app.MarkTodoDoneHandler(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	taskstatus, err := suite.app.todos.SelectTodo(id.Id)
	if err != nil {
		t.Fatal("Failed to select todo by id")
	}
	if taskstatus.Status != true {
		t.Fatal("Expected status to be done")
	}
}

func testUpdateTodo(t *testing.T, suite *TestSuite) {
	suite.InsertTodo()
	id, err := suite.app.todos.GetLastInsertedTodo()
	if err != nil {
		log.Fatal("failed to get inserted todo")
	}
	sid := strconv.Itoa(id.Id)

	req, err := http.NewRequest("GET", suite.server.URL+"/modify/"+sid, nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	suite.app.EditHandlerForm(w, req)
	body := getBody(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, body, "test-todo")
}

func testEditTodo(t *testing.T, suite *TestSuite) {
	suite.InsertTodo()
	id, err := suite.app.todos.GetLastInsertedTodo()
	if err != nil {
		log.Fatal("failed to get inserted todo")
	}
	sid := strconv.Itoa(id.Id)
	formData := "task=modify-todo"
	req, err := http.NewRequest("POST", suite.server.URL+"/edit/"+sid, strings.NewReader(formData))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	suite.app.EditTodoHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	editedtodo, err := suite.app.todos.SelectTodo(id.Id)
	if err != nil {
		log.Fatal("Couldn't not get modified toodo")
	}
	body := getBody(w.Body)
	assert.Contains(t, body, editedtodo.Task_name)
}

func TestValidateIDParam(t *testing.T) {
	for _, tt := range []struct {
		name         string
		idParam      string
		expectedID   int
		expectedErr  string
		expectedCode int
	}{
		{"Valid ID", "123", 123, "", http.StatusOK},
		{"Empty ID", "", 0, "invalid or missing 'id' parameter", http.StatusBadRequest},
		{"Non-numeric ID", "abc", 0, "invalid or missing 'id' parameter", http.StatusBadRequest},
		{"Negative ID", "-5", 0, "invalid 'id' parameter", http.StatusBadRequest},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new test server for each test case
			w := httptest.NewRecorder()

			id, err := validateIDParam(w, tt.idParam)

			// Check the returned ID and error
			if id != tt.expectedID {
				t.Errorf("got %d, want %d", id, tt.expectedID)
			}

			if err != nil && err.Error() != tt.expectedErr {
				t.Errorf("got error %v, want %s", err, tt.expectedErr)
			}

			// Check the HTTP response code
			if w.Code != tt.expectedCode {
				t.Errorf("got HTTP status code %d, want %d", w.Code, tt.expectedCode)
			}
		})
	}
}
