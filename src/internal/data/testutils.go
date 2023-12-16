package internal

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	mg "github.com/oscarsjlh/todo/migrations"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var todo *pgxpool.Pool

const (
	user     = "postgres"
	password = "secret"
	db       = "postgres"
	port     = "5433"
	dialect  = "postgres"
)

type TestDB struct {
	DBInstance *pgxpool.Pool
	DbAddress  string
	container  testcontainers.Container
}

func SetUpTestDB() *TestDB {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	container, dbInstance, dbAddr, err := createContainer(ctx)
	if err != nil {
		log.Fatal("Failed to set up test ", err)
	}
	println("Im failing at main func")
	println(dbAddr)
	err = mg.MigrateDb(dbAddr)
	if err != nil {
		log.Fatal("Failed to perform db migration", err)
	}
	cancel()
	return &TestDB{
		container:  container,
		DBInstance: dbInstance,
		DbAddress:  dbAddr,
	}
}

func (tdb *TestDB) TearDown() {
	tdb.DBInstance.Close()
	_ = tdb.container.Terminate(context.Background())
}

func createContainer(ctx context.Context) (testcontainers.Container, *pgxpool.Pool, string, error) {
	env := map[string]string{
		"POSTGRES_PASSWORD": password,
		"POSTGRES_USER":     user,
		"POSTGRES_DB":       db,
	}
	port := "5432/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16.1",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("Failed to start container: %v", err)
	}
	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())
	time.Sleep(time.Second)
	dbAddr := fmt.Sprintf("localhost:%s", p.Port())
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, dbAddr, db)
	print("Var for dsn")
	println(dsn)
	db, err := NewPool(ctx, dsn)
	if err != nil {
		log.Fatal("Failed in pool")
		return container, db, dbAddr, fmt.Errorf("failed to establish database connection: %v", err)
	}
	return container, db, dsn, nil
}
